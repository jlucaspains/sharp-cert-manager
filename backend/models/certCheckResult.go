package models

import "time"

type OtherCert struct {
	CommonName string `json:"commonName"`
	Issuer     string `json:"issuer"`
	IsCA       bool   `json:"isCA"`
}

type CertCheckResult struct {
	Hostname          string      `json:"hostname"`
	Issuer            string      `json:"issuer"`
	Signature         string      `json:"signature"`
	CertStartDate     time.Time   `json:"certStartDate"`
	CertEndDate       time.Time   `json:"certEndDate"`
	CertDnsNames      []string    `json:"certDnsNames"`
	IsValid           bool        `json:"isValid"`
	TLSVersion        uint16      `json:"tlsVersion"`
	IsCA              bool        `json:"isCA"`
	CommonName        string      `json:"commonName"`
	OtherCerts        []OtherCert `json:"otherCerts"`
	ValidationIssues  []string    `json:"validationIssues"`
	ExpirationWarning bool        `json:"expirationWarning"`
}

type CheckListResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
