package handlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-checker/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetSiteList(t *testing.T) {
	godotenv.Load("../.test.env")
	handlers := new(Handlers)
	handlers.SiteList = handlers.GetConfigSites()
	router := mux.NewRouter()
	router.HandleFunc("/site-list", handlers.GetSiteList).Methods("GET")

	code, body, err := makeRequest[[]models.CheckListResult](router, "GET", "/site-list", nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(*body))
}

func TestGetCheckStatus(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/check-status", handlers.CheckStatus).Methods("GET")

	url := fmt.Sprintf("/check-status?url=%s", "https://blog.lpains.net")
	code, body, err := makeRequest[models.CertCheckResult](router, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, body.IsValid)
	assert.LessOrEqual(t, body.CertStartDate, time.Now())
	assert.GreaterOrEqual(t, body.CertEndDate, time.Now())
	assert.Contains(t, body.Hostname, "blog.lpains.net")
	assert.Contains(t, body.CertDnsNames, "blog.lpains.net")
}

func TestGetCheckStatusNoUrl(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/check-status", handlers.CheckStatus).Methods("GET")

	url := "/check-status"
	code, body, err := makeRequest[models.ErrorResult](router, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 400, code)
	assert.Equal(t, "Url is required", body.Errors[0])
}

func TestGetCheckStatusHttp(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/check-status", handlers.CheckStatus).Methods("GET")

	url := fmt.Sprintf("/check-status?url=%s", "http://blog.lpains.net")
	code, body, err := makeRequest[models.CertCheckResult](router, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.False(t, body.IsValid)
}

func TestGetCheckStatusAllValidations(t *testing.T) {
	handlers := new(Handlers)

	router := mux.NewRouter()
	router.HandleFunc("/check-status", handlers.CheckStatus).Methods("GET")

	mux := http.NewServeMux()
	ts := httptest.NewUnstartedServer(mux)

	certPEM := `-----BEGIN CERTIFICATE-----
MIIDkzCCAnugAwIBAgIUOgggUW2hRYhzI14GIwabMHA4ZWcwDQYJKoZIhvcNAQEF
BQAwWTELMAkGA1UEBhMCVVMxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDESMBAGA1UEAwwJbG9jYWxob3N0MB4X
DTIzMDkwMjIyNTMyMVoXDTIzMDkwMzIyNTMyMVowWTELMAkGA1UEBhMCVVMxEzAR
BgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5
IEx0ZDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAwyLJeolEXCuTndKT9reEBtmvZ5nTfOEuAXfXiSNW1u0hue5ivHu+
oJHkGeRCmAYbrwOhv2SGJBG5BVE8WhBC4IXtR0rKsWVjJrKX6mYCKf12+AlC2bvK
tP7tdq1R6nmARpoTsDcoz7h/jFqXu07oru6W2XNfx1kwDrvZcQB+p9TdXn/kimBx
CGsXKZdWkY2Fcso3rZpNUW22B9cVbQjKxPlt+cm1cYXPUDZCFvF0aw1PPQA4GSSH
2PoHycQdDdA0jyKwypfcdsKgB+TfQnVDFYS7j6y4zGg3wwD+5Cj3Kf+CjUurmPq+
NXUB2i1+gL3Ve72Sf/lzvf1CCVWiHU0wrwIDAQABo1MwUTAdBgNVHQ4EFgQU/l5D
jjUrhbBbjb1wVLrUTkWI4IgwHwYDVR0jBBgwFoAU/l5DjjUrhbBbjb1wVLrUTkWI
4IgwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEAelw6SRb6DvM6
2oEBdk9mTEjHbMRNXmhLBp4UZNIHNSVKrVzNfGUNgGnEdnGcxC2728A4H71Z83dS
o23pM2p3TJVv4Xj2cHAH/XD6vO7jc65UIq1/1F/+QO/8otWWreeM/UM1K5YxyxMp
IkuquUWxZjGtzVeI/3wituLg23Sb4ibAaHcaU+JrD0ySmXZd0mgtslVd+BT6/4HL
S2sqiJP3bhYWHcx3BMe/K3LLLr7D4NiaSeZmcqhotFusvqIedMrxBQ4hvgTJOaCf
HCHBbC/PBypgqvRkCWZTJypMRLph7TOTsH+qQh2OKUr30w4udASYt8poTtvAB6Ih
7hGcF9509g==
-----END CERTIFICATE-----
`

	keyPEM := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDDIsl6iURcK5Od
0pP2t4QG2a9nmdN84S4Bd9eJI1bW7SG57mK8e76gkeQZ5EKYBhuvA6G/ZIYkEbkF
UTxaEELghe1HSsqxZWMmspfqZgIp/Xb4CULZu8q0/u12rVHqeYBGmhOwNyjPuH+M
Wpe7Tuiu7pbZc1/HWTAOu9lxAH6n1N1ef+SKYHEIaxcpl1aRjYVyyjetmk1RbbYH
1xVtCMrE+W35ybVxhc9QNkIW8XRrDU89ADgZJIfY+gfJxB0N0DSPIrDKl9x2wqAH
5N9CdUMVhLuPrLjMaDfDAP7kKPcp/4KNS6uY+r41dQHaLX6AvdV7vZJ/+XO9/UIJ
VaIdTTCvAgMBAAECggEAZQd+vxWQshPRsrWS6/qpvY45FW98IrzHP/VXP2ZvkIln
m8dDkYiT8rh2G72liOYosR01Qk1+cfBHFeywTYT7yxkr92xOszfl9OQkaTR1TF4x
mUvaM7bZxYnzUi18KuTLOEKPjP2SALHqP7Wrt0usht16do0YerK/gfFaK/pwmN1g
qUxDTygplogqApJ5y6b2JnWsHgSEaaHKD3YgmpiDeWd83L6m2ikMSoSeEa7j1GS2
kKDHNAIgxXSi00VAB0xQxke/kOOZa+UmBpgLKyZ8Mntm1IfodP+L7ZjXQY9znJA5
ALCydcJY+V6S0xahjhNZkvD7jq6wdVnlDDhiEXMF+QKBgQD8zbiaqGqyYNLR3bcY
JR7Z2EvIAHtG4Z4VusAc2oyiBp9F2t0IK8oPWch2x5bVN/wrGwxneOWKPjNEBgu6
5RDkik0BrXsLjz8kkTds90ughD197Qa0FKvaYnWQiwe7N7HJLudl0wlmUic5YJLy
aQhX4Yn+yrgOYfixie/BMGLfMwKBgQDFmmf2yCskA0XJpd2hbpeItqjg5QuhyA9l
gNKOEMCqKpbMmBD2Y8UKotGm1iX4jppdaNrpPQuoc1B/wdGDPRfEMCMxL/7p+MOb
AKHcC+kI1DP4HHmiBqGFbGYvEUqKTF+vDkTc3ELxBnkkI26WOjYw8C86Oup43LxV
9v5gv9mYlQKBgQCcjKis3W51WBA1dh9UDGi2boM/L00n7799pVAijhRYodEv6QDH
dpaCOw8wvxhgoXK/Htjnmq5KlYoZrcTFz+ROInbdexifZ+2qL2MrT1i95iZOPOHR
0ps5eY9kGzSGc07dTvZsz+saOfWgSnW1N+W6xig2aELiZTkkeE7IS7ZukQKBgEG5
V7cHYQH7bKzjVFIrXI+GYalbxYCr8CMMs/u4qrxuqfWm5o1tJc6h1SWuuLZxh/pl
s9o8CbKfmDjGGI+UNGF2uV3U3u6nZTga/7sW4w2itx5hKjuwBO1B3sLs92QEfxbU
oibrxAAy7PwOJOwmtHuWh77Qdch5ctMM8hLv/Mn5AoGAbvelLdvk65eBvXSvUDzj
BEKeRm/jPFSfYnWB0eblZ1isA2brhUtd9yP8Kv4YSzYiljFlFW/58nDEmXhyWl8k
VXxOULagYbVOT+gxRDv4eCniZVL4g82k+0NmvzfNVTJgmXWJnXvvFd2fSFEhhfq1
yjbTOuy8KoxNb15g3Ysesbw=
-----END PRIVATE KEY-----
`

	cert, _ := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	ts.StartTLS()
	defer ts.Close()

	url := fmt.Sprintf("/check-status?url=%s", ts.URL)
	code, body, err := makeRequest[models.CertCheckResult](router, "GET", url, nil)

	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.False(t, body.IsValid)
	assert.Equal(t, []string{"Hostname is not valid", "Certificate is not valid yet or expired", "SHA1 is not a secure signature algorithm"}, body.ValidationIssues)
}
