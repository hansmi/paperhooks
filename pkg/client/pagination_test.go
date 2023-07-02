package client

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/go-querystring/query"
)

func TestPageTokenEncodeValues(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      *PageToken
		wantValues url.Values
		want       *PageToken
	}{
		{
			name: "defaults",
			wantValues: url.Values{
				"page":      []string{"1"},
				"page_size": []string{fmt.Sprint(defaultPerPage)},
			},
			want: &PageToken{
				number: 1,
				size:   defaultPerPage,
			},
		},
		{
			name:  "zero",
			input: &PageToken{},
			wantValues: url.Values{
				"page":      []string{"1"},
				"page_size": []string{fmt.Sprint(defaultPerPage)},
			},
			want: &PageToken{
				number: 1,
				size:   defaultPerPage,
			},
		},
		{
			name: "custom values",
			input: &PageToken{
				number: 1,
				size:   2,
			},
			wantValues: url.Values{
				"page":      []string{"1"},
				"page_size": []string{"2"},
			},
			want: &PageToken{
				number: 1,
				size:   2,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values, err := query.Values(struct {
				T *PageToken
			}{tc.input})

			if err != nil {
				t.Errorf("Encoding values failed: %v", err)
			}

			if diff := cmp.Diff(tc.wantValues, values, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Values diff (-want +got):\n%s", diff)
			}

			u := url.URL{
				RawQuery: values.Encode(),
			}

			parsed, err := pageTokenFromURL(u.String())

			if err != nil {
				t.Errorf("Parsing token from URL %q failed: %v", u.String(), err)
			} else if diff := cmp.Diff(tc.want, parsed, cmp.AllowUnexported(PageToken{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Token diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPageTokenFromURL(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    *PageToken
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name:  "number and size",
			input: "/foo?bar=baz&page_size=111&page=222",
			want: &PageToken{
				number: 222,
				size:   111,
			},
		},
		{
			name:  "number only",
			input: "/foo?page=10062",
			want: &PageToken{
				number: 10062,
			},
		},
		{
			name:  "size only",
			input: "/foo?page_size=13434",
			want: &PageToken{
				size: 13434,
			},
		},
		{
			name:    "bad number syntax",
			input:   "/foo?page=0x123",
			wantErr: strconv.ErrSyntax,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pageTokenFromURL(tc.input)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("pageTokenFromURL() error diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(PageToken{}), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Token diff (-want +got):\n%s", diff)
			}
		})
	}
}
