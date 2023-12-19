package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestGetCurrentUser(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		wantErr error
		want    *User
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/ui_settings/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/ui_settings/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/ui_settings/",
					httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]any{
						"user": map[string]any{"id": 123},
					}))
				transport.RegisterResponder(http.MethodGet, "/api/users/123/",
					httpmock.NewJsonResponderOrPanic(http.StatusOK, map[string]any{
						"id":       456,
						"username": "testuser",
					}))
			},
			want: &User{
				ID:       456,
				Username: "testuser",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.GetCurrentUser(context.Background())

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetCurrentUser() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetCurrentUser() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
