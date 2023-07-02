package client

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
)

const defaultPerPage = 25

type PageToken struct {
	// Page number for paginated result sets.
	number int `url:"-"`

	// Number of items on paginated result sets.
	size int `url:"-"`
}

var _ query.Encoder = (*PageToken)(nil)

func (t *PageToken) values() (int, int) {
	number := 1
	size := defaultPerPage

	if t != nil {
		// Zero is an invalid page number.
		if t.number > 0 {
			number = t.number
		}

		if t.size > 0 {
			size = t.size
		}
	}

	return number, size
}

func (t *PageToken) EncodeValues(_ string, v *url.Values) error {
	number, size := t.values()

	if number > 0 {
		v.Set("page", strconv.FormatUint(uint64(number), 10))
	}

	// The page size is always set in the URL to never rely on the server's
	// default.
	v.Set("page_size", strconv.FormatUint(uint64(size), 10))

	return nil
}

func pageTokenFromURL(raw string) (*PageToken, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	query := parsed.Query()

	var t PageToken

	for _, i := range []struct {
		dest *int
		name string
	}{
		{&t.number, "page"},
		{&t.size, "page_size"},
	} {
		raw := query.Get(i.name)
		if raw == "" {
			continue
		}

		num, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing %q value %q: %w", i.name, raw, err)
		}

		*i.dest = int(num)
	}

	return &t, nil
}

type ListOptions struct {
	Page *PageToken
}
