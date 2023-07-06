package client

import (
	"context"
	"log"
)

func Example_pagination() {
	cl := New(Options{ /* … */ })

	var opt ListDocumentsOptions
	var all []Document

	for {
		documents, resp, err := cl.ListDocuments(context.Background(), opt)
		if err != nil {
			log.Fatalf("Listing documents failed: %v", err)
		}

		all = append(all, documents...)

		if resp.NextPage == nil {
			break
		}

		opt.Page = resp.NextPage
	}

	log.Printf("Received %d documents in total.", len(all))
}
