package kpflag

import (
	"github.com/hansmi/paperhooks/internal/kpflagvalue"
	"github.com/hansmi/paperhooks/pkg/postconsume"
)

// RegisterPostConsume adds flags capturing Paperless-ngx post-consumption
// command information.
func RegisterPostConsume(g FlagGroup, f *postconsume.Flags) {
	b := builder{g}

	b.flag("document_id", "Primary database key of the document.").
		PlaceHolder("INT").
		Int64Var(&f.DocumentID)

	b.flag("document_file_name", "Formatted filename, not including paths.").
		PlaceHolder("NAME").
		StringVar(&f.DocumentFilename)

	kpflagvalue.TimeVar(
		b.flag("document_created", "Date and time when the document was created."),
		&f.DocumentCreated)

	kpflagvalue.TimeVar(
		b.flag("document_modified", "Date and time when the document was last modified."),
		&f.DocumentModified)

	kpflagvalue.TimeVar(
		b.flag("document_added", "Date and time when the document was added."),
		&f.DocumentAdded)

	b.flag("document_source_path", "Path to the original document file.").
		PlaceHolder("PATH").
		StringVar(&f.DocumentSourcePath)

	b.flag("document_archive_path", "Path to the generated archive file (if any).").
		PlaceHolder("PATH").
		StringVar(&f.DocumentArchivePath)

	b.flag("document_thumbnail_path", "Path to the generated thumbnail image.").
		PlaceHolder("PATH").
		StringVar(&f.DocumentThumbnailPath)

	b.flag("document_download_url", "URL for document download.").
		PlaceHolder("URL").
		URLVar(&f.DocumentDownloadURL)

	b.flag("document_thumbnail_url", "URL for the document thumbnail image.").
		PlaceHolder("URL").
		URLVar(&f.DocumentThumbnailURL)

	b.flag("document_correspondent", "Assigned correspondent (if any).").
		StringVar(&f.DocumentCorrespondent)

	kpflagvalue.CommaSeparatedStringsVar(
		b.flag("document_tags", "Comma separated list of tags applied (if any).").
			PlaceHolder("TAGS"),
		&f.DocumentTags)

	b.flag("document_original_filename", "Filename of original document.").
		PlaceHolder("NAME").
		StringVar(&f.DocumentOriginalFilename)
}
