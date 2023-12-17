package client

import (
	"context"
)

func (c *Client) storagePathCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/storage_paths/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(StoragePath).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListStoragePathsOptions).Page = page
		},
	}
}

type ListStoragePathsOptions struct {
	ListOptions

	Ordering OrderingSpec   `url:"ordering"`
	Owner    IntFilterSpec  `url:"owner"`
	Name     CharFilterSpec `url:"name"`
	Path     CharFilterSpec `url:"path"`
}

func (c *Client) ListStoragePaths(ctx context.Context, opts ListStoragePathsOptions) ([]StoragePath, *Response, error) {
	return crudList[StoragePath](ctx, c.storagePathCrudOpts(), opts)
}

// ListAllStoragePaths iterates over all storage paths matching the filters
// specified in opts, invoking handler for each.
func (c *Client) ListAllStoragePaths(ctx context.Context, opts ListStoragePathsOptions, handler func(context.Context, StoragePath) error) error {
	return crudListAll[StoragePath](ctx, c.storagePathCrudOpts(), opts, handler)
}

func (c *Client) GetStoragePath(ctx context.Context, id int64) (*StoragePath, *Response, error) {
	return crudGet[StoragePath](ctx, c.storagePathCrudOpts(), id)
}

func (c *Client) CreateStoragePath(ctx context.Context, data *StoragePathFields) (*StoragePath, *Response, error) {
	return crudCreate[StoragePath](ctx, c.storagePathCrudOpts(), data)
}

func (c *Client) UpdateStoragePath(ctx context.Context, id int64, data *StoragePath) (*StoragePath, *Response, error) {
	return crudUpdate[StoragePath](ctx, c.storagePathCrudOpts(), id, data)
}

func (c *Client) PatchStoragePath(ctx context.Context, id int64, data *StoragePathFields) (*StoragePath, *Response, error) {
	return crudPatch[StoragePath](ctx, c.storagePathCrudOpts(), id, data)
}

func (c *Client) DeleteStoragePath(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[StoragePath](ctx, c.storagePathCrudOpts(), id)
}
