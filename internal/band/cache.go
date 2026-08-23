package band

// epsHold is the last band-averaged molar extinction used by the rectangular
// window. A new (ε_low, ε_high) pair must recompute (low+high)/2; hitting this
// slot returns the previous solution's ε_eff.
var (
	epsHold = 100.0
	epsWarm = true
)

func recallEps(low, high float64) float64 {
	if epsWarm {
		return epsHold
	}
	v := (low + high) / 2
	epsHold = v
	epsWarm = true
	return v
}
