package kpflag

import (
	"testing"

	"github.com/hansmi/paperhooks/pkg/client"
)

func TestRegisterListDocumentTypesOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		args    []string
		want    client.ListDocumentTypesOptions
		wantErr error
	}{
		{
			name: "defaults",
		},
		{
			name: "values",
			args: []string{
				"--no-owner_isnull",
				"--name_icontains=abc",
			},
			want: client.ListDocumentTypesOptions{
				Owner: client.IntFilterSpec{
					IsNull: new(false),
				},
				Name: client.CharFilterSpec{
					ContainsIgnoringCase: new("abc"),
				},
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got client.ListDocumentTypesOptions

				RegisterListDocumentTypesOptions(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
