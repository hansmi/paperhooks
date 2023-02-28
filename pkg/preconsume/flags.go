package preconsume

// Flags contains attributes for storing the environment variables given to
// a [pre-consumption script]. The separate "kpflag" package implements
// bindings for [github.com/alecthomas/kingpin/v2].
//
// [pre-consumption script]: https://docs.paperless-ngx.com/advanced_usage/#pre-consume-script
type Flags struct {
	// Original path of the consumed document.
	DocumentSourcePath string

	// Path to a copy of the original that consumption will work on.
	DocumentWorkingPath string
}
