package hostpital

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSleepRandom(t *testing.T) {
	t.Parallel()

	const (
		secMax = 0 // equivalent to 1
		tryMax = 5
	)

	hasDelay := false

	for tryCount := 0; tryCount < tryMax; tryCount++ {
		timeBefore := time.Now()

		SleepRandom(secMax) // Sleep between 0 and 999 milliseconds

		timeAfter := time.Now()

		if timeBefore != timeAfter {
			hasDelay = true

			break
		}
	}

	require.True(t, hasDelay, "failed to sleep randomly in %d tries", tryMax)
}
