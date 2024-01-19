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
MIIBOwIBAAJBAK0uWi2bu8rMLcv+NCs4J4dW0SHFQ6wax6YYQX9SO3YkJtyhNnB+
9r7G0Ei4EVnViXH/WbgoCdgIIfKIP6yJYYsCAwEAAQJAAhsnM5jKPtweznVH8yKa
sHWo021Ptl8ZAHcZDNBWMsiWpS0T1AduvKqWm03eVznRXkReTSLO2y/68H71kSkI
iQIhANYHvzIjUW1SgJ5CcXKASICmlaic/t7hSmvYiWE+mvtdAiEAzyP8/yDEk/2s
5KYVROHZ7r/vIV2dckXPVjrKfgeYagcCIC4yqeBmozLXtg9zBA3VBtFOI8urZ5Aw
TOIOcUjePJG5AiEAqBZfDbTMb/7xFpYDOmM/krLjXKL3yawGhMWuXbjSIG8CIQDU
ZeJO3dIEZiy84+1LTckzPluFpUGzGmsCbYXRBBYAig==
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
