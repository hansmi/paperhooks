package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// GCPServiceAccountKeyAuth uses a Google Cloud Platform service account key
// file to authenticate against an OAuth 2.0-protected Paperless instance using
// the two-legged JWT flow.
//
// The service account key is used to request OpenID Connect (OIDC) ID tokens
// from the Google OAuth 2.0 API. The ID tokens are in turn used for all
// Paperless API requests.
//
// References:
//
//   - https://cloud.google.com/iam/docs/service-account-creds
//   - https://openid.net/specs/openid-connect-core-1_0.html
type GCPServiceAccountKeyAuth struct {
	// Path to a file containing the service account key.
	KeyFile string

	// Service account key in JSON format.
	Key []byte

	// Audience to request for the ID token (case-sensitive). If empty the
	// Paperless URL is used verbatim.
	Audience string

	// Custom HTTP client for requesting tokens.
	HTTPClient *http.Client
}

func (a GCPServiceAccountKeyAuth) Build() (AuthMechanism, error) {
	key := a.Key

	if len(key) == 0 {
		if a.KeyFile == "" {
			return nil, fmt.Errorf("%w: missing key or key path", os.ErrInvalid)
		}

		if content, err := os.ReadFile(a.KeyFile); err != nil {
			return nil, fmt.Errorf("reading service account key: %w", err)
		} else {
			key = content
		}
	}

	config, err := google.JWTConfigFromJSON(key)
	if err != nil {
		return nil, fmt.Errorf("building JWT config from service account key: %w", err)
	}

	if config.PrivateClaims == nil {
		config.PrivateClaims = map[string]any{}
	}

	config.UseIDToken = true

	return &gcpServiceAccountKeyAuthImpl{
		audience:   a.Audience,
		httpClient: a.HTTPClient,
		config:     config,
	}, nil
}

type gcpServiceAccountKeyAuthImpl struct {
	mu         sync.Mutex
	audience   string
	httpClient *http.Client
	config     *jwt.Config
}

var _ AuthMechanism = (*gcpServiceAccountKeyAuthImpl)(nil)

func (o *gcpServiceAccountKeyAuthImpl) authenticate(clientOpts Options, c *resty.Client) {
	ctx := context.Background()

	o.mu.Lock()
	defer o.mu.Unlock()

	audience := o.audience

	if audience == "" {
		audience = clientOpts.BaseURL
	}

	o.config.PrivateClaims["target_audience"] = audience

	if o.httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, o.httpClient)
	}

	c.SetTransport(&oauth2.Transport{
		Base:   c.GetClient().Transport,
		Source: oauth2.ReuseTokenSource(nil, o.config.TokenSource(ctx)),
	})
}
