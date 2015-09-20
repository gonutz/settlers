package game

type randomNumberGenerator struct {
	index int
}

func newRNG(seed int) *randomNumberGenerator {
	return &randomNumberGenerator{seed % len(randomNumbers)}
}

func (r *randomNumberGenerator) next() int {
	r.index = (r.index + 1) % len(randomNumbers)
	return randomNumbers[r.index]
}
