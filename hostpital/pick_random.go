package hostpital

import (
	"math/rand"
	"time"
)

// PickRandom returns a random item from the given slice.
func PickRandom(items []string) string {
	rand.Seed(time.Now().UnixNano())

	//nolint:gosec // not cryptographically secure random but enough for our use case.
	return items[rand.Intn(len(items))]
}
