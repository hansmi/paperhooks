package kpflagvalue

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

type commaSeparatedStrings []string

var _ kingpin.Value = (*commaSeparatedStrings)(nil)

func (s *commaSeparatedStrings) String() string {
	return fmt.Sprint(*s)
}

func (s *commaSeparatedStrings) IsCumulative() bool {
	return true
}

func (s *commaSeparatedStrings) Set(value string) error {
	for _, part := range strings.Split(value, ",") {
		if part = strings.TrimSpace(part); part != "" {
			*s = append(*s, part)
		}
	}

	return nil
}

func CommaSeparatedStringsVar(t kingpin.Settings, target *[]string) {
	t.SetValue((*commaSeparatedStrings)(target))
}
