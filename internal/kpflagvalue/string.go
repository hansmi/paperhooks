package kpflagvalue

import "github.com/alecthomas/kingpin/v2"

type optionalStringValue struct {
	target **string
}

var _ kingpin.Value = (*optionalStringValue)(nil)

func (s *optionalStringValue) String() string {
	if *s.target == nil {
		return "<nil>"
	}

	return **s.target
}

func (s *optionalStringValue) Set(value string) error {
	*s.target = &value
	return nil
}

func OptionalStringVar(t kingpin.Settings, target **string) {
	t.SetValue(&optionalStringValue{target})
}
