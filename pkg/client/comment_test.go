package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestListComments(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		docID   int64
		wantErr error
		want    []Comment
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/123/comments/",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
			docID: 123,
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/0/comments/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "entries",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/18681/comments/",
					httpmock.NewStringResponder(http.StatusOK, `[
						{
							"id": 603,
							"comment": "bar",
							"created": "2023-03-01T22:23:06.183451Z",
							"user": {
								"id": 5,
								"username": "johndoe",
								"firstname": "",
								"lastname": ""
							}
						},
						{
							"id": 485,
							"comment": "foo",
							"created": "2023-03-01T22:23:04.391361Z",
							"user": {
								"id": 2,
								"username": "admin",
								"firstname": "",
								"lastname": ""
							}
						}
					]`))
			},
			docID: 18681,
			want: []Comment{
				{
					ID:      603,
					Text:    "bar",
					Created: time.Date(2023, time.March, 1, 22, 23, 6, 183451000, time.UTC),
					User:    CommentUser{ID: 5, Name: "johndoe"},
				},
				{
					ID:      485,
					Text:    "foo",
					Created: time.Date(2023, time.March, 1, 22, 23, 4, 391361000, time.UTC),
					User:    CommentUser{ID: 2, Name: "admin"},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.ListComments(context.Background(), tc.docID)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("ListComments() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("ListComments() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		docID   int64
		input   *Comment
		wantErr error
	}{
		{
			name:  "empty",
			docID: 7972,
			input: &Comment{},
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPost, "/api/documents/7972/comments/",
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
		},
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodPost, "/api/documents/1457/comments/",
					httpmock.BodyContainsString(`my+text`),
					httpmock.NewStringResponder(http.StatusOK, `[]`))
			},
			docID: 1457,
			input: &Comment{
				Text: "my text",
			},
		},
		{
			name: "error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPost, "/api/documents/24597/comments/",
					httpmock.NewStringResponder(http.StatusTeapot, `{}`))
			},
			docID: 24597,
			input: &Comment{},
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
				Message:    "{}",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			_, err := c.CreateComment(context.Background(), tc.docID, tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("CreateComment() error diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		docID   int64
		id      int64
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodDelete, "/api/documents/9234/comments/",
					"id=7816",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			docID: 9234,
			id:    7816,
		},
		{
			name: "error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponderWithQuery(http.MethodDelete, "/api/documents/7146/comments/",
					"id=12234",
					httpmock.NewStringResponder(http.StatusTeapot, `{ "detail": "error" }`))
			},
			docID: 7146,
			id:    12234,
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
				Message:    `{"detail":"error"}`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			_, err := c.DeleteComment(context.Background(), tc.docID, tc.id)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteComment() error diff (-want +got):\n%s", diff)
			}
		})
	}
}
