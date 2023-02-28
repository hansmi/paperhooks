package client

import "context"

func (c *Client) Ping(ctx context.Context) error {
	resp, err := c.newRequest(ctx).
		Get("api/")

	if err := convertError(err, resp); err != nil {
		return err
	}

	return nil
}
