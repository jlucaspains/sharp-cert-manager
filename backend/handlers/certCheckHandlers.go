package handlers

import (
	"log"
	"net/http"
	"regexp"

	"github.com/jlucaspains/sharp-cert-manager/models"
	"github.com/jlucaspains/sharp-cert-manager/shared"
)

func (h Handlers) GetSiteList(w http.ResponseWriter, r *http.Request) {
	regx := regexp.MustCompile(`https?:\/\/`)
	siteList := []models.CheckListResult{}
	for _, url := range h.SiteList {
		hostName := regx.ReplaceAllString(url, "")
		siteList = append(siteList, models.CheckListResult{Name: hostName, Url: url})
	}

	h.JSON(w, http.StatusOK, siteList)
}

func (h Handlers) CheckStatus(w http.ResponseWriter, r *http.Request) {
	params := models.CertCheckParams{}

	params.Url, _ = h.getQueryParam(r, "url")

	log.Println("Received message for URL: " + params.Url)

	result, err := shared.CheckCertStatus(params, h.ExpirationWarningDays)

	if err != nil {
		h.JSON(w, http.StatusBadRequest, &models.ErrorResult{Errors: []string{err.Error()}})
		return
	}

	h.JSON(w, http.StatusOK, result)
}
