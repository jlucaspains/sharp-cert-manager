package shared

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jlucaspains/sharp-cert-manager/models"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

func GetConfigCerts() (result []models.CheckCertItem) {
	for i := 1; true; i++ {
		rawUrl, ok := os.LookupEnv(fmt.Sprintf("SITE_%d", i))
		if !ok {
			break
		}
		site, err := url.Parse(rawUrl)

		if err != nil {
			continue
		}

		result = append(result, models.CheckCertItem{Name: site.Hostname(), Url: rawUrl, Type: models.CertCheckURL})
	}

	for i := 1; true; i++ {
		rawUrl, ok := os.LookupEnv(fmt.Sprintf("AZUREKEYVAULT_%d", i))
		if !ok {
			break
		}
		akvUrl, err := url.Parse(rawUrl)

		if err != nil {
			continue
		}

		name := akvUrl.Hostname() + "/" + strings.Split(akvUrl.Path, "/")[2]

		result = append(result, models.CheckCertItem{Name: name, Url: rawUrl, Type: models.CertCheckAzure})
	}

	return
}

func CheckCertStatus(cert models.CheckCertItem, expirationWarningDays int) (result *models.CertCheckResult, err error) {
	if cert.Name == "" || cert.Url == "" {
		err = errors.New("name, url, and type are required")
		return
	}

	switch cert.Type {
	case models.CertCheckURL:
		result, err = checkCertByUrlStatus(cert.Name, cert.Url, expirationWarningDays)
	case models.CertCheckAzure:
		result, err = checkAzureCertStatus(cert.Name, cert.Url, expirationWarningDays)
	}

	return
}

func checkCertByUrlStatus(name string, url string, expirationWarningDays int) (result *models.CertCheckResult, err error) {
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
	resp, err := client.Get(url)

	if err != nil || resp.TLS == nil {
		result = &models.CertCheckResult{Hostname: name, CertStartDate: time.Time{}, CertEndDate: time.Time{}, CertDnsNames: []string{}, IsValid: false}
		return
	}

	certStartDate := resp.TLS.PeerCertificates[0].NotBefore
	certEndDate := resp.TLS.PeerCertificates[0].NotAfter

	isValid, errors := validate(resp.TLS.PeerCertificates[0], name, false)
	result = &models.CertCheckResult{
		Hostname:          name,
		Issuer:            resp.TLS.PeerCertificates[0].Issuer.CommonName,
		Signature:         resp.TLS.PeerCertificates[0].SignatureAlgorithm.String(),
		CertStartDate:     certStartDate,
		CertEndDate:       certEndDate,
		CertDnsNames:      resp.TLS.PeerCertificates[0].DNSNames,
		TLSVersion:        resp.TLS.Version,
		IsCA:              resp.TLS.PeerCertificates[0].IsCA,
		CommonName:        resp.TLS.PeerCertificates[0].Subject.CommonName,
		IsValid:           isValid,
		OtherCerts:        getOtherCerts(resp.TLS.PeerCertificates[1:]),
		ValidationIssues:  errors,
		ExpirationWarning: certEndDate.Before(time.Now().AddDate(0, 0, expirationWarningDays)),
	}

	return
}

func checkAzureCertStatus(name string, rawUrl string, expirationWarningDays int) (result *models.CertCheckResult, err error) {
	parsedUrl, _ := url.Parse(rawUrl)
	keyVaultUrl := parsedUrl.Scheme + "://" + parsedUrl.Host
	certName := strings.Split(parsedUrl.Path, "/")[2]

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return
	}

	client, err := azcertificates.NewClient(keyVaultUrl, cred, nil)

	if err != nil {
		return
	}

	log.Printf("Getting certificate from Azure Key Vault: %s", certName)

	response, err := client.GetCertificate(context.Background(), certName, "", nil)

	if err != nil {
		return
	}

	cert, err := x509.ParseCertificate(response.CER)

	if err != nil {
		return
	}

	certStartDate := cert.NotBefore
	certEndDate := cert.NotAfter

	isValid, errors := validate(cert, name, true)
	result = &models.CertCheckResult{
		Hostname:          name,
		Issuer:            cert.Issuer.CommonName,
		Signature:         cert.SignatureAlgorithm.String(),
		CertStartDate:     certStartDate,
		CertEndDate:       certEndDate,
		CertDnsNames:      cert.DNSNames,
		TLSVersion:        uint16(cert.Version),
		IsCA:              cert.IsCA,
		CommonName:        cert.Subject.CommonName,
		IsValid:           isValid,
		OtherCerts:        []models.OtherCert{},
		ValidationIssues:  errors,
		ExpirationWarning: certEndDate.Before(time.Now().AddDate(0, 0, expirationWarningDays)),
	}

	return
}

func validate(cert *x509.Certificate, hostName string, skipHostNameValidation bool) (isValid bool, errors []string) {
	isHostNameValid := skipHostNameValidation || cert.VerifyHostname(hostName) == nil
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
