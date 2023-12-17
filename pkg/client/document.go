package client

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/go-querystring/query"
)

type DocumentVersionMetadata struct {
	Namespace string `json:"namespace"`
	Prefix    string `json:"prefix"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type DocumentMetadata struct {
	OriginalFilename      string                    `json:"original_filename"`
	OriginalMediaFilename string                    `json:"media_filename"`
	OriginalChecksum      string                    `json:"original_checksum"`
	OriginalSize          int64                     `json:"original_size"`
	OriginalMimeType      string                    `json:"original_mime_type"`
	OriginalMetadata      []DocumentVersionMetadata `json:"original_metadata"`

	HasArchiveVersion    bool                      `json:"has_archive_version"`
	ArchiveMediaFilename string                    `json:"archive_media_filename"`
	ArchiveChecksum      string                    `json:"archive_checksum"`
	ArchiveSize          int64                     `json:"archive_size"`
	ArchiveMetadata      []DocumentVersionMetadata `json:"archive_metadata"`

	Language string `json:"lang"`
}

func (c *Client) documentCrudOpts() crudOptions {
	return crudOptions{
		base:       "api/documents/",
		newRequest: c.newRequest,
		getID: func(v any) int64 {
			return v.(Document).ID
		},
		setPage: func(opts any, page *PageToken) {
			opts.(*ListDocumentsOptions).Page = page
		},
	}
}

type ListDocumentsOptions struct {
	ListOptions

	Ordering            OrderingSpec         `url:"ordering"`
	Owner               IntFilterSpec        `url:"owner"`
	Title               CharFilterSpec       `url:"title"`
	Content             CharFilterSpec       `url:"content"`
	ArchiveSerialNumber IntFilterSpec        `url:"archive_serial_number"`
	Created             DateTimeFilterSpec   `url:"created"`
	Added               DateTimeFilterSpec   `url:"added"`
	Modified            DateTimeFilterSpec   `url:"modified"`
	Correspondent       ForeignKeyFilterSpec `url:"correspondent"`
	Tags                ForeignKeyFilterSpec `url:"tags"`
	DocumentType        ForeignKeyFilterSpec `url:"document_type"`
	StoragePath         ForeignKeyFilterSpec `url:"storage_path"`
}

func (c *Client) ListDocuments(ctx context.Context, opts ListDocumentsOptions) ([]Document, *Response, error) {
	return crudList[Document](ctx, c.documentCrudOpts(), opts)
}

// ListAllDocuments iterates over all documents matching the filters specified
// in opts, invoking handler for each.
func (c *Client) ListAllDocuments(ctx context.Context, opts ListDocumentsOptions, handler func(context.Context, Document) error) error {
	return crudListAll[Document](ctx, c.documentCrudOpts(), opts, handler)
}

func (c *Client) GetDocument(ctx context.Context, id int64) (*Document, *Response, error) {
	return crudGet[Document](ctx, c.documentCrudOpts(), id)
}

func (c *Client) UpdateDocument(ctx context.Context, id int64, data *Document) (*Document, *Response, error) {
	return crudUpdate[Document](ctx, c.documentCrudOpts(), id, data)
}

func (c *Client) PatchDocument(ctx context.Context, id int64, data *DocumentFields) (*Document, *Response, error) {
	return crudPatch[Document](ctx, c.documentCrudOpts(), id, data)
}

func (c *Client) DeleteDocument(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[Document](ctx, c.documentCrudOpts(), id)
}

func (c *Client) GetDocumentMetadata(ctx context.Context, id int64) (*DocumentMetadata, *Response, error) {
	resp, err := c.newRequest(ctx).
		SetResult(DocumentMetadata{}).
		Get(fmt.Sprintf("api/documents/%d/metadata/", id))

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return resp.Result().(*DocumentMetadata), wrapResponse(resp), nil
}

type DocumentUploadOptions struct {
	Filename string `url:"-"`

	// Title for the document.
	Title string `url:"title,omitempty"`

	// Datetime at which the document was created.
	Created time.Time `url:"created,omitempty"`

	// ID of a correspondent for the document.
	Correspondent *int64 `url:"correspondent,omitempty"`

	// ID of a document type for the document.
	DocumentType *int64 `url:"document_type,omitempty"`

	// Tag IDs for the document.
	Tags []int64 `url:"tags,omitempty"`

	// Archive serial number to set on the document.
	ArchiveSerialNumber *int64 `url:"archive_serial_number,omitempty"`
}

type DocumentUpload struct {
	TaskID string
}

// Upload a file. Returns immediately and without error if the document
// consumption process was started successfully. No additional status
// information about the consumption process is available immediately. Poll the
// returned task ID to wait for the consumption.
func (c *Client) UploadDocument(ctx context.Context, r io.Reader, opts DocumentUploadOptions) (*DocumentUpload, *Response, error) {
	result := &DocumentUpload{}

	req := c.newRequest(ctx).
		SetResult(&result.TaskID).
		SetFileReader("document", filepath.Base(opts.Filename), r)

	if values, err := query.Values(opts); err != nil {
		return nil, nil, err
	} else {
		req.SetFormDataFromValues(values)
	}

	resp, err := req.Post("api/documents/post_document/")

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	return result, wrapResponse(resp), nil
}
