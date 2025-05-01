package client

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestDownload(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		url     string
		want    *DownloadResult
		wantErr error
	}{
		{
			name: "header",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/test/download?foo=bar",
					httpmock.NewStringResponder(http.StatusOK, `content`).
						HeaderSet(http.Header{
							"Content-Type":        []string{"foo/bar; charset=utf-8"},
							"Content-Disposition": []string{`inline; filename="test.txt"`},
						}))
			},
			url: "/test/download?foo=bar",
			want: &DownloadResult{
				ContentType: "foo/bar",
				ContentTypeParams: map[string]string{
					"charset": "utf-8",
				},
				Filename: "test.txt",
				Length:   7,
			},
		},
		{
			name: "error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/report/error",
					httpmock.NewStringResponder(http.StatusTeapot, ``))
			},
			url: "/report/error",
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
				Message:    `418 I'm a teapot`,
			},
		},
		{
			name: "connection error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/conn/error",
					httpmock.ConnectionFailure)
			},
			url:     "/conn/error",
			wantErr: cmpopts.AnyError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			var buf bytes.Buffer

			got, _, err := c.download(context.Background(), &buf, tc.url, true)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DownloadDocument() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("DownloadDocument() result diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
