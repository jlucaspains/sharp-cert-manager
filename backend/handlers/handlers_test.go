package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestErrorTranslationSuccess(t *testing.T) {
	type User struct {
		Username string `validate:"required"`
		Tagline  string `validate:"required,lt=10"`
		Tagline2 int    `validate:"required,gt=1"`
	}

	user := User{
		Username: "Joeybloggs",
		Tagline:  "Works",
		Tagline2: 1,
	}

	validate := validator.New()
	err := validate.Struct(user)
	handlers := new(Handlers)
	code, result := handlers.ErrorToHttpResult(err)

	assert.Equal(t, http.StatusBadRequest, code)

	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "Tagline2 should be greater than 1", result.Errors[0])
}

func TestErrorTranslationServerError(t *testing.T) {
	handlers := new(Handlers)
	code, result := handlers.ErrorToHttpResult(fmt.Errorf("Something went wrong"))
	assert.Equal(t, http.StatusInternalServerError, code)

	assert.Equal(t, "Unknown error", result.Errors[0])
}

func makeRequest[K any | []any](router *mux.Router, method string, url string, body any) (code int, respBody *K, err error) {
	inputBody := ""

	if body != nil {
		inputBodyJson, _ := json.Marshal(body)
		inputBody = string(inputBodyJson)
	}

	req, _ := http.NewRequest(method, url, bytes.NewReader([]byte(inputBody)))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	result := new(K)
	err = json.Unmarshal(rr.Body.Bytes(), &result)

	return rr.Code, result, err
}
