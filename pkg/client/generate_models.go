//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
)

type modelField struct {
	name      string
	typ       string
	comment   string
	readOnly  bool
	writeOnly bool
}

type model struct {
	name   string
	owned  bool
	fields []modelField
}

func (o model) write(w io.Writer) {
	fmt.Fprintf(w, "type %s struct {\n", strcase.ToCamel(o.name))

	fields := append([]modelField(nil), o.fields...)

	if o.owned {
		fields = append(fields, ownedModelFields...)
	}

	for idx, f := range fields {
		if f.writeOnly {
			continue
		}

		name := strcase.ToCamel(f.name)

		switch f.name {
		case "id":
			name = "ID"
		}

		if f.comment != "" {
			if idx > 0 {
				fmt.Fprintf(w, "\n")
			}
			fmt.Fprintf(w, "  // %s\n", strings.TrimSpace(f.comment))
		}

		fmt.Fprintf(w, "  %s %s `json:%q`\n", name, f.typ, f.name)
	}

	fmt.Fprintf(w, "}\n")

	fieldsStruct := strcase.ToCamel(o.name + "_fields")

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "type %s struct {\n", fieldsStruct)
	fmt.Fprintf(w, "  objectFields\n")
	fmt.Fprintf(w, "}\n")

	fmt.Fprintf(w, "var _ json.Marshaler = (*%s)(nil)\n", fieldsStruct)

	fmt.Fprintf(w, "func New%[1]s() *%[1]s {", fieldsStruct)
	fmt.Fprintf(w, "  return &%s{ objectFields{} }\n", fieldsStruct)
	fmt.Fprintf(w, "}\n")

	for _, f := range fields {
		if f.readOnly {
			continue
		}

		argName := strcase.ToLowerCamel(f.name)

		funcName := fmt.Sprintf("Set%s", strcase.ToCamel(f.name))

		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "// %s sets the %q field.\n", funcName, f.name)
		if f.comment != "" {
			fmt.Fprintf(w, "//\n")
			fmt.Fprintf(w, "// %s\n", strings.TrimSpace(f.comment))
		}
		fmt.Fprintf(w, "func (f *%s) %s(%s %s) *%[1]s {\n", fieldsStruct, funcName, argName, f.typ)
		fmt.Fprintf(w, "  f.set(%q, %s)\n", f.name, argName)
		fmt.Fprintf(w, "  return f\n")
		fmt.Fprintf(w, "}\n")
	}
}

var ownedModelFields = []modelField{
	{
		name:    "owner",
		typ:     "*int64",
		comment: "Object owner; objects without owner can be viewed and edited by all users.",
	},
	{
		name:      "set_permissions",
		typ:       "*ObjectPermissions",
		comment:   "Change object-level permissions.",
		writeOnly: true,
	},
}

var correspondentModel = model{
	name:  "correspondent",
	owned: true,
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "slug", typ: "string", readOnly: true},
		{name: "name", typ: "string"},
		{name: "match", typ: "string"},
		{name: "matching_algorithm", typ: "MatchingAlgorithm"},
		{name: "is_insensitive", typ: "bool"},
		{name: "document_count", typ: "int64", readOnly: true},
		{name: "last_correspondence", typ: "*time.Time", readOnly: true},
	},
}

var customFieldModel = model{
	name:  "customField",
	owned: true,
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "name", typ: "string"},
		{name: "data_type", typ: "string"},
	},
}

var documentModel = model{
	name:  "document",
	owned: true,
	fields: []modelField{
		{name: "id", typ: "int64", comment: "ID of the document.", readOnly: true},
		{name: "title", typ: "string", comment: "Title of the document."},
		{name: "content", typ: "string", comment: "Plain-text content of the document."},
		{name: "tags", typ: "[]int64", comment: "List of tag IDs assigned to this document, or empty list."},
		{name: "document_type", typ: "*int64", comment: "Document type of this document or nil."},
		{name: "correspondent", typ: "*int64", comment: "Correspondent of this document or nil."},
		{name: "storage_path", typ: "*int64", comment: "Storage path of this document or nil."},
		{name: "created", typ: "time.Time", comment: "The date time at which this document was created."},
		{name: "modified", typ: "time.Time", comment: "The date at which this document was last edited in paperless.", readOnly: true},
		{name: "added", typ: "time.Time", comment: "The date at which this document was added to paperless.", readOnly: true},
		{name: "archive_serial_number", typ: "*int64", comment: "The identifier of this document in a physical document archive."},
		{name: "original_file_name", typ: "string", comment: "Verbose filename of the original document.", readOnly: true},
		{name: "archived_file_name", typ: "*string", comment: "Verbose filename of the archived document. Nil if no archived document is available.", readOnly: true},
		{name: "custom_fields", typ: "[]CustomFieldInstance", comment: "Custom fields on the document."},
	},
}

