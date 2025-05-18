package client

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

// Fake X.509 certificate for testing. Generated using the following commands:
//
//	openssl genrsa -out key 512
//
//	faketime '2000-01-01 00:00 UTC' \
//	 openssl req -x509 -new -nodes -key key -days 1 -out cert \
//	   -outform PEM -batch
const fakeCertPEM = `
-----BEGIN CERTIFICATE-----
MIIB4TCCAYugAwIBAgIUDUm2YVOrpISBpfQO5H8o3Kxu8S0wDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0wMDAxMDEwMDAwMDBaFw0wMDAx
MDIwMDAwMDBaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEF
AANLADBIAkEAxMlvlAar74MFUhb9LrqeclDmKWsjWbuiCVdAoj8+Gq+XG3B4H3bL
auNZ+dhyr3eZuHsbw+D3KToeiMJRxAsRZQIDAQABo1MwUTAdBgNVHQ4EFgQURv88
YquEePhNH7s0kP5Gu8rDX0QwHwYDVR0jBBgwFoAURv88YquEePhNH7s0kP5Gu8rD
X0QwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAANBACPTaHtvgMR9yzrb
YhfkbRMvlye1i/xliJihG6kUSd9oPEqtTN6L/6qId2FWENfxijFIceavp6VLYsun
cr9Jj64=
-----END CERTIFICATE-----
`

func newFakeCertPool(t *testing.T) *x509.CertPool {
	t.Helper()

	block, _ := pem.Decode([]byte(fakeCertPEM))

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Parsing test certificate: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(cert)

	return pool
}
