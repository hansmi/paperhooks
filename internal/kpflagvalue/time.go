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

type timeValue time.Time

var _ kingpin.Value = (*timeValue)(nil)

func (t *timeValue) String() string {
	return (*time.Time)(t).String()
}

func (t *timeValue) Set(value string) error {
	var firstErr error

	for _, layout := range timeLayouts {
		ts, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			*(*time.Time)(t) = ts
			return nil
		}

		if firstErr == nil {
			firstErr = err
		}
	}

	return fmt.Errorf("parsing %q as a time value failed (supported layouts: %q): %w", value, timeLayouts, firstErr)
}

func TimeVar(t kingpin.Settings, target *time.Time) {
	t.SetValue((*timeValue)(target))
}
