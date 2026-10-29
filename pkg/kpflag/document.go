package kpflag

import "github.com/hansmi/paperhooks/pkg/client"

// RegisterListDocumentsOptions adds flags for filtering a list of documents.
func RegisterListDocumentsOptions(g FlagGroup, o *client.ListDocumentsOptions) {
	b := toBuilder(g, BuilderOptions{})

	registerIntFilterSpec(b, &o.Owner, "owner", "Filter by document owner ID.")
	registerCharFilterSpec(b, &o.Title, "title", "Filter by document title.")
	registerCharFilterSpec(b, &o.Content, "content", "Filter by document content.")
	registerIntFilterSpec(b, &o.ArchiveSerialNumber, "archive_serial_number", "Filter by document archive serial number.")
	registerDateTimeFilterSpec(b, &o.Created, "created", "Filter by document creation time.")
	registerDateTimeFilterSpec(b, &o.Added, "added", "Filter by document import time.")
	registerDateTimeFilterSpec(b, &o.Modified, "modified", "Filter by document modification time.")
	registerForeignKeyFilterSpec(b, &o.Correspondent, "correspondent", "Filter by document correspondent.")
	registerForeignKeyFilterSpec(b, &o.Tags, "tags", "Filter by document tags.")
	registerForeignKeyFilterSpec(b, &o.DocumentType, "document_type", "Filter by document type.")
	registerForeignKeyFilterSpec(b, &o.StoragePath, "storage_path", "Filter by document storage path.")
}
