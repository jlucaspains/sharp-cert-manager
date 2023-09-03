package handlers

import (
	"crypto/tls"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jlucaspains/sharp-cert-checker/models"
)

func (h Handlers) GetSiteList(c *fiber.Ctx) error {
	regx := regexp.MustCompile(`https?:\/\/`)
	siteList := []models.CheckListResult{}
	for _, url := range h.SiteList {
		hostName := regx.ReplaceAllString(url, "")
		siteList = append(siteList, models.CheckListResult{Name: hostName, Url: url})
	}

	c.JSON(siteList)

	return nil
}

func (h Handlers) CheckStatus(c *fiber.Ctx) error {
	params := &models.CertCheckParams{}
	params.Url = c.Query("url")

	log.Println("Received message for URL: " + params.Url)

	if params.Url == "" {
		c.Status(http.StatusBadRequest)
		c.JSON(&fiber.Map{"error": "Missing URL parameter"})
		return nil
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(params.Url)

	if err != nil {
		return err
	}

	regx := regexp.MustCompile(`https?:\/\/`)
	hostName := regx.ReplaceAllString(params.Url, "")

	if resp.TLS == nil {
		result := &models.CertCheckResult{Hostname: hostName, CertStartDate: "", CertEndDate: "", CertDnsNames: []string{}, IsValid: false}
		c.JSON(result)
		return nil
	}

	certStartDate := resp.TLS.PeerCertificates[0].NotBefore
	certEndDate := resp.TLS.PeerCertificates[0].NotAfter

	isValid := resp.TLS.PeerCertificates[0].VerifyHostname(hostName) == nil && certStartDate.Before(time.Now().UTC()) && certEndDate.After(time.Now().UTC())
	result := &models.CertCheckResult{
		Hostname:      hostName,
		Issuer:        resp.TLS.PeerCertificates[0].Issuer.CommonName,
		Signature:     resp.TLS.PeerCertificates[0].SignatureAlgorithm.String(),
		CertStartDate: certStartDate.String(),
		CertEndDate:   certEndDate.String(),
		CertDnsNames:  resp.TLS.PeerCertificates[0].DNSNames,
		IsValid:       isValid,
	}

	c.JSON(result)

	return nil
}
