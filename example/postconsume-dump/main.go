package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/paperhooks/pkg/kpflag"
	"github.com/hansmi/paperhooks/pkg/postconsume"
	"github.com/kr/pretty"
)

func main() {
	var postconsumeFlags postconsume.Flags

	kpflag.RegisterPostConsume(kingpin.CommandLine, &postconsumeFlags)

	kingpin.Parse()

	pretty.Logf("%# v", postconsumeFlags)
}
