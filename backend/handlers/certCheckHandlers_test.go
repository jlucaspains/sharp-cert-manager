package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/jlucaspains/sharp-cert-manager/shared"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.SiteList = shared.GetConfigSites()
	router := mux.NewRouter()
	router.HandleFunc("/site-list", handlers.GetSiteList).Methods("GET")

	code, body, err := makeRequest[[]models.CheckListResult](router, "GET", "/site-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(*body))
}

func TestGetCheckStatus(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/check-status", handlers.CheckStatus).Methods("GET")

	url := fmt.Sprintf("/check-status?url=%s", "https://blog.lpains.net")
	code, body, err := makeRequest[models.CertCheckResult](router, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, body.IsValid)
	assert.LessOrEqual(t, body.CertStartDate, time.Now())
	assert.GreaterOrEqual(t, body.CertEndDate, time.Now())
	assert.Contains(t, body.Hostname, "blog.lpains.net")
	assert.Contains(t, body.CertDnsNames, "blog.lpains.net")
}

func TestGetCheckStatusNoUrl(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/check-status", handlers.CheckStatus).Methods("GET")

	url := "/check-status"
	code, body, err := makeRequest[models.ErrorResult](router, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 400, code)
	assert.Equal(t, "url is required", body.Errors[0])
}
