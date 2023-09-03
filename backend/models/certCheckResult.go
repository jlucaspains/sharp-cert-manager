package models

type CertCheckResult struct {
	Hostname      string   `json:"hostname"`
	Issuer        string   `json:"issuer"`
	Signature     string   `json:"signature"`
	CertStartDate string   `json:"certStartDate"`
	CertEndDate   string   `json:"certEndDate"`
	CertDnsNames  []string `json:"certDnsNames"`
	IsValid       bool     `json:"isValid"`
}

type CheckListResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
