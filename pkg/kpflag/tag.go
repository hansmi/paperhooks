package kpflag

import "github.com/hansmi/paperhooks/pkg/client"

// RegisterListTagsOptions adds flags for filtering a list of tags.
func RegisterListTagsOptions(g FlagGroup, o *client.ListTagsOptions) {
	b := toBuilder(g, BuilderOptions{})

	registerIntFilterSpec(b, &o.Owner, "owner", "Filter by tag owner ID.")
	registerCharFilterSpec(b, &o.Name, "name", "Filter by tag name.")
}
