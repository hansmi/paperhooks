package client

import (
	"crypto/x509"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/paperhooks/internal/testutil"
)

func TestFlagsBuild(t *testing.T) {
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
				Auth:           &TokenAuth{"abcdef22379"},
				ServerLocation: time.Local,
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
				Auth:           &TokenAuth{"content"},
				ServerLocation: time.Local,
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
				DebugMode:      true,
				ServerLocation: time.Local,
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
		{
			name: "explicit timezone",
			flags: Flags{
				BaseURL:        "http://localhost/timezone",
				ServerTimezone: "UTC",
			},
			want: Options{
				BaseURL:        "http://localhost/timezone",
				ServerLocation: time.UTC,
			},
		},
		{
			name: "trusted CA empty",
			flags: Flags{
				BaseURL: "http://localhost/rootca/empty",
				TrustedRootCAFiles: []string{
					testutil.MustWriteFile(t, filepath.Join(t.TempDir(), "file.txt"), ""),
				},
			},
			want: Options{
				BaseURL:        "http://localhost/rootca/empty",
				ServerLocation: time.Local,
				TrustedRootCAs: x509.NewCertPool(),
			},
		},
		{
			name: "trusted CA",
			flags: Flags{
				BaseURL: "http://localhost/rootca",
				TrustedRootCAFiles: []string{
					testutil.MustWriteFile(t, filepath.Join(t.TempDir(), "file.txt"), fakeCertPEM),
				},
			},
			want: Options{
				BaseURL:        "http://localhost/rootca",
				ServerLocation: time.Local,
				TrustedRootCAs: newFakeCertPool(t),
			},
		},
		{
			name: "trusted CA file not found",
			flags: Flags{
				BaseURL: "http://localhost/rootca/notfound",
				TrustedRootCAFiles: []string{
					filepath.Join(t.TempDir(), "missing"),
				},
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
				if diff := cmp.Diff(tc.want, *got, cmpopts.EquateEmpty(), cmp.AllowUnexported(Options{}), testutil.EquateTimeLocation()); diff != "" {
					t.Errorf("Options diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
