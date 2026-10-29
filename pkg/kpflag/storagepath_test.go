package kpflag

import (
	"testing"

	"github.com/hansmi/paperhooks/pkg/client"
)

func TestRegisterListStoragePathsOptions(t *testing.T) {
	for _, tc := range []struct {
		name    string
		env     map[string]string
		args    []string
		want    client.ListStoragePathsOptions
		wantErr error
	}{
		{
			name: "defaults",
		},
		{
			name: "values",
			args: []string{
				"--owner_equals=8932487",
				"--name_iexact=start",
				"--path_iendswith=end",
			},
			want: client.ListStoragePathsOptions{
				Owner: client.IntFilterSpec{
					Equals: new(int64(8932487)),
				},
				Name: client.CharFilterSpec{
					EqualsIgnoringCase: new("start"),
				},
				Path: client.CharFilterSpec{
					EndsWithIgnoringCase: new("end"),
				},
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got client.ListStoragePathsOptions

				RegisterListStoragePathsOptions(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
