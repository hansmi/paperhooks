package client

import (
	"context"
)

func (c *Client) groupCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/groups/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(Group).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListGroupsOptions).Page = page
		},
	}
}

type ListGroupsOptions struct {
	ListOptions

	Ordering OrderingSpec   `url:"ordering"`
	Name     CharFilterSpec `url:"name"`
}

func (c *Client) ListGroups(ctx context.Context, opts ListGroupsOptions) ([]Group, *Response, error) {
	return crudList[Group](ctx, c.groupCrudOpts(), opts)
}

// ListAllGroups iterates over all groups matching the filters specified in opts,
// invoking handler for each.
func (c *Client) ListAllGroups(ctx context.Context, opts ListGroupsOptions, handler func(context.Context, Group) error) error {
	return crudListAll[Group](ctx, c.groupCrudOpts(), opts, handler)
}

func (c *Client) GetGroup(ctx context.Context, id int64) (*Group, *Response, error) {
	return crudGet[Group](ctx, c.groupCrudOpts(), id)
}
