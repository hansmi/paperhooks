package client

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hansmi/paperhooks/internal/testutil"
)

func TestDateOrDateTimeUnmarshalJSON(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "datetime with offset",
			input: `"2023-06-30T22:00:00+02:00"`,
			want:  time.Date(2023, time.June, 30, 22, 0, 0, 0, time.FixedZone("", 2*60*60)),
		},
		{
			name:  "datetime in UTC",
			input: `"2023-06-30T22:00:00Z"`,
			want:  time.Date(2023, time.June, 30, 22, 0, 0, 0, time.UTC),
		},
		{
			name:  "date only",
			input: `"2026-07-23"`,
			want:  time.Date(2026, time.July, 23, 0, 0, 0, 0, dateOnlyLocation),
		},
		{
			name:  "null",
			input: `null`,
		},
		{
			name:    "invalid",
			input:   `"not a timestamp"`,
			wantErr: true,
		},
		{
			name:    "number",
			input:   `12345`,
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got dateOrDateTime

			err := got.UnmarshalJSON([]byte(tc.input))

			if tc.wantErr != (err != nil) {
				t.Errorf("UnmarshalJSON(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}

			if err == nil {
				if !got.Time.Equal(tc.want) {
					t.Errorf("UnmarshalJSON(%q) = %v, want %v", tc.input, got.Time, tc.want)
				}

				if wantLoc := tc.want.Location(); got.Time.Location() != wantLoc && wantLoc == dateOnlyLocation {
					t.Errorf("UnmarshalJSON(%q) location = %v, want date-only marker", tc.input, got.Time.Location())
				}
			}
		})
	}
}

func TestNormalizeDateOnly(t *testing.T) {
	plus2 := time.FixedZone("UTC+2", 2*60*60)

	for _, tc := range []struct {
		name string
		t    time.Time
		loc  *time.Location
		want time.Time
	}{
		{
			name: "date only to UTC",
			t:    time.Date(2026, time.July, 23, 0, 0, 0, 0, dateOnlyLocation),
			loc:  time.UTC,
			want: time.Date(2026, time.July, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "date only to fixed offset",
			t:    time.Date(2026, time.July, 23, 0, 0, 0, 0, dateOnlyLocation),
			loc:  plus2,
			want: time.Date(2026, time.July, 23, 0, 0, 0, 0, plus2),
		},
		{
			name: "datetime remains unchanged",
			t:    time.Date(2023, time.June, 30, 22, 0, 0, 0, time.UTC),
			loc:  plus2,
			want: time.Date(2023, time.June, 30, 22, 0, 0, 0, time.UTC),
		},
		{
			name: "zero value remains unchanged",
			loc:  plus2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeDateOnly(tc.t, tc.loc)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("normalizeDateOnly() diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want.Location(), got.Location(), testutil.EquateTimeLocation()); diff != "" {
				t.Errorf("normalizeDateOnly() location diff (-want +got):\n%s", diff)
			}
		})
	}
}
