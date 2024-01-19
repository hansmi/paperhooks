package kpflag

import "github.com/hansmi/paperhooks/pkg/preconsume"

// RegisterPreConsume adds flags capturing Paperless-ngx pre-consumption
// command information.
func RegisterPreConsume(g FlagGroup, f *preconsume.Flags) {
	b := builder{g}

	b.flag("document_source_path", "Original path of the consumed document.").
		PlaceHolder("PATH").
		StringVar(&f.DocumentSourcePath)

	b.flag("document_working_path", "Path to a copy of the original that consumption will work on.").
		PlaceHolder("PATH").
		StringVar(&f.DocumentWorkingPath)
}
