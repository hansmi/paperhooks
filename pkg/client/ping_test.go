package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestClientPing(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name      string
		responder httpmock.Responder
		wantErr   error
	}{
		{
			name:      "success",
			responder: httpmock.NewJsonResponderOrPanic(http.StatusOK, nil),
		},
		{
			name:      "not found",
			responder: httpmock.NewJsonResponderOrPanic(http.StatusNotFound, nil),
			wantErr: &RequestError{
				StatusCode: http.StatusNotFound,
				Message:    "null",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)
			transport.RegisterResponder(http.MethodGet, "/api/", tc.responder)

			c := New(Options{
				transport: transport,
			})

			err := c.Ping(ctx)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Ping() error diff (-want +got):\n%s", diff)
			}
		})
	}
}
