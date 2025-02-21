package client

import (
	"context"
)

type Status struct {
	PNGXVersion string         `json:"pngx_version"`
	ServerOS    string         `json:"server_os"`
	InstallType string         `json:"install_type"`
	Storage     StorageStatus  `json:"storage"`
	Database    DatabaseStatus `json:"database"`
	Tasks       TasksStatus    `json:"tasks"`
}

type StorageStatus struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
}

type DatabaseStatus struct {
	Type            string                  `json:"type"`
	URL             string                  `json:"url"`
	Status          string                  `json:"status"`
	Error           string                  `json:"error"`
	MigrationStatus DatabaseMigrationStatus `json:"migration_status"`
}

type DatabaseMigrationStatus struct {
	LatestMigration     string   `json:"latest_migration"`
	UnappliedMigrations []string `json:"unapplied_migrations"`
}

type TasksStatus struct {
	RedisURL              string `json:"redis_url"`
	RedisStatus           string `json:"redis_status"`
	RedisError            string `json:"redis_error"`
	CeleryStatus          string `json:"celery_status"`
	IndexStatus           string `json:"index_status"`
	IndexLastModified     string `json:"index_last_modified"`
	IndexError            string `json:"index_error"`
	ClassifierStatus      string `json:"classifier_status"`
	ClassifierLastTrained string `json:"classifier_last_trained"`
	ClassifierError       string `json:"classifier_error"`
}

func (c *Client) GetStatus(ctx context.Context) (*Status, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult(&Status{}).
		Get("api/status/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*Status), wrapResponse(resp), nil
}
