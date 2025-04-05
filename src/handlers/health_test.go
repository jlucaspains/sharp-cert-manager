package handlers

import (
	"net/http"
	"testing"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	handlers := new(Handlers)

	router := http.NewServeMux()
	router.HandleFunc("GET /health", handlers.HealthCheck)

	code, body, err, _ := makeRequest[models.HealthResult](router, "GET", "/health", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, body.Healthy)
}
