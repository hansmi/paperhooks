package client

import (
	"context"
)

type DocumentType struct {
	ID                int64             `json:"id"`
	Slug              string            `json:"slug"`
	Name              string            `json:"name"`
	Match             string            `json:"match"`
	MatchingAlgorithm MatchingAlgorithm `json:"matching_algorithm"`
	IsInsensitive     bool              `json:"is_insensitive"`
	DocumentCount     int64             `json:"document_count"`
}

func (c *Client) documentTypeCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/document_types/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(DocumentType).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListDocumentTypesOptions).Page = page
		},
	}
}

type ListDocumentTypesOptions struct {
	ListOptions

	Ordering OrderingSpec   `url:"ordering"`
	Name     CharFilterSpec `url:"name"`
}

func (c *Client) ListDocumentTypes(ctx context.Context, opts ListDocumentTypesOptions) ([]DocumentType, *Response, error) {
	return crudList[DocumentType](ctx, c.documentTypeCrudOpts(), opts)
}

// ListAllDocumentTypes iterates over all document types matching the filters
// specified in opts, invoking handler for each.
func (c *Client) ListAllDocumentTypes(ctx context.Context, opts ListDocumentTypesOptions, handler func(context.Context, DocumentType) error) error {
	return crudListAll[DocumentType](ctx, c.documentTypeCrudOpts(), opts, handler)
}

func (c *Client) GetDocumentType(ctx context.Context, id int64) (*DocumentType, *Response, error) {
	return crudGet[DocumentType](ctx, c.documentTypeCrudOpts(), id)
}

func (c *Client) CreateDocumentType(ctx context.Context, data *DocumentType) (*DocumentType, *Response, error) {
	return crudCreate[DocumentType](ctx, c.documentTypeCrudOpts(), data)
}

func (c *Client) UpdateDocumentType(ctx context.Context, id int64, data *DocumentType) (*DocumentType, *Response, error) {
	return crudUpdate[DocumentType](ctx, c.documentTypeCrudOpts(), id, data)
}

func (c *Client) DeleteDocumentType(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[DocumentType](ctx, c.documentTypeCrudOpts(), id)
}
