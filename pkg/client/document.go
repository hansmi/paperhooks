package client

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/go-querystring/query"
)

type Document struct {
	// ID of the document. Read-only.
	ID int64 `json:"id"`

	// Title of the document.
	Title string `json:"title"`

	// Plain-text content of the document.
	Content string `json:"content"`

	// List of tag IDs assigned to this document, or empty list.
	Tags []int64 `json:"tags"`

	// Document type of this document, or nil.
	DocumentType *int64 `json:"document_type"`

	// Correspondent of this document or nil.
	Correspondent *int64 `json:"correspondent"`

	// Storage path of this document or nil.
	StoragePath *int64 `json:"storage_path"`

	// The date time at which this document was created.
	Created time.Time `json:"created"`

	// The date at which this document was last edited in paperless. Read-only.
	Modified time.Time `json:"modified"`

	// The date at which this document was added to paperless. Read-only.
	Added time.Time `json:"added"`

	// The identifier of this document in a physical document archive.
	ArchiveSerialNumber *int64 `json:"archive_serial_number"`

	// Verbose filename of the original document. Read-only.
	OriginalFileName string `json:"original_file_name"`

	// Verbose filename of the archived document. Read-only. Nil if no archived document is available.
	ArchivedFileName *string `json:"archived_file_name"`
}

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
	}
}

type ListDocumentsOptions struct {
	ListOptions

	Ordering            OrderingSpec         `url:"ordering"`
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

func (c *Client) ListDocuments(ctx context.Context, opts *ListDocumentsOptions) ([]Document, *Response, error) {
	return crudList[Document](ctx, c.documentCrudOpts(), opts)
}

func (c *Client) GetDocument(ctx context.Context, id int64) (*Document, *Response, error) {
	return crudGet[Document](ctx, c.documentCrudOpts(), id)
}

func (c *Client) UpdateDocument(ctx context.Context, id int64, data *Document) (*Document, *Response, error) {
	return crudUpdate[Document](ctx, c.documentCrudOpts(), id, data)
}

func (c *Client) DeleteDocument(ctx context.Context, id int64) (*Response, error) {
	return crudDelete[Document](ctx, c.documentCrudOpts(), id)
}

func (c *Client) GetDocumentMetadata(ctx context.Context, id int64) (*DocumentMetadata, error) {
	resp, err := c.newRequest(ctx).
		SetResult(DocumentMetadata{}).
		Get(fmt.Sprintf("api/documents/%d/metadata/", id))

	if err := convertError(err, resp); err != nil {
		return nil, err
	}

	return resp.Result().(*DocumentMetadata), nil
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
func (c *Client) UploadDocument(ctx context.Context, r io.Reader, opts *DocumentUploadOptions) (*DocumentUpload, *Response, error) {
	req := c.newRequest(ctx).
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

	result := &DocumentUpload{
		TaskID: string(resp.Body()),
	}

	return result, wrapResponse(resp), nil
}
