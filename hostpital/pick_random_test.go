package hostpital

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPickRandom(t *testing.T) {
	t.Parallel()

	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	const (
		itemFirst = "a"
		itemLast  = "j"
		tryMax    = 1000
	)

	foundFirst := false
	foundLast := false

	for range tryMax {
		SleepRandom(1)

		picked := PickRandom(items)
		if picked == itemFirst {
			foundFirst = true
		}

		if picked == itemLast {
			foundLast = true
		}

		if foundFirst && foundLast {
			break
		}
	}

	require.True(t, foundFirst && foundLast, "failed to pick all the items randomly in %d tries", tryMax)
}
