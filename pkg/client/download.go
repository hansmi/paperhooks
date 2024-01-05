package client

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"go.uber.org/multierr"
)

type DownloadResult struct {
	// MIME content type (e.g. "application/pdf").
	ContentType string

	// Parameters for body content type (e.g. "charset").
	ContentTypeParams map[string]string

	// The preferred filename as reported by the server (if any).
	Filename string

	// Length of the downloaded body in bytes.
	Length int64
}

func (c *Client) download(ctx context.Context, w io.Writer, url string, expectDisposition bool) (_ *DownloadResult, _ *Response, err error) {
	req := c.newRequest(ctx).
		SetDoNotParseResponse(true)

	resp, err := req.Get(url)

	if !(resp == nil || resp.RawBody() == nil) {
		defer multierr.AppendFunc(&err, resp.RawBody().Close)
	}

	if err := convertError(err, resp); err != nil {
		return nil, wrapResponse(resp), err
	}

	result := &DownloadResult{}

	if result.ContentType, result.ContentTypeParams, err = mime.ParseMediaType(resp.Header().Get("Content-Type")); err != nil {
		return nil, wrapResponse(resp), fmt.Errorf("invalid content-type header: %w", err)
	}

	if contentDisposition := resp.Header().Get("Content-Disposition"); contentDisposition == "" {
		if expectDisposition {
			c.logger.Warnf("%s: missing Content-Disposition header", req.URL)
		}
	} else if _, params, err := mime.ParseMediaType(contentDisposition); err != nil {
		c.logger.Warnf("%s: parsing content-disposition header failed: %w", req.URL, err)
	} else if filename, ok := params["filename"]; ok && filename != "" {
		result.Filename = filepath.Base(filepath.Clean(params["filename"]))
	}

	result.Length, err = io.Copy(w, resp.RawBody())
	if err != nil {
		return nil, wrapResponse(resp), err
	}

	return result, wrapResponse(resp), nil
}

// DownloadDocumentOriginal retrieves the document in the format originally
// consumed by Paperless. The file format can be determined using
// [DownloadResult.ContentType].
//
// The content of the document is written to the given writer. To verify that
// the document is complete (the HTTP request may have been terminated early)
// the size and/or checksum can be verified with [GetDocumentMetadata].
func (c *Client) DownloadDocumentOriginal(ctx context.Context, w io.Writer, id int64) (*DownloadResult, *Response, error) {
	return c.download(ctx, w, fmt.Sprintf("api/documents/%d/download/?original=true", id), true)
}

// DownloadDocumentArchived retrieves an archived PDF/A file generated from the
// originally consumed file. The archived version may not be available and the
// API may return the original. [DownloadDocumentOriginal] for additional
// details.
func (c *Client) DownloadDocumentArchived(ctx context.Context, w io.Writer, id int64) (*DownloadResult, *Response, error) {
	return c.download(ctx, w, fmt.Sprintf("api/documents/%d/download/", id), true)
}

// DownloadDocumentThumbnail retrieves a preview image of the document. See
// [DownloadDocumentOriginal] for additional details.
func (c *Client) DownloadDocumentThumbnail(ctx context.Context, w io.Writer, id int64) (*DownloadResult, *Response, error) {
	return c.download(ctx, w, fmt.Sprintf("api/documents/%d/thumb/", id), false)
}