var storagePathModel = model{
	name:  "storagePath",
	owned: true,
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "slug", typ: "string", readOnly: true},
		{name: "name", typ: "string"},
		{name: "match", typ: "string"},
		{name: "matching_algorithm", typ: "MatchingAlgorithm"},
		{name: "is_insensitive", typ: "bool"},
		{name: "document_count", typ: "int64", readOnly: true},
	},
}

var tagModel = model{
	name:  "tag",
	owned: true,
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "slug", typ: "string", readOnly: true},
		{name: "name", typ: "string"},
		{name: "color", typ: "Color"},
		{name: "text_color", typ: "Color"},
		{name: "match", typ: "string"},
		{name: "matching_algorithm", typ: "MatchingAlgorithm"},
		{name: "is_insensitive", typ: "bool"},
		{name: "is_inbox_tag", typ: "bool"},
		{name: "document_count", typ: "int64", readOnly: true},
	},
}

var documentTypeModel = model{
	name:  "documentType",
	owned: true,
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "slug", typ: "string", readOnly: true},
		{name: "name", typ: "string"},
		{name: "match", typ: "string"},
		{name: "matching_algorithm", typ: "MatchingAlgorithm"},
		{name: "is_insensitive", typ: "bool"},
		{name: "document_count", typ: "int64", readOnly: true},
	},
}

var userModel = model{
	name: "user",
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "username", typ: "string"},
		{name: "email", typ: "string"},
		{name: "first_name", typ: "string"},
		{name: "last_name", typ: "string"},
		{name: "is_active", typ: "bool"},
		{name: "is_staff", typ: "bool"},
		{name: "is_superuser", typ: "bool"},
	},
}

var groupModel = model{
	name: "group",
	fields: []modelField{
		{name: "id", typ: "int64", readOnly: true},
		{name: "name", typ: "string"},
	},
}

var statisticsModel = model{
	name: "statistics",
	fields: []modelField{
		{name: "documents_total", typ: "int64", readOnly: true},
		{name: "documents_inbox", typ: "int64", readOnly: true},
		{name: "inbox_tag", typ: "int64", readOnly: true},
		{name: "inbox_tags", typ: "[]int64", readOnly: true},
		{name: "document_file_type_counts", typ: "[]DocumentFileType", readOnly: true},
		{name: "character_count", typ: "int64", readOnly: true},
		{name: "tag_count", typ: "int64", readOnly: true},
		{name: "correspondent_count", typ: "int64", readOnly: true},
		{name: "document_type_count", typ: "int64", readOnly: true},
		{name: "storage_path_count", typ: "int64", readOnly: true},
		{name: "current_asn", typ: "int64", readOnly: true},
	},
}

func main() {
	outputFile := flag.String("output", "", "Destination file")

	flag.Parse()

	var buf bytes.Buffer

	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(&buf, "// Code generated by %q; DO NOT EDIT.\n",
		strings.Join(append([]string{filepath.Base(exe)}, os.Args[1:]...), " "))
	buf.WriteString("\n")
	buf.WriteString("package client\n")

	imports := []string{
		"encoding/json",
		"time",
	}

	sort.Strings(imports)

	for _, i := range imports {
		fmt.Fprintf(&buf, "import %q\n", i)
	}

	models := []model{
		correspondentModel,
		customFieldModel,
		documentModel,
		documentTypeModel,
		storagePathModel,
		tagModel,
		userModel,
		groupModel,
		statisticsModel,
	}

	sort.Slice(models, func(a, b int) bool {
		return models[a].name < models[b].name
	})

	for _, i := range models {
		i.write(&buf)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Formatting code failed: %v\n%s", err, buf.String())
	}

	if *outputFile == "" || *outputFile == "-" {
		os.Stdout.Write(formatted)
	} else if err := os.WriteFile(*outputFile, formatted, 0o644); err != nil {
		log.Fatal("Writing output failed: %v", err)
	}
}
