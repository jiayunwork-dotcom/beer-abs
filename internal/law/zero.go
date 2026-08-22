package law

// ZeroConcentration describes a sample whose concentration is exactly zero.
// Such a sample has no absorbing species: A = 0, T = 1, and no stray-light
// correction can change that, because the observed transmittance collapses to
// (1 + s)/(1 + s) = 1 regardless of the stray fraction.
type ZeroConcentration struct {
	// PathLength is the optical path still in place.
	PathLength float64
}

// Absorbance returns 0 for a zero-concentration sample.
func (z ZeroConcentration) Absorbance() float64 {
	return 0
}

// Transmittance returns 1 for a zero-concentration sample.
func (z ZeroConcentration) Transmittance() float64 {
	return 1
}

// ObservedTransmittance returns 1 independently of the stray fraction s.
func (z ZeroConcentration) ObservedTransmittance(strayFraction float64) float64 {
	return 1
}

// Result renders the zero-concentration sample as a full Result.
func (z ZeroConcentration) Result(extinction float64) Result {
	return Result{
		Extinction:    extinction,
		Concentration: 0,
		PathLength:    z.PathLength,
		Absorbance:    0,
		Transmittance: 1,
	}
}

// IsZeroConcentration reports whether a concentration is exactly zero.
func IsZeroConcentration(concentration float64) bool {
	return concentration == 0
}

// EvaluateZero returns a Result for c = 0 without touching the extinction
// coefficient, which plays no role when there is nothing to absorb.
func EvaluateZero(pathLength float64) (Result, error) {
	if err := ValidatePathLength(pathLength); err != nil {
		return Result{}, err
	}
	return Result{
		Extinction:    0,
		Concentration: 0,
		PathLength:    pathLength,
		Absorbance:    0,
		Transmittance: 1,
	}, nil
}
