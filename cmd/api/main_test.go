package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	godotenv.Load("../../.test.env")
	loadEnv()

	assert.Equal(t, "testing", os.Getenv("ENV"))
}

func TestGetJobNotifier(t *testing.T) {
	godotenv.Load("../../.test.env")
	notifier := getJobNotifier()

	assert.NotNil(t, notifier)
}

func TestGetCertExpirationWarningDays(t *testing.T) {
	godotenv.Load("../../.test.env")
	warningDays := getCertExpirationWarningDays()

	assert.Equal(t, 30, warningDays)
}

func TestGetCORSOrigins(t *testing.T) {
	godotenv.Load("../../.test.env")
	origins := getCORSOrigins()

	assert.Equal(t, "https://localhost", origins)
}
