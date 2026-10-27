package kpflagvalue

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOptionalString(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   string
		wantErr error
		want    *string
	}{
		{
			name: "empty",
			want: new(""),
		},
		{
			name:  "success",
			value: "hello",
			want:  new("hello"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := new("---")

			err := (&optionalStringValue{&got}).Set(tc.value)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Set() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("Parsed string diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
