package client

import "context"

// Ping tests whether the API is available.
func (c *Client) Ping(ctx context.Context) error {
	resp, err := c.newRequest(ctx).
		Get("api/")

	if err := convertError(err, resp); err != nil {
		return err
	}

	return nil
}
