package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/jlucaspains/sharp-cert-manager/shared"
)

var indexTemplate *template.Template
var templatePath string = "frontend"

func initTemplates() {
	if indexTemplate != nil {
		return
	}

	indexTemplate = template.Must(template.ParseGlob(fmt.Sprintf("%s/*", templatePath)))
}

func (h Handlers) Index(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	err := indexTemplate.ExecuteTemplate(w, "index.html", h.CertList)

	handleError(w, err)
}

func (h Handlers) GetItem(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	name, _ := h.getQueryParam(r, "name")

	log.Println("Received get item for name: " + name)

	if name == "" {
		h.HTML(w, http.StatusBadRequest, "name is required")
		return
	}

	idx := slices.IndexFunc(h.CertList, func(c models.CheckCertItem) bool { return c.Name == name })

	if idx < 0 {
		h.HTML(w, http.StatusBadRequest, "the provided cert name is not configured")
		return
	}

	item := h.CertList[idx]
	result, err := shared.CheckCertStatus(item, h.ExpirationWarningDays)

	if err != nil {
		handleError(w, err)
		return
	}

	err = indexTemplate.ExecuteTemplate(w, "itemLoaded.html", result)

	handleError(w, err)
}

func (h Handlers) GetItemDetail(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	name, _ := h.getQueryParam(r, "name")

	log.Println("Received detail message for name: " + name)

	if name == "" {
		h.HTML(w, http.StatusBadRequest, "name is required")
		return
	}

	idx := slices.IndexFunc(h.CertList, func(c models.CheckCertItem) bool { return c.Name == name })

	if idx < 0 {
		h.HTML(w, http.StatusBadRequest, "the provided cert name is not configured")
		return
	}

	item := h.CertList[idx]
	result, err := shared.CheckCertStatus(item, h.ExpirationWarningDays)

	if err != nil {
		handleError(w, err)
		return
	}

	err = indexTemplate.ExecuteTemplate(w, "itemModal.html", result)

	handleError(w, err)
}

func (h Handlers) GetEmpty(w http.ResponseWriter, r *http.Request) {
	h.HTML(w, http.StatusOK, "")
}

func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("Processing error: %v", err)
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
	}
}
