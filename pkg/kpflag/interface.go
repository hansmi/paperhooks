package kpflag

import "github.com/alecthomas/kingpin/v2"

// Implemented by [kingpin.CommandLine] (an instance of [kingpin.Application]).
type FlagGroup interface {
	Flag(name, help string) *kingpin.FlagClause
}

var _ FlagGroup = (*kingpin.Application)(nil)
var _ FlagGroup = kingpin.CommandLine
