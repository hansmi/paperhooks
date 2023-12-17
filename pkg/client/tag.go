package client

import (
	"context"
	"encoding/json"
)

type Tag struct {
	ID                int64             `json:"id"`
	Slug              string            `json:"slug"`
	Name              string            `json:"name"`
	Color             Color             `json:"color"`
	TextColor         Color             `json:"text_color"`
	Match             string            `json:"match"`
	MatchingAlgorithm MatchingAlgorithm `json:"matching_algorithm"`
	IsInsensitive     bool              `json:"is_insensitive"`
	IsInboxTag        bool              `json:"is_inbox_tag"`
	DocumentCount     int64             `json:"document_count"`
}

type TagFields struct {
	objectFields
}

var _ json.Marshaler = (*TagFields)(nil)

func NewTagFields() *TagFields {
	return &TagFields{
		objectFields: objectFields{},
	}
}

func (f *TagFields) Name(name string) *TagFields {
	f.set("name", name)
	return f
}

func (f *TagFields) Color(c Color) *TagFields {
	f.set("color", c)
	return f
}

func (f *TagFields) TextColor(c Color) *TagFields {
	f.set("text_color", c)
	return f
}

func (c *Client) tagCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/tags/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(Tag).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListTagsOptions).Page = page
		},
	}
}

type ListTagsOptions struct {
	ListOptions

	Ordering OrderingSpec   `url:"ordering"`
	Owner    IntFilterSpec  `url:"owner"`
	Name     CharFilterSpec `url:"name"`
}

func (c *Client) ListTags(ctx context.Context, opts ListTagsOptions) ([]Tag, *Response, error) {
	return crudList[Tag](ctx, c.tagCrudOpts(), opts)
}

// ListAllTags iterates over all tags matching the filters specified in opts,
// invoking handler for each.
func (c *Client) ListAllTags(ctx context.Context, opts ListTagsOptions, handler func(context.Context, Tag) error) error {
	return crudListAll[Tag](ctx, c.tagCrudOpts(), opts, handler)
}

func (c *Client) GetTag(ctx context.Context, id int64) (*Tag, *Response, error) {
	return crudGet[Tag](ctx, c.tagCrudOpts(), id)
}

func (c *Client) CreateTag(ctx context.Context, data *TagFields) (*Tag, *Response, error) {
	return crudCreate[Tag](ctx, c.tagCrudOpts(), data)
}

func (c *Client) UpdateTag(ctx context.Context, id int64, data *Tag) (*Tag, *Response, error) {
	return crudUpdate[Tag](ctx, c.tagCrudOpts(), id, data)
}

func (c *Client) PatchTag(ctx context.Context, id int64, data *TagFields) (*Tag, *Response, error) {
	return crudPatch[Tag](ctx, c.tagCrudOpts(), id, data)
}

func (c *Client) DeleteTag(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[Tag](ctx, c.tagCrudOpts(), id)
}
