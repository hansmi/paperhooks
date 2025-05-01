package client

import (
	"context"
	"net/http"
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
