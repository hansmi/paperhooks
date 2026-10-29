package kpflag

import "github.com/hansmi/paperhooks/pkg/client"

// RegisterListStoragePathsOptions adds flags for filtering a list of storage
// paths.
func RegisterListStoragePathsOptions(g FlagGroup, o *client.ListStoragePathsOptions) {
	b := toBuilder(g, BuilderOptions{})

	registerIntFilterSpec(b, &o.Owner, "owner", "Filter by storage path owner ID.")
	registerCharFilterSpec(b, &o.Name, "name", "Filter by storage path name.")
	registerCharFilterSpec(b, &o.Path, "path", "Filter by storage path.")
}
