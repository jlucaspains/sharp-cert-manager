package handlers

import (
	"net/http"
	"testing"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/stretchr/testify/assert"
)

func TestInitTemplate(t *testing.T) {
	templatePath = "../frontend"
	initTemplates()

	assert.NotNil(t, indexTemplate.Lookup("index.html"))
	assert.NotNil(t, indexTemplate.Lookup("body.html"))
	assert.NotNil(t, indexTemplate.Lookup("head.html"))
	assert.NotNil(t, indexTemplate.Lookup("item.html"))
	assert.NotNil(t, indexTemplate.Lookup("itemLoaded.html"))
	assert.NotNil(t, indexTemplate.Lookup("itemModal.html"))
}

func TestRendersIndex(t *testing.T) {
	templatePath = "../frontend"
	handlers := new(Handlers)
	handlers.CertList = []models.CheckCertItem{
		{Name: "blog.lpains.net", Url: "https://blog.lpains.net", Type: models.CertCheckURL},
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /", handlers.Index)

	code, _, body, _, err := makeRequest[string](router, "GET", "/", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	// Testing the full HTML content is not practical, so we check for specific elements
	assert.Contains(t, body, "data-testid=\"result-item\"")
	assert.Contains(t, body, "hx-get=\"/item?name=blog.lpains.net\"")
}

func TestRendersItem(t *testing.T) {
	templatePath = "../frontend"
	handlers := new(Handlers)
	handlers.CertList = []models.CheckCertItem{
		{Name: "blog.lpains.net", Url: "https://blog.lpains.net", Type: models.CertCheckURL},
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /item", handlers.GetItem)

	code, _, body, _, err := makeRequest[string](router, "GET", "/item?name=blog.lpains.net", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	// Testing the full HTML content is not practical, so we check for specific elements
	assert.Contains(t, body, "hx-get=\"/itemDetail?name=blog.lpains.net\" hx-trigger=\"click, keyup[key=='Enter']\" hx-target=\"#modal\"")
	assert.Contains(t, body, "<h2 class=\"text-white text-lg font-medium\">blog.lpains.net</h2>")
}

func TestRendersItemDetail(t *testing.T) {
	templatePath = "../frontend"
	handlers := new(Handlers)
	handlers.CertList = []models.CheckCertItem{
		{Name: "blog.lpains.net", Url: "https://blog.lpains.net", Type: models.CertCheckURL},
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /itemDetail", handlers.GetItemDetail)

	code, _, body, _, err := makeRequest[string](router, "GET", "/itemDetail?name=blog.lpains.net", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	// Testing the full HTML content is not practical, so we check for specific elements
	assert.Contains(t, body, "<td class=\"px-4 py-2\">blog.lpains.net</td>")
}
