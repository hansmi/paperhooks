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
	var listOpts client.ListDocumentsOptions

	kpflag.RegisterClient(kingpin.CommandLine, &flags)
	kpflag.RegisterListDocumentsOptions(kingpin.CommandLine, &listOpts)

	kingpin.Parse()

	ctx := context.Background()

	cl, err := flags.Build()
	if err != nil {
		log.Fatal(err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', 0)

	defer tw.Flush()

	fmt.Fprintf(tw, "ID\tTitle\t\n")

	if err := cl.ListAllDocuments(ctx, listOpts, func(_ context.Context, doc client.Document) error {
		fmt.Fprintf(tw, "%d\t%s\t\n", doc.ID, doc.Title)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
