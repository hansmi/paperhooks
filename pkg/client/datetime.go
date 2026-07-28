package client

import (
	"encoding/json"
	"time"
)

// dateOnlyLocation marks timestamps parsed from a date-only string until they
// are normalized to a real location. It uses a zero offset and is only
// compared by pointer identity.
var dateOnlyLocation = time.FixedZone("DATEONLY", 0)

// dateOrDateTime accepts timestamps both with and without a time component.
// Paperless REST API version 9 changed the document "created" field from
// a datetime to a date-only string.
type dateOrDateTime struct {
	time.Time
}

func (d *dateOrDateTime) UnmarshalJSON(data []byte) error {
	origErr := d.Time.UnmarshalJSON(data)
	if origErr == nil {
		return nil
	}

	var s string

	if json.Unmarshal(data, &s) == nil {
		if parsed, err := time.ParseInLocation(time.DateOnly, s, dateOnlyLocation); err == nil {
			d.Time = parsed
			return nil
		}
	}

	return origErr
}

// normalizeDateOnly converts a timestamp parsed from a date-only string to
// midnight of the same calendar date in the given location. Timestamps with
// a time component are returned unchanged.
func normalizeDateOnly(t time.Time, loc *time.Location) time.Time {
	if t.Location() == dateOnlyLocation {
		year, month, day := t.Date()

		return time.Date(year, month, day, 0, 0, 0, 0, loc)
	}

	return t
}
