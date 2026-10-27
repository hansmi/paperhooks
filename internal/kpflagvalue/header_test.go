package kpflagvalue

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHTTPHeader(t *testing.T) {
	for _, tc := range []struct {
		name    string
		values  []string
		want    http.Header
		wantErr error
	}{
		{
			name: "nothing",
		},
		{
			name:    "empty",
			values:  []string{""},
			wantErr: cmpopts.AnyError,
		},
		{
			name:   "one",
			values: []string{"a:b"},
			want: http.Header{
				"A": []string{"b"},
			},
		},
		{
			name: "multiple",
			values: []string{
				"A:b",
				"host:localhost",
				"a:aaa",
			},
			want: http.Header{
				"A":    []string{"b", "aaa"},
				"Host": []string{"localhost"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var fv httpHeader
			var err error

			for _, i := range tc.values {
				if err = fv.Set(i); err != nil {
					break
				}
			}

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Set() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, (http.Header)(fv), cmpopts.EquateEmpty()); diff != "" {
					t.Errorf("Parsed values diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
