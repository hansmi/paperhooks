package client

import (
	"context"
	"time"
)

type SystemStatus struct {
	PNGXVersion string               `json:"pngx_version"`
	ServerOS    string               `json:"server_os"`
	InstallType string               `json:"install_type"`
	Storage     SystemStatusStorage  `json:"storage"`
	Database    SystemStatusDatabase `json:"database"`
	Tasks       SystemStatusTasks    `json:"tasks"`
}

type SystemStatusStorage struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
}

type SystemStatusDatabase struct {
	Type            string                        `json:"type"`
	URL             string                        `json:"url"`
	Status          string                        `json:"status"`
	Error           string                        `json:"error"`
	MigrationStatus SystemStatusDatabaseMigration `json:"migration_status"`
}

type SystemStatusDatabaseMigration struct {
	LatestMigration     string   `json:"latest_migration"`
	UnappliedMigrations []string `json:"unapplied_migrations"`
}

type SystemStatusTasks struct {
	RedisURL              string    `json:"redis_url"`
	RedisStatus           string    `json:"redis_status"`
	RedisError            string    `json:"redis_error"`
	CeleryStatus          string    `json:"celery_status"`
	CeleryURL             string    `json:"celery_url"`
	CeleryError           string    `json:"celery_error"`
	IndexStatus           string    `json:"index_status"`
	IndexLastModified     time.Time `json:"index_last_modified"`
	IndexError            string    `json:"index_error"`
	ClassifierStatus      string    `json:"classifier_status"`
	ClassifierLastTrained time.Time `json:"classifier_last_trained"`
	ClassifierError       string    `json:"classifier_error"`
	SanityCheckStatus     string    `json:"sanity_check_status"`
	SanityCheckLastRun    time.Time `json:"sanity_check_last_run"`
	SanityCheckError      string    `json:"sanity_check_error"`
}

func (c *Client) GetStatus(ctx context.Context) (*SystemStatus, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult(&SystemStatus{}).
		Get("api/status/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*SystemStatus), wrapResponse(resp), nil
}
