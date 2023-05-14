package client

import (
	"context"
	"log"
)

func Example_filter() {
	cl := New(Options{ /* … */ })

	var opt ListStoragePathsOptions

	opt.Ordering.Field = "name"
	opt.Name.ContainsIgnoringCase = String("sales")
	opt.Path.StartsWithIgnoringCase = String("2019/")

	for {
		got, resp, err := cl.ListStoragePaths(context.Background(), &opt)
		if err != nil {
			log.Fatalf("Listing storage paths failed: %v", err)
		}

		for _, i := range got {
			log.Printf("%s (%d documents)", i.Name, i.DocumentCount)
		}

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}
}
