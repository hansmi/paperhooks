package client

import (
	"context"
)

func (c *Client) GetRemoteVersion(ctx context.Context) (*RemoteVersion, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult(&RemoteVersion{}).
		Get("api/remote_version/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*RemoteVersion), wrapResponse(resp), nil
}
