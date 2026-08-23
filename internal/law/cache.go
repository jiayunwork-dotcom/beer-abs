package law

import "math"

// tSession remembers the last transmittance so a repeated absorbance can skip
// the power. A new A must miss this slot; otherwise T_obs of the current
// solution is the previous sample's T.
type tSession struct {
	T    float64
	warm bool
}

var tCache = tSession{}

func recallT(absorbance float64) float64 {
	t := math.Pow(10, -absorbance)
	tCache.T = t
	tCache.warm = true
	return t
}
