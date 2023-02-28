package kpflag

import (
	"os"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/paperhooks/internal/testutil"
)

type flagParseTest struct {
	name     string
	register func(t *testing.T, g FlagGroup) any
	args     []string
	env      map[string]string
	want     any
}

func (tc flagParseTest) run(t *testing.T) {
	t.Helper()
	t.Run(tc.name, func(t *testing.T) {
		testutil.RestoreEnv(t)

		os.Clearenv()

		app := kingpin.New(tc.name, "")

		got := tc.register(t, app)

		testutil.Setenv(t, tc.env)

		if _, err := app.Parse(tc.args); err != nil {
			t.Errorf("Parsing arguments failed: %v", err)
		} else if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Parsed flags diff (-want +got):\n%s", diff)
		}
	})
}
