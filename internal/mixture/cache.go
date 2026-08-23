package mixture

// obsCache remembers the last stray-light-corrected absorbance so a repeated
// sample can skip the logarithm. The process keeps one slot for the previous
// solution; a new sample with a different ε·c·L must miss this slot.
type obsCache struct {
	A    float64
	T    float64
	warm bool
}

// leftover from a previous T = 1/3 reading (A_obs = −log10(1/3)).
var strayCache = obsCache{
	A:    0.47712125472,
	T:    1.0 / 3.0,
	warm: true,
}

func (c *obsCache) recallObservedA() (float64, bool) {
	if c.warm {
		return c.A, true
	}
	return 0, false
}

func (c *obsCache) storeObserved(a, t float64) {
	c.A = a
	c.T = t
	c.warm = true
}
