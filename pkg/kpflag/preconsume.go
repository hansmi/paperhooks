package kpflag

import "github.com/hansmi/paperhooks/pkg/preconsume"

// RegisterPreConsume adds flags capturing Paperless-ngx pre-consumption
// command information.
func RegisterPreConsume(g FlagGroup, f *preconsume.Flags) {
	g.Flag("document_source_path", "Original path of the consumed document.").
		PlaceHolder("PATH").
		Envar("DOCUMENT_SOURCE_PATH").StringVar(&f.DocumentSourcePath)

	g.Flag("document_working_path", "Path to a copy of the original that consumption will work on.").
		PlaceHolder("PATH").
		Envar("DOCUMENT_WORKING_PATH").StringVar(&f.DocumentWorkingPath)
}
