package client

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestGetDocumentMetadata(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		id      int64
		want    *DocumentMetadata
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/7124/metadata/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"original_mime_type": "text/plain"
					}`))
			},
			id: 7124,
			want: &DocumentMetadata{
				OriginalMimeType: "text/plain",
			},
		},
		{
			name: "error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/25650/metadata/",
					httpmock.NewStringResponder(http.StatusTeapot, `{ "detail": "wrong" }`))
			},
			id: 25650,
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
				Message:    `{"detail":"wrong"}`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.GetDocumentMetadata(context.Background(), tc.id)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetDocumentMetadata() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetDocumentMetadata() result diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestUploadDocument(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		r       io.Reader
		opts    DocumentUploadOptions
		want    *DocumentUpload
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPost, "/api/documents/post_document/",
					httpmock.NewStringResponder(http.StatusOK, `"e068eb08-cf70-4755-8087-3cf0644f3c7b"`))
			},
			r: strings.NewReader("test content"),
			want: &DocumentUpload{
				TaskID: "e068eb08-cf70-4755-8087-3cf0644f3c7b",
			},
		},
		{
			name: "options",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterMatcherResponder(http.MethodPost, "/api/documents/post_document/",
					httpmock.BodyContainsString("\ndoctitle"),
					httpmock.NewStringResponder(http.StatusOK, `"0dbf0a2b-3a09-4d7b-96bf-51544dda8427"`))
			},
			r: strings.NewReader("more content"),
			opts: DocumentUploadOptions{
				Filename:            filepath.Join(t.TempDir(), "myfile.txt"),
				Title:               "doctitle",
				Created:             time.Date(2020, time.December, 31, 1, 2, 3, 0, time.UTC),
				Correspondent:       Int64(100),
				DocumentType:        Int64(200),
				Tags:                []int64{300, 301, 302},
				ArchiveSerialNumber: Int64(400),
			},
			want: &DocumentUpload{
				TaskID: "0dbf0a2b-3a09-4d7b-96bf-51544dda8427",
			},
		},
		{
			name: "error",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodPost, "/api/documents/post_document/",
					httpmock.NewStringResponder(http.StatusTeapot, `{}`))
			},
			r: strings.NewReader(""),
			wantErr: &RequestError{
				StatusCode: http.StatusTeapot,
				Message:    `{}`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.UploadDocument(context.Background(), tc.r, tc.opts)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("UploadDocument() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("UploadDocument() result diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
