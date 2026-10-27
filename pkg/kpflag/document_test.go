package kpflag

import (
	"testing"

	"github.com/hansmi/paperhooks/pkg/client"
)

func TestRegisterListDocumentsOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		args    []string
		want    client.ListDocumentsOptions
		wantErr error
	}{
		{
			name: "defaults",
		},
		{
			name: "values",
			args: []string{
				"--owner_isnull",
				"--no-storage_path_isnull",
				"--tags_id=1234",
				"--content_icontains=hello",
			},
			want: client.ListDocumentsOptions{
				Owner: client.IntFilterSpec{
					IsNull: new(true),
				},
				Content: client.CharFilterSpec{
					ContainsIgnoringCase: new("hello"),
				},
				Tags: client.ForeignKeyFilterSpec{
					ID: new(int64(1234)),
				},
				StoragePath: client.ForeignKeyFilterSpec{
					IsNull: new(false),
				},
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got client.ListDocumentsOptions

				RegisterListDocumentsOptions(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
