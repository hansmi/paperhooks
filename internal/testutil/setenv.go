package testutil

import (
	"os"
	"testing"
)

func Setenv(t *testing.T, env map[string]string) {
	t.Helper()

	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			t.Errorf("Setenv(%q, %q) failed: %v", key, value, err)
		}
	}
}
