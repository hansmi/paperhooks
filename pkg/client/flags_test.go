package client

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/paperhooks/internal/testutil"
)

func TestRegisterClient(t *testing.T) {
	for _, tc := range []struct {
		name    string
		flags   Flags
		want    Options
		wantErr error
	}{
		{
			name:    "defaults",
			wantErr: cmpopts.AnyError,
		},
		{
			name: "token auth",
			flags: Flags{
				BaseURL:   "http://localhost:1234",
				AuthToken: "abcdef22379",
				Header: http.Header{
					"x-header": []string{"value"},
					"another":  []string{"value2"},
				},
			},
			want: Options{
				BaseURL: "http://localhost:1234",
				Header: http.Header{
					"X-Header": []string{"value"},
					"Another":  []string{"value2"},
				},
				Auth: &TokenAuth{"abcdef22379"},
			},
		},
		{
			name: "token file",
			flags: Flags{
				BaseURL:       "http://localhost:1234/tokenfile",
				AuthToken:     "mytoken",
				AuthTokenFile: testutil.MustWriteFile(t, filepath.Join(t.TempDir(), "file.txt"), "content\n"),
				Header: http.Header{
					"x-header": []string{"foobar"},
				},
			},
			want: Options{
				BaseURL: "http://localhost:1234/tokenfile",
				Header: http.Header{
					"X-Header": []string{"foobar"},
				},
				Auth: &TokenAuth{"content"},
			},
		},
		{
			name: "password auth",
			flags: Flags{
				DebugMode:    true,
				BaseURL:      "http://localhost:9999/pw",
				AuthUsername: "admin",
				AuthPassword: "password",
			},
			want: Options{
				BaseURL: "http://localhost:9999/pw",
				Auth: &UsernamePasswordAuth{
					Username: "admin",
					Password: "password",
				},
				DebugMode: true,
			},
		},
		{
			name: "token file not found",
			flags: Flags{
				BaseURL:       "http://localhost/notfound",
				AuthTokenFile: filepath.Join(t.TempDir(), "missing"),
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "password file not found",
			flags: Flags{
				BaseURL:          "http://localhost/notfound",
				AuthUsername:     "user",
				AuthPasswordFile: filepath.Join(t.TempDir(), "missing"),
			},
			wantErr: os.ErrNotExist,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.flags.BuildOptions()

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("BuildOptions() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, *got, cmpopts.EquateEmpty(), cmp.AllowUnexported(Options{})); diff != "" {
					t.Errorf("Options diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
