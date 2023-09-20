package shared

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jlucaspains/sharp-cert-manager/models"
)

func GetConfigSites() []string {
	siteList := []string{}
	for i := 1; true; i++ {
		site, ok := os.LookupEnv(fmt.Sprintf("SITE_%d", i))
		if !ok {
			break
		}

		siteList = append(siteList, site)
	}

	return siteList
}

func CheckCertStatus(params models.CertCheckParams) (result models.CertCheckResult, err error) {
	if params.Url == "" {
		err = errors.New("url is required")
		return
	}

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
		result = models.CertCheckResult{Hostname: hostName, CertStartDate: time.Time{}, CertEndDate: time.Time{}, CertDnsNames: []string{}, IsValid: false}
		return
	}

	certStartDate := resp.TLS.PeerCertificates[0].NotBefore
	certEndDate := resp.TLS.PeerCertificates[0].NotAfter

	isValid, errors := validate(resp.TLS.PeerCertificates[0], hostName)
	result = models.CertCheckResult{
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

	return
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
