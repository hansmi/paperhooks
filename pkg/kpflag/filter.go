package kpflag

import (
	"fmt"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/paperhooks/internal/kpflagvalue"
	"github.com/hansmi/paperhooks/pkg/client"
)

type filterFlagName string

func (n filterFlagName) combine(suffix string) string {
	return fmt.Sprintf("%s_%s", n, suffix)
}

func registerBoolFilter(b *Builder, target **bool, name, help string) {
	var value bool

	b.Flag(name, help).
		PreAction(func(context *kingpin.ParseContext) error {
			*target = &value
			return nil
		}).
		BoolVar(&value)
}

func registerCharFilterSpec(b *Builder, s *client.CharFilterSpec, name, help string) {
	baseName := filterFlagName(name)

	for _, i := range []struct {
		suffix string
		help   string
		target **string
	}{
		{suffix: "iexact", help: help, target: &s.EqualsIgnoringCase},
		{suffix: "istartswith", target: &s.StartsWithIgnoringCase},
		{suffix: "iendswith", target: &s.EndsWithIgnoringCase},
		{suffix: "icontains", target: &s.ContainsIgnoringCase},
	} {
		kpflagvalue.OptionalStringVar(
			b.Flag(baseName.combine(i.suffix), i.help).
				PlaceHolder("TEXT"),
			i.target)
	}
}

func registerDateTimeFilterSpec(b *Builder, s *client.DateTimeFilterSpec, name, help string) {
	baseName := filterFlagName(name)

	for _, i := range []struct {
		suffix string
		help   string
		target **time.Time
	}{
		{suffix: "lt", help: help, target: &s.Lt},
		{suffix: "gt", target: &s.Gt},
	} {
		kpflagvalue.OptionalTimeVar(
			b.Flag(baseName.combine(i.suffix), i.help).
				PlaceHolder("DATE"),
			i.target)
	}
}

func registerIntFilterSpec(b *Builder, s *client.IntFilterSpec, name, help string) {
	baseName := filterFlagName(name)

	registerBoolFilter(b, &s.IsNull, baseName.combine("isnull"), help)

	for _, i := range []struct {
		suffix string
		target **int64
	}{
		{suffix: "exact", target: &s.Equals},
		{suffix: "lt", target: &s.Lt},
		{suffix: "lte", target: &s.Lte},
		{suffix: "gte", target: &s.Gte},
		{suffix: "gt", target: &s.Gt},
	} {
		kpflagvalue.OptionalInt64Var(
			b.Flag(baseName.combine(i.suffix), "").
				PlaceHolder("NUM"),
			i.target)
	}
}

func registerForeignKeyFilterSpec(b *Builder, s *client.ForeignKeyFilterSpec, name, help string) {
	baseName := filterFlagName(name)

	registerBoolFilter(b, &s.IsNull, baseName.combine("isnull"), help)

	kpflagvalue.OptionalInt64Var(
		b.Flag(baseName.combine("id"), "").
			PlaceHolder("NUM"),
		&s.ID)

	registerCharFilterSpec(b, &s.Name, baseName.combine("name"), "")
}
