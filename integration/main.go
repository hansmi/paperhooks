package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/paperhooks/pkg/client"
	"github.com/hansmi/paperhooks/pkg/kpflag"
)

func main() {
	var clientFlags client.Flags

	rand.Seed(time.Now().UnixNano())

	destructive := kingpin.Flag("destructive",
		"Execute potentially destructive tests. Do not use with production instances.").
		Bool()

	kpflag.RegisterClient(kingpin.CommandLine, &clientFlags)

	kingpin.CommandLine.Help = "Integration tests for the paperhooks library."
	kingpin.Parse()

	ctx := context.Background()

	client, err := clientFlags.Build()
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx); err != nil {
		log.Fatalf("Connection test failed: %v", err)
	}

	ro := readOnlyTests{
		logger: log.Default(),
		client: client,
	}

	tests := []func(context.Context) error{
		ro.tags,
		ro.correspondents,
		ro.documentTypes,
		ro.storagePaths,
		ro.customFields,
		ro.documents,
		ro.tasks,
		ro.logs,
		ro.currentUser,
		ro.users,
		ro.groups,
	}

	if *destructive {
		dt := &destructiveTests{
			logger: ro.logger,
			client: ro.client,
			mark:   fmt.Sprintf("test%x", rand.Int63()),
		}

		// Destructive tests (create, update, delete, etc.)
		tests = append(tests,
			dt.uploadDocument,
			dt.tags,
		)
	}

	for _, fn := range tests {
		if err := fn(ctx); err != nil {
			log.Fatal(err)
		}
	}

	log.Print("All tests completed successfully.")
}
