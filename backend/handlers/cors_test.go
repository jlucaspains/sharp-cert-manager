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
	router.HandleFunc("OPTIONS /api/check-url", handlers.CORS)

	code, body, err, headers := makeRequest[string](router, "OPTIONS", "/api/check-url", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Empty(t, body)
	assert.Equal(t, "http://localhost:5173", headers["Access-Control-Allow-Origin"][0])
}
func TestCORSGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.SiteList = shared.GetConfigSites()
	handlers.CORSOrigins = "http://localhost:5173"

	router := http.NewServeMux()
	router.HandleFunc("GET /site-list", handlers.GetSiteList)

	code, body, err, headers := makeRequest[[]models.CheckListResult](router, "GET", "/site-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(*body))
	assert.Equal(t, "http://localhost:5173", headers["Access-Control-Allow-Origin"][0])
}

func TestNoCORSGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.SiteList = shared.GetConfigSites()

	router := http.NewServeMux()
	router.HandleFunc("GET /site-list", handlers.GetSiteList)

	code, body, err, headers := makeRequest[[]models.CheckListResult](router, "GET", "/site-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(*body))
	assert.Empty(t, headers["Access-Control-Allow-Origin"])
}
