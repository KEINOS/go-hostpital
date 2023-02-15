package hostpital

import (
	"math/rand"
	"time"
)

// PickRandom returns a random item from the given slice.
func PickRandom(items []string) string {
	//nolint:gosec // not cryptographically secure random but enough for our use case.
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	return items[randGen.Intn(len(items))]
}
