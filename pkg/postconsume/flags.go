package postconsume

import (
	"net/url"
	"time"
)

// Flags contains attributes for storing the environment variables given to
// a [post-consumption script]. The separate "kpflag" package implements
// bindings for [github.com/alecthomas/kingpin/v2].
//
// [post-consumption script]: https://docs.paperless-ngx.com/advanced_usage/#post-consume-script
type Flags struct {
	// Primary database key of the document.
	DocumentID int64

	// Formatted filename, not including paths.
	DocumentFilename string

	// Date and time when the document was created.
	DocumentCreated time.Time

	// Date and time when the document was last modified.
	DocumentModified time.Time

	// Date and time when the document was added.
	DocumentAdded time.Time

	// Path to the original document file.
	DocumentSourcePath string

	// Path to the generated archive file.
	DocumentArchivePath string

	// Path to the generated thumbnail image.
	DocumentThumbnailPath string

	// URL for document download.
	DocumentDownloadURL *url.URL

	// URL for the document thumbnail image.
	DocumentThumbnailURL *url.URL

	// Name of the assigned correspondent.
	DocumentCorrespondent string

	// Names of tags applied to document.
	DocumentTags []string

	// Filename of original document.
	DocumentOriginalFilename string
}
