package client

import "context"

type DocumentFileType struct {
	MimeType      string `json:"mime_type"`
	MimeTypeCount int64  `json:"mime_type_count"`
}

func (c *Client) GetStatistics(ctx context.Context) (*Statistics, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult(&Statistics{}).
		Get("api/statistics/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*Statistics), wrapResponse(resp), nil
}
