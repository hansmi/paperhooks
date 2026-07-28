package client

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/paperhooks/internal/testutil"
	"github.com/jarcoal/httpmock"
)

func TestGetDocument(t *testing.T) {
	plus2 := time.FixedZone("UTC+2", 2*60*60)

	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		loc     *time.Location
		id      int64
		want    *Document
		wantErr error
	}{
		{
			name: "created with datetime",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/8127/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"id": 8127,
						"title": "first",
						"created": "2023-06-30T22:00:00Z"
					}`))
			},
			id: 8127,
			want: &Document{
				ID:      8127,
				Title:   "first",
				Created: time.Date(2023, time.June, 30, 22, 0, 0, 0, time.UTC),
			},
		},
		{
			// Paperless REST API version 9 changed the "created" field to
			// a date-only string.
			name: "created with date only",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/8128/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"id": 8128,
						"title": "second",
						"created": "2026-07-23",
						"modified": "2026-07-23T08:09:10Z"
					}`))
			},
			loc: time.UTC,
			id:  8128,
			want: &Document{
				ID:       8128,
				Title:    "second",
				Created:  time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC),
				Modified: time.Date(2026, time.July, 23, 8, 9, 10, 0, time.UTC),
			},
		},
		{
			// Date-only values are interpreted as midnight in the server's
			// timezone.
			name: "created with date only in server timezone",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/8130/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"id": 8130,
						"created": "2026-07-23"
					}`))
			},
			loc: plus2,
			id:  8130,
			want: &Document{
				ID:      8130,
				Created: time.Date(2026, time.July, 23, 0, 0, 0, 0, plus2),
			},
		},
		{
			name: "invalid created",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/documents/8129/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"id": 8129,
						"created": "not a timestamp"
					}`))
			},
			id:      8129,
			wantErr: cmpopts.AnyError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport:      transport,
				ServerLocation: tc.loc,
			})

			got, _, err := c.GetDocument(t.Context(), tc.id)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetDocument() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetDocument() result diff (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(tc.want.Created.Location(), got.Created.Location(), testutil.EquateTimeLocation()); diff != "" {
					t.Errorf("GetDocument() Created location diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestListDocumentsDateOnly(t *testing.T) {
	plus2 := time.FixedZone("UTC+2", 2*60*60)

	transport := newMockTransport(t)
	transport.RegisterResponder(http.MethodGet, "/api/documents/",
		httpmock.NewStringResponder(http.StatusOK, `{
			"count": 2,
			"next": null,
			"previous": null,
			"results": [
				{ "id": 1, "created": "2026-07-23" },
				{ "id": 2, "created": "2023-06-30T22:00:00Z" }
			]
		}`))

	c := New(Options{
		transport:      transport,
		ServerLocation: plus2,
	})

	got, _, err := c.ListDocuments(t.Context(), ListDocumentsOptions{})
	if err != nil {
		t.Fatalf("ListDocuments() failed: %v", err)
	}

	want := []Document{
		{ID: 1, Created: time.Date(2026, time.July, 23, 0, 0, 0, 0, plus2)},
		{ID: 2, Created: time.Date(2023, time.June, 30, 22, 0, 0, 0, time.UTC)},
	}

	if diff := cmp.Diff(want, got, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("ListDocuments() diff (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(plus2, got[0].Created.Location(), testutil.EquateTimeLocation()); diff != "" {
		t.Errorf("ListDocuments() Created location diff (-want +got):\n%s", diff)
	}
}

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

			got, _, err := c.GetDocumentMetadata(t.Context(), tc.id)

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
				StoragePath:         Int64(500),
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

			got, _, err := c.UploadDocument(t.Context(), tc.r, tc.opts)

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
