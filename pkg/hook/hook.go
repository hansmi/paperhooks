package hook

import (
	"log"
	"os"
)

func SuccessOrDie(fn func() error) {
	if err := fn(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	os.Exit(0)
}
