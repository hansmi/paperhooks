package kpflagvalue

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTime(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   string
		wantErr error
		want    time.Time
	}{
		{
			name:    "empty",
			wantErr: cmpopts.AnyError,
		},
		{
			name:  "success",
			value: "2020-12-31T13:07:14-05:00",
			want:  time.Date(2020, time.December, 31, 13+5, 07, 14, 0, time.UTC),
		},
		{
			name:  "day only",
			value: "2004-03-02",
			want:  time.Date(2004, time.March, 02, 0, 0, 0, 0, time.Local),
		},
		{
			name:  "paperless format",
			value: "2023-02-27 23:03:50.127675+00:00",
			want:  time.Date(2023, time.February, 27, 23, 03, 50, 127675, time.UTC),
		},
		{
			name:  "paperless format short",
			value: "2023-02-24 23:00:00+00:00",
			want:  time.Date(2023, time.February, 24, 23, 00, 00, 0, time.UTC),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var fv timeValue

			err := fv.Set(tc.value)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Set() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if diff := cmp.Diff(tc.want, (time.Time)(fv), cmpopts.EquateApproxTime(time.Second)); diff != "" {
					t.Errorf("Parsed time diff (-want +got):\n%s", diff)
				}
			}
		})
	}
}
