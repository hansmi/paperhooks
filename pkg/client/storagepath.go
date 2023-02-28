package client

import (
	"context"
)

type StoragePath struct {
	ID                int64             `json:"id"`
	Slug              string            `json:"slug"`
	Name              string            `json:"name"`
	Match             string            `json:"match"`
	MatchingAlgorithm MatchingAlgorithm `json:"matching_algorithm"`
	IsInsensitive     bool              `json:"is_insensitive"`
	DocumentCount     int64             `json:"document_count"`
}

func (c *Client) storagePathCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/storage_paths/",
		newRequest: c.newRequest,
	}
}

type ListStoragePathsOptions struct {
	ListOptions

	Ordering OrderingSpec   `url:"ordering"`
	Name     CharFilterSpec `url:"name"`
	Path     CharFilterSpec `url:"path"`
}

func (c *Client) ListStoragePaths(ctx context.Context, opts *ListStoragePathsOptions) ([]StoragePath, *Response, error) {
	return crudList[StoragePath](ctx, c.storagePathCrudOpts(), opts)
}

func (c *Client) GetStoragePath(ctx context.Context, id int64) (*StoragePath, *Response, error) {
	return crudGet[StoragePath](ctx, c.storagePathCrudOpts(), id)
}

func (c *Client) CreateStoragePath(ctx context.Context, data *StoragePath) (*StoragePath, *Response, error) {
	return crudCreate[StoragePath](ctx, c.storagePathCrudOpts(), data)
}

func (c *Client) UpdateStoragePath(ctx context.Context, id int64, data *StoragePath) (*StoragePath, *Response, error) {
	return crudUpdate[StoragePath](ctx, c.storagePathCrudOpts(), id, data)
}

func (c *Client) DeleteStoragePath(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[StoragePath](ctx, c.storagePathCrudOpts(), id)
}
