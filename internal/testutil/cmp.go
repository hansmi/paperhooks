package testutil

import (
	"time"

	"github.com/google/go-cmp/cmp"
)

func EquateTimeLocation() cmp.Option {
	return cmp.Comparer(func(a, b time.Location) bool {
		return a.String() == b.String()
	})
}
