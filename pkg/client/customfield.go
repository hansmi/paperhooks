package client

import (
	"context"
)

type CustomField struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	DataType string `json:"data_type"`
}

func (c *Client) customFieldCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/custom_fields/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(CustomField).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListCustomFieldsOptions).Page = page
		},
	}
}

type ListCustomFieldsOptions struct {
	ListOptions
}

func (c *Client) ListCustomFields(ctx context.Context, opts ListCustomFieldsOptions) ([]CustomField, *Response, error) {
	return crudList[CustomField](ctx, c.customFieldCrudOpts(), opts)
}

// ListAllCustomFields iterates over all custom fields matching the filters
// specified in opts, invoking handler for each.
func (c *Client) ListAllCustomFields(ctx context.Context, opts ListCustomFieldsOptions, handler func(context.Context, CustomField) error) error {
	return crudListAll[CustomField](ctx, c.customFieldCrudOpts(), opts, handler)
}

func (c *Client) GetCustomField(ctx context.Context, id int64) (*CustomField, *Response, error) {
	return crudGet[CustomField](ctx, c.customFieldCrudOpts(), id)
}

func (c *Client) CreateCustomField(ctx context.Context, data *CustomField) (*CustomField, *Response, error) {
	return crudCreate[CustomField](ctx, c.customFieldCrudOpts(), data)
}

func (c *Client) UpdateCustomField(ctx context.Context, id int64, data *CustomField) (*CustomField, *Response, error) {
	return crudUpdate[CustomField](ctx, c.customFieldCrudOpts(), id, data)
}

func (c *Client) DeleteCustomField(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[CustomField](ctx, c.customFieldCrudOpts(), id)
}

type CustomFieldInstance struct {
	Field int64 `json:"field"`
	Value any   `json:"value"`
}
