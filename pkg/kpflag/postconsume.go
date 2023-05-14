package kpflag

import (
	"github.com/hansmi/paperhooks/internal/kpflagvalue"
	"github.com/hansmi/paperhooks/pkg/postconsume"
)

// RegisterPostConsume adds flags capturing Paperless-ngx post-consumption
// command information.
func RegisterPostConsume(g FlagGroup, f *postconsume.Flags) {
	g.Flag("document_id", "Primary database key of the document.").
		PlaceHolder("INT").
		Envar("DOCUMENT_ID").Int64Var(&f.DocumentID)

	g.Flag("document_file_name", "Formatted filename, not including paths.").
		PlaceHolder("NAME").
		Envar("DOCUMENT_FILE_NAME").StringVar(&f.DocumentFilename)

	kpflagvalue.TimeVar(
		g.Flag("document_created", "Date and time when the document was created.").
			Envar("DOCUMENT_CREATED"), &f.DocumentCreated)

	kpflagvalue.TimeVar(
		g.Flag("document_modified", "Date and time when the document was last modified.").
			Envar("DOCUMENT_MODIFIED"), &f.DocumentModified)

	kpflagvalue.TimeVar(
		g.Flag("document_added", "Date and time when the document was added.").
			Envar("DOCUMENT_ADDED"), &f.DocumentAdded)

	g.Flag("document_source_path", "Path to the original document file.").
		PlaceHolder("PATH").
		Envar("DOCUMENT_SOURCE_PATH").StringVar(&f.DocumentSourcePath)

	g.Flag("document_archive_path", "Path to the generated archive file (if any).").
		PlaceHolder("PATH").
		Envar("DOCUMENT_ARCHIVE_PATH").StringVar(&f.DocumentArchivePath)

	g.Flag("document_thumbnail_path", "Path to the generated thumbnail image.").
		PlaceHolder("PATH").
		Envar("DOCUMENT_THUMBNAIL_PATH").StringVar(&f.DocumentThumbnailPath)

	g.Flag("document_download_url", "URL for document download.").
		PlaceHolder("URL").
		Envar("DOCUMENT_DOWNLOAD_URL").URLVar(&f.DocumentDownloadURL)

	g.Flag("document_thumbnail_url", "URL for the document thumbnail image.").
		PlaceHolder("URL").
		Envar("DOCUMENT_THUMBNAIL_URL").URLVar(&f.DocumentThumbnailURL)

	g.Flag("document_correspondent", "Assigned correspondent (if any).").
		Envar("DOCUMENT_CORRESPONDENT").StringVar(&f.DocumentCorrespondent)

	kpflagvalue.CommaSeparatedStringsVar(
		g.Flag("document_tags", "Comma separated list of tags applied (if any).").
			PlaceHolder("TAGS").
			Envar("DOCUMENT_TAGS"), &f.DocumentTags)

	g.Flag("document_original_filename", "Filename of original document.").
		PlaceHolder("NAME").
		Envar("DOCUMENT_ORIGINAL_FILENAME").StringVar(&f.DocumentOriginalFilename)
}
