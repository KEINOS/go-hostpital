package hostpital

import (
	"math/rand"
	"time"
)

// SleepRandom sleeps for a random amount of time. If secMax is 0 or 1, it
// sleeps for a random amount of time between 0 and 999 milliseconds.
func SleepRandom(secMax int) {
	const mil = 1000

	if secMax == 0 {
		secMax = 1
	}

	//nolint:gosec // not cryptographically secure random but enough for our use case.
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// In case of secMax = 1, we get a random number between 0 and 999.
	sec := randGen.Intn(secMax * mil)

	time.Sleep(time.Millisecond * time.Duration(sec))
}
