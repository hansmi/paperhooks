package client

import (
	"context"
	"errors"
)

func (c *Client) userCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/users/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(User).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListUsersOptions).Page = page
		},
	}
}

type ListUsersOptions struct {
	ListOptions

	Ordering OrderingSpec   `url:"ordering"`
	Username CharFilterSpec `url:"username"`
}

func (c *Client) ListUsers(ctx context.Context, opts ListUsersOptions) ([]User, *Response, error) {
	return crudList[User](ctx, c.userCrudOpts(), opts)
}

// ListAllUsers iterates over all users matching the filters specified in opts,
// invoking handler for each.
func (c *Client) ListAllUsers(ctx context.Context, opts ListUsersOptions, handler func(context.Context, User) error) error {
	return crudListAll[User](ctx, c.userCrudOpts(), opts, handler)
}

func (c *Client) GetUser(ctx context.Context, id int64) (*User, *Response, error) {
	return crudGet[User](ctx, c.userCrudOpts(), id)
}

// GetCurrentUser looks up the authenticated user.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, *Response, error) {
	type uiSettings struct {
		User struct {
			ID *int64 `json:"id"`
		} `json:"user"`
	}

	req := c.newRequest(ctx).SetResult(uiSettings{})

	resp, err := req.Get("api/ui_settings/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	userID := resp.Result().(*uiSettings).User.ID

	if userID == nil {
		return nil, wrapResponse(resp), errors.New("missing user ID in response")
	}

	return c.GetUser(ctx, *userID)
}
