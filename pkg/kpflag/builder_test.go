package kpflag

import (
	"os"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/paperhooks/internal/testutil"
)

func TestBuilder(t *testing.T) {
	for _, tc := range []struct {
		name     string
		env      map[string]string
		defaults BuilderOptions
		setup    func(*testing.T, *kingpin.Application) FlagGroup
		args     []string
		want     string
	}{
		{
			name: "defaults",
		},
		{
			name: "env empty",
			defaults: BuilderOptions{
				EnableEnvVars: true,
			},
		},
		{
			name: "env",
			env: map[string]string{
				"MYFLAG": "foobar",
			},
			defaults: BuilderOptions{
				EnableEnvVars: true,
			},
			want: "foobar",
		},
		{
			name: "env with prefix",
			env: map[string]string{
				"ABC_MYFLAG": "foobar",
			},
			defaults: BuilderOptions{
				Prefix:        "abc_",
				EnableEnvVars: true,
			},
			want: "foobar",
		},
		{
			name: "env and arg with prefix",
			env: map[string]string{
				"ABC_MYFLAG": "unused",
			},
			defaults: BuilderOptions{
				Prefix:        "abc_",
				EnableEnvVars: true,
			},
			args: []string{"--abc_myflag=arg"},
			want: "arg",
		},
		{
			name: "arg with prefix",
			defaults: BuilderOptions{
				Prefix: "test_",
			},
			args: []string{"--test_myflag=123"},
			want: "123",
		},
		{
			name: "custom",
			setup: func(_ *testing.T, app *kingpin.Application) FlagGroup {
				return NewBuilder(app, BuilderOptions{
					Prefix: "custom_",
				})
			},
			defaults: BuilderOptions{
				Prefix: "unused",
			},
			args: []string{"--custom_myflag=value"},
			want: "value",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testutil.RestoreEnv(t)

			os.Clearenv()

			testutil.Setenv(t, tc.env)

			if tc.setup == nil {
				tc.setup = func(_ *testing.T, app *kingpin.Application) FlagGroup {
					return app
				}
			}

			app := kingpin.New(tc.name, "")

			b := toBuilder(tc.setup(t, app), tc.defaults)

			var got string

			flag := b.Flag("myflag", "myhelp")
			flag.StringVar(&got)

			if _, err := app.Parse(tc.args); err != nil {
				t.Errorf("Parsing arguments failed: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Parsed value diff (-want +got):\n%s", diff)
			}
		})
	}
}
