package client

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/paperhooks/internal/testutil"
	"github.com/jarcoal/httpmock"
	"golang.org/x/oauth2/jws"
)

func TestGCPServiceAccountKeyAuth(t *testing.T) {
	const fakeRsaPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsPnoGUOnrpiSqt4XynxA+HRP7S+BSObI6qJ7fQAVSPtRkqso
tWxQYLEYzNEx5ZSHTGypibVsJylvCfuToDTfMul8b/CZjP2Ob0LdpYrNH6l5hvFE
89FU1nZQF15oVLOpUgA7wGiHuEVawrGfey92UE68mOyUVXGweJIVDdxqdMoPvNNU
l86BU02vlBiESxOuox+dWmuVV7vfYZ79Toh/LUK43YvJh+rhv4nKuF7iHjVjBd9s
B6iDjj70HFldzOQ9r8SRI+9NirupPTkF5AKNe6kUhKJ1luB7S27ZkvB3tSTT3P59
3VVJvnzOjaA1z6Cz+4+eRvcysqhrRgFlwI9TEwIDAQABAoIBAEEYiyDP29vCzx/+
dS3LqnI5BjUuJhXUnc6AWX/PCgVAO+8A+gZRgvct7PtZb0sM6P9ZcLrweomlGezI
FrL0/6xQaa8bBr/ve/a8155OgcjFo6fZEw3Dz7ra5fbSiPmu4/b/kvrg+Br1l77J
aun6uUAs1f5B9wW+vbR7tzbT/mxaUeDiBzKpe15GwcvbJtdIVMa2YErtRjc1/5B2
BGVXyvlJv0SIlcIEMsHgnAFOp1ZgQ08aDzvilLq8XVMOahAhP1O2A3X8hKdXPyrx
IVWE9bS9ptTo+eF6eNl+d7htpKGEZHUxinoQpWEBTv+iOoHsVunkEJ3vjLP3lyI/
fY0NQ1ECgYEA3RBXAjgvIys2gfU3keImF8e/TprLge1I2vbWmV2j6rZCg5r/AS0u
pii5CvJ5/T5vfJPNgPBy8B/yRDs+6PJO1GmnlhOkG9JAIPkv0RBZvR0PMBtbp6nT
Y3yo1lwamBVBfY6rc0sLTzosZh2aGoLzrHNMQFMGaauORzBFpY5lU50CgYEAzPHl
u5DI6Xgep1vr8QvCUuEesCOgJg8Yh1UqVoY/SmQh6MYAv1I9bLGwrb3WW/7kqIoD
fj0aQV5buVZI2loMomtU9KY5SFIsPV+JuUpy7/+VE01ZQM5FdY8wiYCQiVZYju9X
Wz5LxMNoz+gT7pwlLCsC4N+R8aoBk404aF1gum8CgYAJ7VTq7Zj4TFV7Soa/T1eE
k9y8a+kdoYk3BASpCHJ29M5R2KEA7YV9wrBklHTz8VzSTFTbKHEQ5W5csAhoL5Fo
qoHzFFi3Qx7MHESQb9qHyolHEMNx6QdsHUn7rlEnaTTyrXh3ifQtD6C0yTmFXUIS
CW9wKApOrnyKJ9nI0HcuZQKBgQCMtoV6e9VGX4AEfpuHvAAnMYQFgeBiYTkBKltQ
XwozhH63uMMomUmtSG87Sz1TmrXadjAhy8gsG6I0pWaN7QgBuFnzQ/HOkwTm+qKw
AsrZt4zeXNwsH7QXHEJCFnCmqw9QzEoZTrNtHJHpNboBuVnYcoueZEJrP8OnUG3r
UjmopwKBgAqB2KYYMUqAOvYcBnEfLDmyZv9BTVNHbR2lKkMYqv5LlvDaBxVfilE0
2riO4p6BaAdvzXjKeRrGNEKoHNBpOSfYCOM16NjL8hIZB1CaV3WbT5oY+jp7Mzd5
7d56RZOE+ERK2uz/7JX9VSsM/LbH9pJibd4e8mikDS9ntciqOH/3
-----END RSA PRIVATE KEY-----
`
	const fakeIDToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	for _, tc := range []struct {
		name    string
		a       GCPServiceAccountKeyAuth
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: os.ErrInvalid,
		},
		{
			name: "file not found",
			a: GCPServiceAccountKeyAuth{
				KeyFile: filepath.Join(t.TempDir(), "missing"),
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "direct",
			a: GCPServiceAccountKeyAuth{
				Key: []byte(`{
"type": "service_account",
"project_id": "myproject",
"private_key_id": "keyid1234",
"private_key": ` + strconv.QuoteToASCII(fakeRsaPrivateKey) + `,
"client_email": "user@example.com",
"client_id": "clientid1234",
"auth_uri": "https://example.com/o/oauth2/auth",
"token_uri": "https://example.com/token",
"auth_provider_x509_cert_url": "https://example.com/oauth2/v1/certs"
				}`),
			},
		},
		{
			name: "file with audience",
			a: GCPServiceAccountKeyAuth{
				Audience: "testaudience",
				KeyFile: testutil.MustWriteFile(t, filepath.Join(t.TempDir(), "key"), `{
"type": "service_account",
"project_id": "keyfromfile",
"private_key_id": "keyid10250",
"private_key": `+strconv.QuoteToASCII(fakeRsaPrivateKey)+`,
"client_email": "file@example.com",
"client_id": "clientid18686",
"auth_uri": "https://example.com/o/oauth2/auth",
"token_uri": "https://example.com/token",
"auth_provider_x509_cert_url": "https://example.com/oauth2/v1/certs"
}`),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)
			transport.RegisterMatcherResponder(http.MethodPost, "https://example.com/token",
				httpmock.NewMatcher("", func(req *http.Request) bool {
					if err := req.ParseForm(); err != nil {
						t.Errorf("ParseForm() failed: %v", err)
					}

					if _, err := jws.Decode(req.Form.Get("assertion")); err != nil {
						t.Errorf("Decode() failed: %v", err)
					}

					return true
				}),
				httpmock.NewStringResponder(http.StatusOK, `{ "id_token": "`+fakeIDToken+`" }`))
			transport.RegisterMatcherResponder(http.MethodGet, "http://localhost/",
				httpmock.HeaderIs("Authorization", "Bearer "+fakeIDToken),
				httpmock.NewStringResponder(http.StatusOK, "success:"+t.Name()))

			tc.a.HTTPClient = &http.Client{
				Transport: transport,
			}

			got, err := tc.a.Build()

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Build() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				r := resty.New().
					SetBaseURL("http://localhost").
					SetTransport(transport)

				got.authenticate(Options{}, r)

				for range [3]struct{}{} {
					if resp, err := r.R().Get("/"); err != nil {
						t.Errorf("Get() failed: %v", err)
					} else if diff := cmp.Diff("success:"+t.Name(), string(resp.Body())); diff != "" {
						t.Errorf("Response body diff (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
