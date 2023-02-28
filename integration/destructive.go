package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/paperhooks/pkg/client"
)

type destructiveTests struct {
	logger *log.Logger
	client *client.Client
	mark   string
}

func (t *destructiveTests) tags(ctx context.Context) error {
	name := fmt.Sprintf("%s test tag", t.mark)
	compareOpts := []cmp.Option{
		// Fields controlled by server
		cmpopts.IgnoreFields(client.Tag{}, "ID", "Slug", "TextColor"),
	}

	t.logger.Printf("Create tag %q", name)

	tag, _, err := t.client.CreateTag(ctx, &client.Tag{
		Name:              name,
		MatchingAlgorithm: client.MatchAny,
	})
	if err != nil {
		return fmt.Errorf("creating tag %s failed: %w", name, err)
	}

	if diff := cmp.Diff(client.Tag{
		Name:              name,
		MatchingAlgorithm: client.MatchAny,
	}, *tag, compareOpts...); diff != "" {
		return fmt.Errorf("tag diff (-want +got):\n%s", diff)
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

	origname := name
	name += " modified"

	t.logger.Printf("List tags with name %q", name)

	if matches, _, err := t.client.ListTags(ctx, &client.ListTagsOptions{
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

	if diff := cmp.Diff(client.Tag{
		Name:              name,
		Color:             client.Color{R: 0xFF},
		IsInboxTag:        true,
		IsInsensitive:     true,
		MatchingAlgorithm: client.MatchFuzzy,
		Match:             origname,
	}, *tag, compareOpts...); diff != "" {
		return fmt.Errorf("tag diff (-want +got):\n%s", diff)
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
