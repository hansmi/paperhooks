package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hansmi/paperhooks/pkg/client"
)

type destructiveTests struct {
	logger *log.Logger
	client *client.Client
	mark   string
}

func (t *destructiveTests) tags(ctx context.Context) error {
	name := fmt.Sprintf("%s test tag", t.mark)

	t.logger.Printf("Create tag %q", name)

	tag, _, err := t.client.CreateTag(ctx, client.NewTagFields().
		Name(name).
		MatchingAlgorithm(client.MatchAny))
	if err != nil {
		return fmt.Errorf("creating tag %s failed: %w", name, err)
	}

	if !(tag.Name == name && tag.MatchingAlgorithm == client.MatchAny) {
		return fmt.Errorf("tag settings differ from configuration: %#v", *tag)
	}

	t.logger.Printf("Update tag %q without making changes: %#v", name, *tag)

	tag, _, err = t.client.UpdateTag(ctx, tag.ID, tag)
	if err != nil {
		return fmt.Errorf("updating tag %s without changes failed: %w", name, err)
	}

	tag.Name = name + " modified"
	tag.Color = client.Color{R: 0xFF}
	tag.IsInboxTag = true
	tag.IsInsensitive = true
	tag.MatchingAlgorithm = client.MatchFuzzy
	tag.Match = name

	t.logger.Printf("Update tag %q with changes: %#v", name, *tag)

	tag, _, err = t.client.UpdateTag(ctx, tag.ID, tag)
	if err != nil {
		return fmt.Errorf("updating tag %s with changes failed: %w", name, err)
	}

	name += " modified"

	t.logger.Printf("List tags with name %q", name)

	if matches, _, err := t.client.ListTags(ctx, client.ListTagsOptions{
		Ordering: client.OrderingSpec{
			Field: "name",
		},
		Name: client.CharFilterSpec{
			EqualsIgnoringCase: client.String(name),
		},
	}); err != nil {
		return fmt.Errorf("listing tags failed: %v", err)
	} else if len(matches) != 1 {
		return fmt.Errorf("listing tags did not return exactly one match for %q: %+v", name, matches)
	}

	t.logger.Printf("Delete tag %q", name)

	if _, err := t.client.DeleteTag(ctx, tag.ID); err != nil {
		return fmt.Errorf("deleting tag %s failed: %w", name, err)
	}

	_, _, err = t.client.GetTag(ctx, tag.ID)
	if detail, ok := err.(*client.RequestError); !(ok && detail.StatusCode == http.StatusNotFound) {
		return fmt.Errorf("getting tag %s did not return HTTP 404: %w", name, err)
	}

	return nil
}

func (t *destructiveTests) uploadDocument(ctx context.Context) error {
	imgBytes, err := makeRandomImage(100, 100)
	if err != nil {
		return err
	}

	t.logger.Printf("Upload document with a generated image")

	result, _, err := t.client.UploadDocument(ctx,
		bytes.NewReader(imgBytes),
		client.DocumentUploadOptions{
			Filename: t.mark + ".png",
			Title:    t.mark,
		})
	if err != nil {
		return fmt.Errorf("uploading document failed: %w", err)
	}

	t.logger.Printf("Task ID from document upload: %s", result.TaskID)

	if task, err := t.client.WaitForTask(ctx, result.TaskID, client.WaitForTaskOptions{}); err != nil {
		return fmt.Errorf("document upload: %w", err)
	} else {
		t.logger.Printf("Task finished: %+v", task)
	}

	return nil
}
