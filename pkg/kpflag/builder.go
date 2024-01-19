package kpflag

import (
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

type builder struct {
	g FlagGroup
}

func (b *builder) flag(name, help string) *kingpin.FlagClause {
	return b.g.Flag(name, help).Envar(strings.ToUpper(name))
}
