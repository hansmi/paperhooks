package kpflag

import "github.com/hansmi/paperhooks/pkg/client"

// RegisterListCorrespondentsOptions adds flags for filtering a list of documents.
func RegisterListCorrespondentsOptions(g FlagGroup, o *client.ListCorrespondentsOptions) {
	b := toBuilder(g, BuilderOptions{})

	registerIntFilterSpec(b, &o.Owner, "owner", "Filter by correspondent owner ID.")
	registerCharFilterSpec(b, &o.Name, "name", "Filter by correspondent name.")
}
