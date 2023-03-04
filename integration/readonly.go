package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hansmi/paperhooks/pkg/client"
)

type readOnlyTests struct {
	logger *log.Logger
	client *client.Client
}

func (t *readOnlyTests) tags(ctx context.Context) error {
	opt := &client.ListTagsOptions{}
	all := []client.Tag{}

	for {
		tags, resp, err := t.client.ListTags(ctx, opt)
		if err != nil {
			return fmt.Errorf("listing tags failed: %w", err)
		}

		all = append(all, tags...)

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}

	t.logger.Printf("Received %d tags. Fetching all of them.", len(all))

	for _, i := range all {
		if _, _, err := t.client.GetTag(ctx, i.ID); err != nil {
			return fmt.Errorf("getting tag %d failed: %w", i.ID, err)
		}
	}

	return nil
}

func (t *readOnlyTests) correspondents(ctx context.Context) error {
	opt := &client.ListCorrespondentsOptions{}
	all := []client.Correspondent{}

	for {
		correspondents, resp, err := t.client.ListCorrespondents(ctx, opt)
		if err != nil {
			return fmt.Errorf("listing correspondents failed: %w", err)
		}

		all = append(all, correspondents...)

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}

	t.logger.Printf("Received %d correspondents.", len(all))

	for _, i := range all {
		if _, _, err := t.client.GetCorrespondent(ctx, i.ID); err != nil {
			return fmt.Errorf("getting correspondent %d failed: %w", i.ID, err)
		}
	}

	return nil
}

func (t *readOnlyTests) documentTypes(ctx context.Context) error {
	opt := &client.ListDocumentTypesOptions{}
	all := []client.DocumentType{}

	for {
		documentTypes, resp, err := t.client.ListDocumentTypes(ctx, opt)
		if err != nil {
			return fmt.Errorf("listing document types failed: %w", err)
		}

		all = append(all, documentTypes...)

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}

	t.logger.Printf("Received %d document types.", len(all))

	for _, i := range all {
		if _, _, err := t.client.GetDocumentType(ctx, i.ID); err != nil {
			return fmt.Errorf("getting document type %d failed: %w", i.ID, err)
		}
	}

	return nil
}

func (t *readOnlyTests) storagePaths(ctx context.Context) error {
	opt := &client.ListStoragePathsOptions{}
	all := []client.StoragePath{}

	for {
		storagePaths, resp, err := t.client.ListStoragePaths(ctx, opt)
		if err != nil {
			return fmt.Errorf("listing storage paths failed: %w", err)
		}

		all = append(all, storagePaths...)

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}

	t.logger.Printf("Received %d storage paths.", len(all))

	for _, i := range all {
		if _, _, err := t.client.GetStoragePath(ctx, i.ID); err != nil {
			return fmt.Errorf("getting storage path %d failed: %w", i.ID, err)
		}
	}

	return nil
}

func (t *readOnlyTests) documents(ctx context.Context) error {
	const examineCount = 10

	opt := &client.ListDocumentsOptions{}
	all := []client.Document{}

	for {
		documents, resp, err := t.client.ListDocuments(ctx, opt)
		if err != nil {
			return fmt.Errorf("listing documents failed: %w", err)
		}

		all = append(all, documents...)

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}

	t.logger.Printf("Received %d documents.", len(all))

	for idx, i := range all {
		if _, _, err := t.client.GetDocument(ctx, i.ID); err != nil {
			return fmt.Errorf("getting document %d failed: %w", i.ID, err)
		}

		if md, err := t.client.GetDocumentMetadata(ctx, i.ID); err != nil {
			return fmt.Errorf("getting document %d metadata failed: %w", i.ID, err)
		} else {
			t.logger.Printf("Document %d metadata: %+v", i.ID, *md)
		}

		for _, x := range []struct {
			name string
			fn   func(context.Context, io.Writer, int64) (*client.DownloadResult, *client.Response, error)
		}{
			{
				name: "original",
				fn:   t.client.DownloadDocumentOriginal,
			},
			{
				name: "archived",
				fn:   t.client.DownloadDocumentArchived,
			},
			{
				name: "thumbnail",
				fn:   t.client.DownloadDocumentThumbnail,
			},
		} {
			t.logger.Printf("Download %s version of document %d.", x.name, i.ID)

			if r, _, err := x.fn(ctx, io.Discard, i.ID); err != nil {
				return fmt.Errorf("download of document %d failed: %w", i.ID, err)
			} else {
				t.logger.Printf("Received %d bytes for filename %q with content type %q.",
					r.Length, r.Filename, r.ContentType)
			}
		}

		if idx >= examineCount {
			break
		}
	}

	return nil
}

func (t *readOnlyTests) tasks(ctx context.Context) error {
	tasks, _, err := t.client.ListTasks(ctx)
	if err != nil {
		return fmt.Errorf("listing tasks failed: %w", err)
	}

	t.logger.Printf("Received %d tasks.", len(tasks))

	return nil
}

func (t *readOnlyTests) logs(ctx context.Context) error {
	logs, _, err := t.client.ListLogs(ctx)
	if err != nil {
		return fmt.Errorf("listing logs failed: %w", err)
	}

	t.logger.Printf("Received log names: %v", logs)

	for _, name := range logs {
		entries, _, err := t.client.GetLog(ctx, name)
		if err != nil {
			var re *client.RequestError

			if !(errors.As(err, &re) && re.StatusCode == http.StatusNotFound) {
				return fmt.Errorf("fetching entries for log %q failed: %w", name, err)
			}

			entries = nil
		}

		t.logger.Printf("Received %d entries for log %q.", len(entries), name)
	}

	return nil
}
