package kpflag

import (
	"testing"

	"github.com/hansmi/paperhooks/pkg/client"
)

func TestRegisterListCorrespondentsOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		args    []string
		want    client.ListCorrespondentsOptions
		wantErr error
	}{
		{
			name: "defaults",
		},
		{
			name: "values",
			args: []string{
				"--owner_equals=999",
				"--name_iexact=foo",
			},
			want: client.ListCorrespondentsOptions{
				Owner: client.IntFilterSpec{
					Equals: new(int64(999)),
				},
				Name: client.CharFilterSpec{
					EqualsIgnoringCase: new("foo"),
				},
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got client.ListCorrespondentsOptions

				RegisterListCorrespondentsOptions(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
