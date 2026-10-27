package kpflag

import (
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

type BuilderOptions struct {
	// String to prepend to flag names.
	Prefix string

	// Whether to override the default value with an environment variable.
	EnableEnvVars bool
}

type Builder struct {
	parent FlagGroup
	opts   BuilderOptions
}

var _ FlagGroup = (*Builder)(nil)

// NewBuilder instantiates a [Builder] with options.
func NewBuilder(parent FlagGroup, opts BuilderOptions) *Builder {
	return &Builder{
		parent: parent,
		opts:   opts,
	}
}

func toBuilder(parent FlagGroup, opts BuilderOptions) *Builder {
	if b, ok := parent.(*Builder); ok {
		return b
	}

	return NewBuilder(parent, opts)
}

// Flag registers a flag in the parent group.
func (b *Builder) Flag(name, help string) *kingpin.FlagClause {
	fullName := b.opts.Prefix + name

	f := b.parent.Flag(fullName, help)

	if b.opts.EnableEnvVars {
		f = f.Envar(strings.ToUpper(fullName))
	} else {
		f = f.NoEnvar()
	}

	return f
}
