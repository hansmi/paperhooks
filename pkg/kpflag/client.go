package kpflag

import (
	"fmt"
	"time"

	"github.com/hansmi/paperhooks/internal/kpflagvalue"
	"github.com/hansmi/paperhooks/pkg/client"
)

// RegisterClient adds flags for creating a Paperless-ngx API client.
func RegisterClient(g FlagGroup, f *client.Flags) {
	b := builder{g}

	b.flag("paperless_url", "Base URL for accessing Paperless.").
		PlaceHolder("URL").
		StringVar(&f.BaseURL)

	b.flag("paperless_max_concurrent_requests", "Number of requests allowed to be in flight at the same time. Defaults to zero (disabled).").
		PlaceHolder("NUM").
		IntVar(&f.MaxConcurrentRequests)

	b.flag("paperless_auth_token", "Authentication token for Paperless. Reading the token from a file is preferrable.").
		PlaceHolder("TOKEN").
		StringVar(&f.AuthToken)

	b.flag("paperless_auth_token_file", "File containing authentication token for Paperless.").
		PlaceHolder("PATH").
		StringVar(&f.AuthTokenFile)

	b.flag("paperless_auth_username", "Username for HTTP basic authentication.").
		PlaceHolder("NAME").
		StringVar(&f.AuthUsername)

	b.flag("paperless_auth_password", "Password for HTTP basic authentication. Reading the password from a file is preferrable.").
		PlaceHolder("PASSWORD").
		StringVar(&f.AuthPassword)

	b.flag("paperless_auth_password_file", "Username for HTTP basic authentication.").
		PlaceHolder("PATH").
		StringVar(&f.AuthPasswordFile)

	b.flag("paperless_auth_gcp_service_account_key_file", "Authenticate using OpenID Connect (OIDC) ID tokens derived from a Google Cloud Platform service account key file.").
		PlaceHolder("PATH").
		StringVar(&f.AuthGCPServiceAccountKeyFile)

	b.flag("paperless_auth_oidc_id_token_audience", "Target audience for OpenID Connect (OIDC) ID tokens. Defaults to the base URL.").
		PlaceHolder("STRING").
		StringVar(&f.AuthOIDCIDTokenAudience)

	kpflagvalue.HTTPHeaderVar(
		b.flag("paperless_header", "HTTP headers to set on all requests to Paperless.").
			PlaceHolder("KEY:VALUE"),
		&f.Header)

	b.flag("paperless_server_timezone", fmt.Sprintf("Timezone for parsing timestamps. Defaults to %q.", time.Local.String())).
		PlaceHolder("AREA/LOCATION").
		StringVar(&f.ServerTimezone)

	b.flag("paperless_client_debug", "Enable verbose logging messages.").
		BoolVar(&f.DebugMode)
}
