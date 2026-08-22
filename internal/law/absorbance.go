package law

// Absorbance computes A = ε·c·L for a single component. The inputs are
// validated first; a violation returns an error instead of a number.
func Absorbance(extinction, concentration, pathLength float64) (float64, error) {
	if err := Validate(extinction, concentration, pathLength); err != nil {
		return 0, err
	}
	return extinction * concentration * pathLength, nil
}

// Absorbance computes the absorbance of p.
func (p Params) Absorbance() (float64, error) {
	return Absorbance(p.Extinction, p.Concentration, p.PathLength)
}

// AbsorbanceUnchecked multiplies ε·c·L without validation. Callers that have
// already validated their inputs (or that only probe trends) may use this,
// but the checked variants should be preferred for user-facing entry points.
func AbsorbanceUnchecked(extinction, concentration, pathLength float64) float64 {
	return extinction * concentration * pathLength
}

// AbsorbanceUnchecked computes the absorbance of p without validation.
func (p Params) AbsorbanceUnchecked() float64 {
	return p.Extinction * p.Concentration * p.PathLength
}

// ProductOf separates the factors of an absorbance: given the absorbance of a
// valid single-component system and the concentration/path length already used,
// it returns the extinction coefficient that reproduces it. This is only
// meaningful for A > 0 and c > 0, L > 0; the result is otherwise undefined.
func ProductOf(absorbance, concentration, pathLength float64) float64 {
	return absorbance / (concentration * pathLength)
}
