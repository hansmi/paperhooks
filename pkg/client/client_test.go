package client

import (
	"context"
	"crypto/x509"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func newMockTransport(t *testing.T) *httpmock.MockTransport {
	t.Helper()

	transport := httpmock.NewMockTransport()
	transport.RegisterNoResponder(httpmock.NewNotFoundResponder(t.Fatal))

	return transport
}

func TestClient(t *testing.T) {
	for _, tc := range []struct {
		name    string
		opts    Options
		setup   func(*testing.T, *httpmock.MockTransport)
		wantErr error
	}{
		{
			name: "defaults",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodGet, "/api/",
					httpmock.HeaderIs("Accept", "application/json; version=2"),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, nil))
			},
		},
		{
			name: "customized",
			opts: Options{
				BaseURL: "http://localhost:1234/path/",
				Auth:    &TokenAuth{"foo26175bar"},
				Header: http.Header{
					"X-Custom": []string{"aaa"},
				},
			},
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodGet, "http://localhost:1234/path/api/",
					httpmock.HeaderIs("Authorization", "Token foo26175bar").And(
						httpmock.HeaderIs("X-Custom", "aaa"),
					),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, nil))
			},
		},
		{
			name: "basic auth",
			opts: Options{
				BaseURL: "http://localhost:1234////",
				Auth: &UsernamePasswordAuth{
					Username: "user",
					Password: "password",
				},
			},
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodGet, "http://localhost:1234/api/",
					httpmock.HeaderIs("Authorization", "Basic dXNlcjpwYXNzd29yZA=="),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, nil))
			},
		},
		{
			name: "redirect",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/",
					httpmock.NewJsonResponderOrPanic(http.StatusSeeOther, nil))
			},
			wantErr: &RequestError{
				StatusCode: http.StatusSeeOther,
				Message:    "303 See Other",
			},
		},
		{
			name: "internal server error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/",
					httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, nil))
			},
			wantErr: &RequestError{
				StatusCode: http.StatusInternalServerError,
				Message:    "null",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			tc.opts.transport = transport

			err := New(tc.opts).Ping(context.Background())

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Ping() error diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClientWithTLS(t *testing.T) {
	for _, tc := range []struct {
		name  string
		trust bool
		opts  Options
	}{
		{
			name:  "trusted",
			trust: true,
		},
		{
			name: "untrusted",
		},
		{
			name:  "max concurrent",
			trust: true,
			opts: Options{
				MaxConcurrentRequests: 100,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}))
			t.Cleanup(srv.Close)

			srv.Config.ErrorLog = log.New(io.Discard, "", 0)
			srv.StartTLS()

			opts := tc.opts
			opts.BaseURL = srv.URL
			opts.TrustedRootCAs = x509.NewCertPool()

			if tc.trust {
				opts.TrustedRootCAs.AddCert(srv.Certificate())
			}

			err := New(opts).Ping(t.Context())

			if !tc.trust {
				var caErr x509.UnknownAuthorityError

				if err == nil || !errors.As(err, &caErr) {
					t.Errorf("Ping() should report bad X.509 CA, got: %v", err)
				}
			} else if err != nil {
				t.Errorf("Ping() failed: %v", err)
			}
		})
	}
}

func TestWrapResponse(t *testing.T) {
	for _, tc := range []struct {
		name string
		resp *resty.Response
	}{
		{
			name: "nil",
		},
		{
			name: "response",
			resp: &resty.Response{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			wrapResponse(tc.resp)
		})
	}
}
