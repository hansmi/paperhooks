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

func TestGetStatus(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*testing.T, *httpmock.MockTransport)
		wantErr error
		want    *SystemStatus
	}{
		{
			name: "empty",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/status/",
					httpmock.NewStringResponder(http.StatusOK, `{}`))
			},
			want: &SystemStatus{},
		},
		{
			name: "bad JSON",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/status/",
					httpmock.NewStringResponder(http.StatusOK, `{`))
			},
			wantErr: cmpopts.AnyError,
		},
		{
			name: "status",
			setup: func(t *testing.T, transport *httpmock.MockTransport) {
				transport.RegisterResponder(http.MethodGet, "/api/status/",
					httpmock.NewStringResponder(http.StatusOK, `{
						"pngx_version": "2.14.7",
						"server_os": "Linux-6.8.12-8-pve-x86_64-with-glibc2.36",
						"install_type": "bare-metal",
						"storage": {
							"total": 21474836480,
							"available": 13406437376
						},
						"database": {
							"type": "postgresql",
							"url": "paperlessdb",
							"status": "OK",
							"error": null,
							"migration_status": {
								"latest_migration": "mfa.0003_authenticator_type_uniq",
								"unapplied_migrations": []
							}
						},
						"tasks": {
							"redis_url": "redis://localhost:6379",
							"redis_status": "OK",
							"redis_error": null,
							"celery_status": "OK",
							"index_status": "OK",
							"index_last_modified": "2025-02-21T00:01:54.773392Z",
							"index_error": null,
							"classifier_status": "OK",
							"classifier_last_trained": "2025-02-21T20:05:01.589548Z",
							"classifier_error": null
						}
						}`))
			},
			want: &SystemStatus{
				PNGXVersion: "2.14.7",
				ServerOS:    "Linux-6.8.12-8-pve-x86_64-with-glibc2.36",
				InstallType: "bare-metal",
				Storage: SystemStatusStorage{
					Total:     21474836480,
					Available: 13406437376,
				},
				Database: SystemStatusDatabase{
					Type:   "postgresql",
					URL:    "paperlessdb",
					Status: "OK",
					Error:  "",
					MigrationStatus: SystemStatusDatabaseMigration{
						LatestMigration:     "mfa.0003_authenticator_type_uniq",
						UnappliedMigrations: []string{},
					},
				},
				Tasks: SystemStatusTasks{
					RedisURL:              "redis://localhost:6379",
					RedisStatus:           "OK",
					CeleryStatus:          "OK",
					IndexStatus:           "OK",
					IndexLastModified:     time.Date(2025, time.February, 21, 0, 1, 54, 773392000, time.UTC),
					ClassifierStatus:      "OK",
					ClassifierLastTrained: time.Date(2025, time.February, 21, 20, 5, 1, 589548000, time.UTC),
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

			got, _, err := c.GetStatus(context.Background())

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetStatus() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("GetStatus() diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
