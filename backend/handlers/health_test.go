package handlers

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	code, body, err := makeRequest[models.HealthResult](router, "GET", "/health", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, body.Healthy)
}
