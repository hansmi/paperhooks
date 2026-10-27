package kpflagvalue

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOptionalInt64(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   string
		wantErr error
		want    *int64
	}{
		{
			name: "empty",
		},
		{
			name:    "bad value",
			value:   "nope",
			wantErr: cmpopts.AnyError,
		},
		{
			name:  "success",
			value: "999",
			want:  new(int64(999)),
		},
		{
			name:  "negative",
			value: "-1",
			want:  new(int64(-1)),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := new(int64(math.MinInt64))

			err := (&optionalInt64Value{&got}).Set(tc.value)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Set() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("Parsed int64 diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
