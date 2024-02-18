package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"
)

func readFile(name string) (string, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimRightFunc(content, unicode.IsSpace)), nil
}

// Flags contains attributes to construct a Paperless client instance. The
// separate "kpflag" package implements bindings for
// [github.com/alecthomas/kingpin/v2].
type Flags struct {
	// Whether to enable verbose log messages.
	DebugMode bool

	// HTTP(S) URL for Paperless.
	BaseURL string

	// Number of concurrent requests allowed to be in flight.
	MaxConcurrentRequests int

	// Authenticate via token.
	AuthToken string

	// Read the authentication token from a file.
	AuthTokenFile string

	// Authenticate via HTTP basic authentication (username and password).
	AuthUsername string
	AuthPassword string

	// Read the password from a file.
	AuthPasswordFile string

	// Authenticate using OpenID Connect (OIDC) ID tokens derived from a Google
	// Cloud Platform service account key file.
	AuthGCPServiceAccountKeyFile string

	// Target audience for OpenID Connect (OIDC) ID tokens. May be left empty,
	// in which case the Paperless URL is used verbatim.
	AuthOIDCIDTokenAudience string

	// HTTP headers to set on all requests.
	Header http.Header

	// Timezone for parsing timestamps without offset.
	ServerTimezone string
}

// This function makes no attempt to deconflict different authentication
// options. Tokens from a files are preferred.
func (f *Flags) buildAuth() (AuthMechanism, error) {
	var err error

	token := f.AuthToken

	if f.AuthTokenFile != "" {
		token, err = readFile(f.AuthTokenFile)
		if err != nil {
			return nil, fmt.Errorf("reading authentication token failed: %w", err)
		}
	}

	if token != "" {
		return &TokenAuth{token}, nil
	}

	if f.AuthUsername != "" {
		password := f.AuthPassword

		if f.AuthPasswordFile != "" {
			password, err = readFile(f.AuthPasswordFile)
			if err != nil {
				return nil, fmt.Errorf("reading password failed: %w", err)
			}
		}

		return &UsernamePasswordAuth{
			Username: f.AuthUsername,
			Password: password,
		}, nil
	}

	if f.AuthGCPServiceAccountKeyFile != "" {
		a, err := GCPServiceAccountKeyAuth{
			KeyFile:  f.AuthGCPServiceAccountKeyFile,
			Audience: f.AuthOIDCIDTokenAudience,
		}.Build()
		if err != nil {
			return nil, fmt.Errorf("GCP service account key authentication: %w", err)
		}

		return a, nil
	}

	return nil, nil
}

// BuildOptions returns the client options derived from flags.
func (f *Flags) BuildOptions() (*Options, error) {
	if f.BaseURL == "" {
		return nil, errors.New("Paperless URL is not specified")
	}

	opts := &Options{
		BaseURL:               f.BaseURL,
		MaxConcurrentRequests: f.MaxConcurrentRequests,
		DebugMode:             f.DebugMode,
		Header:                http.Header{},
		ServerLocation:        time.Local,
	}

	for name, values := range f.Header {
		name = http.CanonicalHeaderKey(name)
		for _, value := range values {
			opts.Header.Add(name, value)
		}
	}

	if auth, err := f.buildAuth(); err != nil {
		return nil, err
	} else {
		opts.Auth = auth
	}

	if f.ServerTimezone != "" {
		if loc, err := time.LoadLocation(f.ServerTimezone); err != nil {
			return nil, err
		} else {
			opts.ServerLocation = loc
		}
	}

	return opts, nil
}

// Build returns a fully configured Paperless client derived from flags.
func (f *Flags) Build() (*Client, error) {
	opts, err := f.BuildOptions()
	if err != nil {
		return nil, err
	}

	return New(*opts), nil
}
