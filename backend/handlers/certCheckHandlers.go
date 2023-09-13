package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jlucaspains/sharp-cert-manager/models"
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
	params := &models.CertCheckParams{}

	params.Url, _ = h.getQueryParam(r, "url")

	if params.Url == "" {
		h.JSON(w, http.StatusBadRequest, &models.ErrorResult{Errors: []string{"Url is required"}})
		return
	}

	log.Println("Received message for URL: " + params.Url)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		// disabling security here is fine
		// the purpose of the client is to pull certs
		// no data exchange is happening
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, _ := client.Get(params.Url)

	hostName := strings.Split(params.Url, ":")[1]
	hostName = strings.Trim(hostName, "/")

	if resp.TLS == nil {
		result := &models.CertCheckResult{Hostname: hostName, CertStartDate: time.Time{}, CertEndDate: time.Time{}, CertDnsNames: []string{}, IsValid: false}
		h.JSON(w, http.StatusOK, result)
		return
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

	h.JSON(w, http.StatusOK, result)
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
