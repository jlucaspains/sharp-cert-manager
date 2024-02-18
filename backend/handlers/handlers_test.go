package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestErrorTranslationSuccess(t *testing.T) {
	type TestStruct struct {
		Req   string `validate:"required"`
		Lt    string `validate:"required,lt=10"`
		Lte   string `validate:"required,lte=1"`
		Gt    int    `validate:"required,gt=1"`
		Gte   int    `validate:"required,gte=10"`
		Min   string `validate:"min=10"`
		Max   string `validate:"max=9"`
		Alpha string `validate:"alpha"`
	}

	user := TestStruct{
		Req:   "",
		Lt:    "0123456789",
		Lte:   "012345678",
		Gt:    1,
		Gte:   1,
		Min:   "012345678",
		Max:   "0123456789",
		Alpha: "0123456789",
	}

	validate := validator.New()
	err := validate.Struct(user)
	handlers := new(Handlers)
	code, result := handlers.ErrorToHttpResult(err)

	assert.Equal(t, http.StatusBadRequest, code)

	assert.Len(t, result.Errors, 8)
	assert.Equal(t, "Req is required", result.Errors[0])
	assert.Equal(t, "Lt should be less than 10", result.Errors[1])
	assert.Equal(t, "Lte should be less than or equal to 1", result.Errors[2])
	assert.Equal(t, "Gt should be greater than 1", result.Errors[3])
	assert.Equal(t, "Gte should be greater than or equal to 10", result.Errors[4])
	assert.Equal(t, "Min should have minimum length of 10", result.Errors[5])
	assert.Equal(t, "Max should have maximum length of 9", result.Errors[6])
	assert.Equal(t, "Alpha should contain alpha characters only", result.Errors[7])
}

func TestErrorTranslationServerError(t *testing.T) {
	handlers := new(Handlers)
	code, result := handlers.ErrorToHttpResult(fmt.Errorf("Something went wrong"))
	assert.Equal(t, http.StatusInternalServerError, code)

	assert.Equal(t, "Unknown error", result.Errors[0])
}

func makeRequest[K any | []any](router *http.ServeMux, method string, url string, body any) (code int, respBody *K, err error, headers http.Header) {
	inputBody := ""

	if body != nil {
		inputBodyJson, _ := json.Marshal(body)
		inputBody = string(inputBodyJson)
	}

	req, _ := http.NewRequest(method, url, bytes.NewReader([]byte(inputBody)))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	result := new(K)

	switch any(result).(type) {
	case *string:
		// do nothing as we don't care about string
	default:
		err = json.Unmarshal(rr.Body.Bytes(), &result)
	}

	return rr.Code, result, err, rr.Result().Header
}
