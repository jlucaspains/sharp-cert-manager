package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"regexp"
	"strings"
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

	hostName := strings.Split(params.Url, ":")[1]
	hostName = strings.Trim(hostName, "/")

	if resp.TLS == nil {
		result := &models.CertCheckResult{Hostname: hostName, CertStartDate: time.Time{}, CertEndDate: time.Time{}, CertDnsNames: []string{}, IsValid: false}
		c.JSON(result)
		return nil
	}

	certStartDate := resp.TLS.PeerCertificates[0].NotBefore
	certEndDate := resp.TLS.PeerCertificates[0].NotAfter

	isValid, errors := validate(resp.TLS.PeerCertificates[0], hostName)
	result := &models.CertCheckResult{
		Hostname:         hostName,
		Issuer:           resp.TLS.PeerCertificates[0].Issuer.CommonName,
		Signature:        resp.TLS.PeerCertificates[0].SignatureAlgorithm.String(),
		CertStartDate:    certStartDate,
		CertEndDate:      certEndDate,
		CertDnsNames:     resp.TLS.PeerCertificates[0].DNSNames,
		TLSVersion:       resp.TLS.Version,
		IsCA:             resp.TLS.PeerCertificates[0].IsCA,
		CommonName:       resp.TLS.PeerCertificates[0].Subject.CommonName,
		IsValid:          isValid,
		OtherCerts:       getOtherCerts(resp.TLS.PeerCertificates[1:]),
		ValidationIssues: errors,
	}

	c.JSON(result)

	return nil
}

func validate(cert *x509.Certificate, hostName string) (isValid bool, errors []string) {
	isHostNameValid := cert.VerifyHostname(hostName) == nil
	areDatesValid := cert.NotBefore.Before(time.Now().UTC()) && cert.NotAfter.After(time.Now().UTC())
	isSignatureValid := cert.SignatureAlgorithm.String() != "SHA1-RSA"
	isValid = isHostNameValid && areDatesValid && isSignatureValid

	if !isHostNameValid {
		errors = append(errors, "Hostname is not valid")
	}
	if !areDatesValid {
		errors = append(errors, "Certificate is not valid yet or expired")
	}
	if !isSignatureValid {
		errors = append(errors, "SHA1 is not a secure signature algorithm")
	}

	return
}

func getOtherCerts(certs []*x509.Certificate) []models.OtherCert {
	results := []models.OtherCert{}
	for _, v := range certs {
		results = append(results, models.OtherCert{
			IsCA:       v.IsCA,
			CommonName: v.Subject.CommonName,
			Issuer:     v.Issuer.CommonName,
		})
	}

	return results
}
