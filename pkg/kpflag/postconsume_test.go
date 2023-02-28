package kpflag

import (
	"net/url"
	"testing"
	"time"

	"github.com/hansmi/paperhooks/pkg/postconsume"
)

func TestRegisterPostConsume(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
		args []string
		want postconsume.Flags
	}{
		{
			name: "defaults",
		},
		{
			name: "flags",
			args: []string{
				"--document_id=1234",
				"--document_file_name=filename",
				"--document_created=2000-01-01",
				"--document_modified=2010-02-03T12:34:56-05:00",
				"--document_added=2005-11-22 09:08",
				"--document_source_path=source",
				"--document_archive_path=archive",
				"--document_thumbnail_path=thumbnail",
				"--document_download_url=/down/load",
				"--document_thumbnail_url=/thumb/nail",
				"--document_correspondent=Mail Ltd.",
				"--document_tags=x,y,z",
				"--document_tags=more, tags",
				"--document_original_filename=original",
			},
			want: postconsume.Flags{
				DocumentID:               1234,
				DocumentFilename:         "filename",
				DocumentCreated:          time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
				DocumentModified:         time.Date(2010, time.February, 3, 12+5, 34, 56, 0, time.UTC),
				DocumentAdded:            time.Date(2005, time.November, 22, 9, 8, 0, 0, time.Local),
				DocumentSourcePath:       "source",
				DocumentArchivePath:      "archive",
				DocumentThumbnailPath:    "thumbnail",
				DocumentDownloadURL:      &url.URL{Path: "/down/load"},
				DocumentThumbnailURL:     &url.URL{Path: "/thumb/nail"},
				DocumentCorrespondent:    "Mail Ltd.",
				DocumentTags:             []string{"x", "y", "z", "more", "tags"},
				DocumentOriginalFilename: "original",
			},
		},
		{
			name: "env",
			env: map[string]string{
				"DOCUMENT_ID":                "14903",
				"DOCUMENT_FILE_NAME":         "envfile",
				"DOCUMENT_SOURCE_PATH":       "envsource",
				"DOCUMENT_ARCHIVE_PATH":      "envarchive",
				"DOCUMENT_THUMBNAIL_PATH":    "envthumbnail",
				"DOCUMENT_TAGS":              "foo, bar, baz",
				"DOCUMENT_ORIGINAL_FILENAME": "envorig",
			},
			want: postconsume.Flags{
				DocumentID:               14903,
				DocumentFilename:         "envfile",
				DocumentSourcePath:       "envsource",
				DocumentArchivePath:      "envarchive",
				DocumentThumbnailPath:    "envthumbnail",
				DocumentTags:             []string{"foo", "bar", "baz"},
				DocumentOriginalFilename: "envorig",
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got postconsume.Flags

				RegisterPostConsume(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
