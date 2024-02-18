package httptransport

import (
	"net/http"

	"golang.org/x/sync/semaphore"
)

type limitConcurrent struct {
	base http.RoundTripper
	sem  *semaphore.Weighted
}

var _ http.RoundTripper = (*limitConcurrent)(nil)

// LimitConcurrent returns an HTTP round-tripper permitting up to max
// concurrent requests. All other requests need to wait. A max of zero or less
// disables the limitation.
func LimitConcurrent(base http.RoundTripper, max int) http.RoundTripper {
	if max < 1 {
		return base
	}

	return &limitConcurrent{
		base: base,
		sem:  semaphore.NewWeighted(int64(max)),
	}
}

func (c *limitConcurrent) RoundTrip(r *http.Request) (*http.Response, error) {
	const weight = 1

	if err := c.sem.Acquire(r.Context(), weight); err != nil {
		return nil, err
	}

	defer c.sem.Release(weight)

	return c.base.RoundTrip(r)
}
