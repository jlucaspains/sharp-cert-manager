package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jlucaspains/sharp-cert-checker/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.SiteList = handlers.GetConfigSites()
	app := fiber.New()
	defer app.Shutdown()
	app.Get("/site-list", handlers.GetSiteList)

	code, respBody, err := getJsonTestRequestResponse[[]models.CertCheckResult](app, "GET", "/site-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(respBody))
}

func TestGetCheckStatus(t *testing.T) {
	handlers := new(Handlers)
	app := fiber.New()
	defer app.Shutdown()

	app.Get("/check-status", handlers.CheckStatus)

	url := fmt.Sprintf("/check-status?url=%s", "https://blog.lpains.net")
	code, respBody, err := getJsonTestRequestResponse[models.CertCheckResult](app, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, respBody.IsValid)
	assert.LessOrEqual(t, respBody.CertStartDate, time.Now().String())
	assert.GreaterOrEqual(t, respBody.CertEndDate, time.Now().String())
	assert.Contains(t, respBody.Hostname, "blog.lpains.net")
	assert.Contains(t, respBody.CertDnsNames, "blog.lpains.net")
}

func TestGetCheckStatusNoUrl(t *testing.T) {
	handlers := new(Handlers)
	app := fiber.New()
	defer app.Shutdown()

	app.Get("/check-status", handlers.CheckStatus)

	url := "/check-status"
	code, respBody, err := getJsonTestRequestResponse[map[string]string](app, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 400, code)
	assert.Equal(t, respBody["error"], "Missing URL parameter")
}

func TestGetCheckStatusHttp(t *testing.T) {
	handlers := new(Handlers)
	app := fiber.New()
	defer app.Shutdown()

	app.Get("/check-status", handlers.CheckStatus)

	url := fmt.Sprintf("/check-status?url=%s", "http://blog.lpains.net")
	code, respBody, err := getJsonTestRequestResponse[models.CertCheckResult](app, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.False(t, respBody.IsValid)
}

func getJsonTestRequestResponse[K any | []any](app *fiber.App, method string, url string, reqBody any) (code int, respBody K, err error) {
	bodyJson, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(method, url, bytes.NewReader(bodyJson))
	resp, err := app.Test(req, 1000)
	// If error we're done
	if err != nil {
		return
	}
	code = resp.StatusCode
	// If no body content, we're done
	if resp.ContentLength == 0 {
		return
	}
	bodyData := make([]byte, resp.ContentLength)
	_, _ = resp.Body.Read(bodyData)
	err = json.Unmarshal(bodyData, &respBody)
	return
}
