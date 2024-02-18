package client

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// Options for constructing a Paperless client.
type Options struct {
	// Paperless URL. May include a path.
	BaseURL string

	// API authentication.
	Auth AuthMechanism

	// Enable debug mode with many details logged.
	DebugMode bool

	// Logger for writing log messages. If debug mode is enabled and no logger
	// is configured all messages are written to standard library's default
	// logger (log.Default()).
	Logger Logger

	// HTTP headers to set on all requests.
	Header http.Header

	// Server's timezone for parsing timestamps without explicit offset.
	// Defaults to [time.Local].
	ServerLocation *time.Location

	// Override the default HTTP transport.
	transport http.RoundTripper
}

type Client struct {
	logger Logger
	loc    *time.Location
	r      *resty.Client
}

// New creates a new client instance.
func New(opts Options) *Client {
	if opts.Logger == nil {
		if opts.DebugMode {
			opts.Logger = &wrappedStdLogger{log.Default()}
		} else {
			opts.Logger = &discardLogger{}
		}
	}

	if opts.ServerLocation == nil {
		opts.ServerLocation = time.Local
	}

	r := resty.New().
		SetDebug(opts.DebugMode).
		SetLogger(&prefixLogger{
			wrapped: opts.Logger,
			prefix:  "Resty: ",
		}).
		SetDisableWarn(true).
		SetBaseURL(opts.BaseURL).
		SetHeader("Accept", "application/json; version=2").
		SetRedirectPolicy(resty.NoRedirectPolicy())

	if opts.transport != nil {
		r.SetTransport(opts.transport)
	}

	if opts.Auth != nil {
		opts.Auth.authenticate(opts, r)
	}

	if len(opts.Header) > 0 {
		r.SetPreRequestHook(func(_ *resty.Client, req *http.Request) error {
			for name, values := range opts.Header {
				req.Header[http.CanonicalHeaderKey(name)] = values
			}

			return nil
		})
	}

	return &Client{
		logger: opts.Logger,
		loc:    opts.ServerLocation,
		r:      r,
	}
}

func (c *Client) newRequest(ctx context.Context) *resty.Request {
	return c.r.R().
		SetContext(ctx).
		SetError(requestError{}).
		ExpectContentType("application/json")
}

type Response struct {
	*http.Response

	// Token for fetching next page in paginated result sets.
	NextPage *PageToken

	// Token for fetching previous page in paginated result sets.
	PrevPage *PageToken
}

func wrapResponse(r *resty.Response) *Response {
	if r == nil {
		return nil
	}

	return &Response{Response: r.RawResponse}
}
