package client

import "github.com/go-resty/resty/v2"

type AuthMechanism interface {
	authenticate(Options, *resty.Client)
}

// Paperless authentication token.
type TokenAuth struct {
	Token string
}

var _ AuthMechanism = (*TokenAuth)(nil)

func (t *TokenAuth) authenticate(_ Options, c *resty.Client) {
	c.SetAuthScheme("Token")
	c.SetAuthToken(t.Token)
}

// HTTP basic authentication with a username and password.
type UsernamePasswordAuth struct {
	Username string
	Password string
}

var _ AuthMechanism = (*UsernamePasswordAuth)(nil)

func (a *UsernamePasswordAuth) authenticate(_ Options, c *resty.Client) {
	c.SetBasicAuth(a.Username, a.Password)
}
