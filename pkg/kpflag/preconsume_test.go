package kpflag

import (
	"testing"

	"github.com/hansmi/paperhooks/pkg/preconsume"
)

func TestRegisterPreConsume(t *testing.T) {
	for _, tc := range []struct {
		name string
		env  map[string]string
		args []string
		want preconsume.Flags
	}{
		{
			name: "defaults",
		},
		{
			name: "flags",
			args: []string{
				"--document_source_path=source/path",
				"--document_working_path=working/path",
			},
			want: preconsume.Flags{
				DocumentSourcePath:  "source/path",
				DocumentWorkingPath: "working/path",
			},
		},
		{
			name: "env",
			env: map[string]string{
				"DOCUMENT_SOURCE_PATH":  "env/source",
				"DOCUMENT_WORKING_PATH": "env/working",
			},
			want: preconsume.Flags{
				DocumentSourcePath:  "env/source",
				DocumentWorkingPath: "env/working",
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got preconsume.Flags

				RegisterPreConsume(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
