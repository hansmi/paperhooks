package kpflag

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/hansmi/paperhooks/internal/testutil"
	"github.com/hansmi/paperhooks/pkg/client"
)

func TestRegisterClient(t *testing.T) {
	tokenfile := testutil.MustWriteFile(t, filepath.Join(t.TempDir(), "token.txt"), "filetoken\n")

	for _, tc := range []struct {
		name    string
		env     map[string]string
		args    []string
		want    client.Flags
		wantErr error
	}{
		{
			name: "defaults",
		},
		{
			name: "token auth",
			args: []string{
				"--paperless_url=http://localhost:1234",
				"--paperless_auth_token=abcdef22379",
				"--paperless_header=x-header:value",
				"--paperless_header=another:value2",
			},
			want: client.Flags{
				BaseURL: "http://localhost:1234",
				Header: http.Header{
					"X-Header": []string{"value"},
					"Another":  []string{"value2"},
				},
				AuthToken: "abcdef22379",
			},
		},
		{
			name: "env with token file",
			env: map[string]string{
				"PAPERLESS_URL":             "http://localhost:1234/env",
				"PAPERLESS_AUTH_TOKEN":      "envtoken",
				"PAPERLESS_AUTH_TOKEN_FILE": tokenfile,
				"PAPERLESS_HEADER":          "x-header:foobar",
				"PAPERLESS_CLIENT_DEBUG":    "1",
			},
			want: client.Flags{
				BaseURL: "http://localhost:1234/env",
				Header: http.Header{
					"X-Header": []string{"foobar"},
				},
				AuthToken:     "envtoken",
				AuthTokenFile: tokenfile,
				DebugMode:     true,
			},
		},
		{
			name: "env with password",
			env: map[string]string{
				"PAPERLESS_URL":           "http://localhost:9999/pw",
				"PAPERLESS_AUTH_USERNAME": "admin",
				"PAPERLESS_AUTH_PASSWORD": "password",
			},
			args: []string{
				"--paperless_client_debug",
			},
			want: client.Flags{
				BaseURL:      "http://localhost:9999/pw",
				AuthUsername: "admin",
				AuthPassword: "password",
				DebugMode:    true,
			},
		},
		{
			name: "server timezone",
			env: map[string]string{
				"PAPERLESS_SERVER_TIMEZONE": "Australia/Sydney",
			},
			want: client.Flags{
				ServerTimezone: "Australia/Sydney",
			},
		},
		{
			name: "max concurrent requests",
			args: []string{
				"--paperless_max_concurrent_requests=123",
			},
			want: client.Flags{
				MaxConcurrentRequests: 123,
			},
		},
	} {
		flagParseTest{
			name: tc.name,
			register: func(t *testing.T, g FlagGroup) any {
				var got client.Flags

				RegisterClient(g, &got)

				return &got
			},
			args: tc.args,
			env:  tc.env,
			want: &tc.want,
		}.run(t)
	}
}
