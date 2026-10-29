package kpflag

import (
	"testing"

	"github.com/hansmi/paperhooks/pkg/client"
)

func TestRegisterListTagsOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		args    []string
		want    client.ListTagsOptions
		wantErr error
	}{
		{
			name: "defaults",
		},
		{
			name: "values",
			args: []string{
				"--owner_equals=12124",
				"--name_istartswith=9ef3a65fb25",
			},
			want: client.ListTagsOptions{
				Owner: client.IntFilterSpec{
					Equals: new(int64(12124)),
				},
				Name: client.CharFilterSpec{
					StartsWithIgnoringCase: new("9ef3a65fb25"),
				},
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got client.ListTagsOptions

				RegisterListTagsOptions(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
