package testutil

import (
	"os"
	"strings"
	"testing"
)

// RestoreEnv restores all environment variables at the end of a test.
func RestoreEnv(t *testing.T) {
	t.Helper()

	previous := os.Environ()

	t.Cleanup(func() {
		os.Clearenv()

		for _, i := range previous {
			parts := strings.SplitN(i, "=", 2)

			if err := os.Setenv(parts[0], parts[1]); err != nil {
				t.Errorf("Setenv(%q, %q) failed: %v", parts[0], parts[1], err)
			}
		}
	})
}
