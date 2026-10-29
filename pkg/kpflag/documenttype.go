package kpflag

import "github.com/hansmi/paperhooks/pkg/client"

// RegisterListDocumentTypesOptions adds flags for filtering a list of document
// types.
func RegisterListDocumentTypesOptions(g FlagGroup, o *client.ListDocumentTypesOptions) {
	b := toBuilder(g, BuilderOptions{})

	registerIntFilterSpec(b, &o.Owner, "owner", "Filter by document type owner ID.")
	registerCharFilterSpec(b, &o.Name, "name", "Filter by document type name.")
}
