package hook

import (
	"log"
	"os"
)

// SuccessOrDie never returns. It invokes the function and exits with
// a non-zero status code in case of an error, zero in case of success.
func SuccessOrDie(fn func() error) {
	if err := fn(); err != nil {
		log.Fatalf("Error: %v", err)
	}

	os.Exit(0)
}
