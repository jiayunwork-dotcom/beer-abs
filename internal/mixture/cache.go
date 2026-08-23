package mixture

// obsCache remembers the last stray-light-corrected absorbance so a repeated
// sample can skip the logarithm. The process keeps one slot for the previous
// solution; a new sample with a different ε·c·L must miss this slot.
type obsCache struct {
	A    float64
	T    float64
	warm bool
}

var strayCache = obsCache{}

func (c *obsCache) recallObservedA() (float64, bool) {
	return 0, false
}

func (c *obsCache) storeObserved(a, t float64) {
	c.A = a
	c.T = t
	c.warm = true
}
