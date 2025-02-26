package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
)

func TestGetStatistics(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		wantErr error
		want    *Statistics
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/statistics/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			want: &Statistics{},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/statistics/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "remote_version",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/statistics/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"documents_total": 1447,
						"documents_inbox": 273,
						"inbox_tag": 1,
						"inbox_tags": [
							1
						],
						"document_file_type_counts": [
							{
							"mime_type": "application/pdf",
							"mime_type_count": 1397
							},
							{
							"mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
							"mime_type_count": 37
							}
						],
						"character_count": 11048372,
						"tag_count": 55,
						"correspondent_count": 201,
						"document_type_count": 42,
						"storage_path_count": 0,
						"current_asn": 0
					}`))
			},
			want: &Statistics{
				DocumentsTotal: 1447,
				DocumentsInbox: 273,
				InboxTag:       1,
				InboxTags:      []int64{1},
				DocumentFileTypeCounts: []StatisticsDocumentFileType{
					{
						MimeType:      "application/pdf",
						MimeTypeCount: 1397,
					},
					{
						MimeType:      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
						MimeTypeCount: 37,
					},
				},
				CharacterCount:     11048372,
				TagCount:           55,
				CorrespondentCount: 201,
				DocumentTypeCount:  42,
				StoragePathCount:   0,
				CurrentAsn:         0,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport := newMockTransport(t)

			tc.setup(t, transport)

			c := New(Options{
				transport: transport,
			})

			got, _, err := c.GetStatistics(context.Background())

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetStatistics() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetStatistics() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
