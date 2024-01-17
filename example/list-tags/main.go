package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/paperhooks/pkg/client"
	"github.com/hansmi/paperhooks/pkg/kpflag"
)

func main() {
	var flags client.Flags

	kpflag.RegisterClient(kingpin.CommandLine, &flags)

	kingpin.Parse()

	ctx := context.Background()

	cl, err := flags.Build()
	if err != nil {
		log.Fatal(err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)

	defer tw.Flush()

	fmt.Fprintf(tw, "ID\tName\tCount\t\n")

	if err := cl.ListAllTags(ctx, client.ListTagsOptions{}, func(_ context.Context, tag client.Tag) error {
		fmt.Fprintf(tw, "%d\t%s\t%d\t\n", tag.ID, tag.Name, tag.DocumentCount)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
