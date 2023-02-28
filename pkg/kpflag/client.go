package kpflag

import (
	"github.com/hansmi/paperhooks/internal/kpflagvalue"
	"github.com/hansmi/paperhooks/pkg/client"
)

func RegisterClient(g FlagGroup, f *client.Flags) {
	g.Flag("paperless_url", "Base URL for accessing Paperless.").
		PlaceHolder("URL").
		Envar("PAPERLESS_URL").StringVar(&f.BaseURL)

	g.Flag("paperless_auth_token", "Authentication token for Paperless. Reading the token from a file is preferrable.").
		PlaceHolder("TOKEN").
		Envar("PAPERLESS_AUTH_TOKEN").StringVar(&f.AuthToken)

	g.Flag("paperless_auth_token_file", "File containing authentication token for Paperless.").
		PlaceHolder("PATH").
		Envar("PAPERLESS_AUTH_TOKEN_FILE").StringVar(&f.AuthTokenFile)

	g.Flag("paperless_auth_username", "Username for HTTP basic authentication.").
		PlaceHolder("NAME").
		Envar("PAPERLESS_AUTH_USERNAME").StringVar(&f.AuthUsername)

	g.Flag("paperless_auth_password", "Password for HTTP basic authentication. Reading the password from a file is preferrable.").
		PlaceHolder("PASSWORD").
		Envar("PAPERLESS_AUTH_PASSWORD").StringVar(&f.AuthPassword)

	g.Flag("paperless_auth_password_file", "Username for HTTP basic authentication.").
		PlaceHolder("PATH").
		Envar("PAPERLESS_AUTH_PASSWORD_FILE").StringVar(&f.AuthPasswordFile)

	kpflagvalue.HTTPHeaderVar(
		g.Flag("paperless_header", "HTTP headers to set on all requests to Paperless.").
			PlaceHolder("KEY:VALUE").
			Envar("PAPERLESS_HEADER"), &f.Header)

	g.Flag("paperless_client_debug", "Enable verbose logging messages.").
		BoolVar(&f.DebugMode)
}
