package model

import (
	"math/rand"
	"testing"
)

func TestRandomGeneraterWithSameSeedAllTimes(t *testing.T) {
	var initialSeed int64
	initialSeed = 10
	randomGenerator := rand.NewSource(initialSeed)
	antoherRandomGenerator := rand.NewSource(initialSeed)

	for i := 0; i < 10; i++ {
		if randomGenerator.Int63() != antoherRandomGenerator.Int63() {
			t.Errorf("Both random generator gave a different number")
		}
	}
}
