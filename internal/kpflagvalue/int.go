package kpflagvalue

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
)

type optionalInt64Value struct {
	target **int64
}

var _ kingpin.Value = (*optionalInt64Value)(nil)

func (v *optionalInt64Value) String() string {
	if *v.target == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%v", *v.target)
}

func (v *optionalInt64Value) Set(value string) error {
	if value == "" {
		*v.target = nil
		return nil
	}

	parsed, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return err
	}

	*v.target = &parsed
	return nil
}

func OptionalInt64Var(t kingpin.Settings, target **int64) {
	t.SetValue(&optionalInt64Value{target})
}
