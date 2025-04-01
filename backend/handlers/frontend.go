package handlers

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"slices"
	"time"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/jlucaspains/sharp-cert-manager/shared"
)

var indexTemplate *template.Template

func initTemplates() {
	if indexTemplate != nil {
		return
	}

	indexTemplate = template.Must(template.ParseGlob("frontend/*"))
}

func (h Handlers) Index(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	err := indexTemplate.ExecuteTemplate(w, "index.html", h.CertList)

	if err != nil {
		log.Fatal(err)
	}
}

type FragmentResult struct {
	Name      string
	Issuer    string
	Signature string
	Validity  int
	IsValid   bool
}

func (h Handlers) Fragment(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	name, _ := h.getQueryParam(r, "name")

	log.Println("Received message for name: " + name)

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
		h.HTML(w, http.StatusBadRequest, "Failed to process request")
		return
	}

	validity := time.Now().UTC().Sub(result.CertStartDate)
	// create object to be passed to template
	parsedResult := &FragmentResult{
		Name:      item.Name,
		Issuer:    result.Issuer,
		Signature: result.Signature,
		IsValid:   result.IsValid,
		Validity:  int(math.Abs(float64(validity.Hours() / 24))),
	}

	err = indexTemplate.ExecuteTemplate(w, "itemLoaded.html", parsedResult)

	if err != nil {
		h.HTML(w, http.StatusBadRequest, "Failed to process request")
		return
	}
}
