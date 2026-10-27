package kpflagvalue

import (
	"fmt"
	"time"

	"github.com/alecthomas/kingpin/v2"
)

var timeLayouts = []string{
	// Formats used by Paperless
	"2006-01-02 15:04:05.000000Z07:00",
	"2006-01-02 15:04:05Z07:00",

	// Other formats
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
	"2006-01-02",
}

func parseTime(value string) (time.Time, error) {
	var firstErr error

	for _, layout := range timeLayouts {
		ts, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return ts, nil
		}

		if firstErr == nil {
			firstErr = err
		}
	}

	return time.Time{}, fmt.Errorf("parsing %q as a time value failed (supported layouts: %q): %w", value, timeLayouts, firstErr)
}

type timeValue time.Time

var _ kingpin.Value = (*timeValue)(nil)

func (t *timeValue) String() string {
	return (*time.Time)(t).String()
}

func (t *timeValue) Set(value string) error {
	ts, err := parseTime(value)
	if err != nil {
		return err
	}

	*(*time.Time)(t) = ts

	return nil
}

func TimeVar(t kingpin.Settings, target *time.Time) {
	t.SetValue((*timeValue)(target))
}

type optionalTimeValue struct {
	target **time.Time
}

var _ kingpin.Value = (*optionalTimeValue)(nil)

func (t *optionalTimeValue) String() string {
	if *t.target == nil {
		return ""
	}

	return (*time.Time)(*t.target).String()
}

func (t *optionalTimeValue) Set(value string) error {
	if value == "" {
		*t.target = nil
		return nil
	}

	ts, err := parseTime(value)
	if err != nil {
		return err
	}

	*t.target = &ts

	return nil
}

func OptionalTimeVar(t kingpin.Settings, target **time.Time) {
	t.SetValue(&optionalTimeValue{target})
}
