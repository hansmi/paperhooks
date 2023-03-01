package client

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type CommentUser struct {
	ID   int64  `json:"id"`
	Name string `json:"username"`
}

type Comment struct {
	ID      int64       `json:"id"`
	Text    string      `json:"comment"`
	Created time.Time   `json:"created"`
	User    CommentUser `json:"user"`
}

func (c *Client) commentsURL(documentID int64) string {
	return fmt.Sprintf("api/documents/%d/comments/", documentID)
}

func (c *Client) ListComments(ctx context.Context, documentID int64) ([]Comment, *Response, error) {
	req := c.newRequest(ctx).SetResult([]Comment(nil))

	resp, err := req.Get(c.commentsURL(documentID))

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return *resp.Result().(*[]Comment), wrapResponse(resp), nil
}

func (c *Client) CreateComment(ctx context.Context, documentID int64, data *Comment) (*Response, error) {
	req := c.newRequest(ctx).SetFormData(map[string]string{
		"comment": data.Text,
	})

	resp, err := req.Post(c.commentsURL(documentID))

	if err := convertError(err, resp); err != nil {
		return wrapResponse(resp), err
	}

	return wrapResponse(resp), nil
}

func (c *Client) DeleteComment(ctx context.Context, documentID, id int64) (*Response, error) {
	req := c.newRequest(ctx).SetQueryParam("id", strconv.FormatInt(id, 10))

	resp, err := req.Delete(c.commentsURL(documentID))

	if err := convertError(err, resp); err != nil {
		return wrapResponse(resp), err
	}

	return wrapResponse(resp), nil
}
