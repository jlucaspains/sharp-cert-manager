package handlers

import (
	"net/http"
	"testing"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/jlucaspains/sharp-cert-manager/shared"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCORSCheckURL(t *testing.T) {
	handlers := new(Handlers)
	handlers.CORSOrigins = "http://localhost:5173"

	router := http.NewServeMux()
	router.HandleFunc("OPTIONS /api/check-cert", handlers.CORS)

	code, body, err, headers := makeRequest[string](router, "OPTIONS", "/api/check-cert", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Empty(t, body)
	assert.Equal(t, "http://localhost:5173", headers["Access-Control-Allow-Origin"][0])
}

func TestCORSGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.CertList = shared.GetConfigCerts()
	handlers.CORSOrigins = "http://localhost:5173"

	router := http.NewServeMux()
	router.HandleFunc("GET /cert-list", handlers.GetCertList)

	code, body, err, headers := makeRequest[[]models.CheckCertItem](router, "GET", "/cert-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 2, len(*body))
	assert.Equal(t, "http://localhost:5173", headers["Access-Control-Allow-Origin"][0])
}

func TestNoCORSGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.CertList = shared.GetConfigCerts()

	router := http.NewServeMux()
	router.HandleFunc("GET /cert-list", handlers.GetCertList)

	code, body, err, headers := makeRequest[[]models.CheckCertItem](router, "GET", "/cert-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 2, len(*body))
	assert.Empty(t, headers["Access-Control-Allow-Origin"])
}
