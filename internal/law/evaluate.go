package law

// Result is the outcome of a single-component evaluation: the input
// parameters and the two derived quantities A and T.
type Result struct {
	// Extinction, Concentration, PathLength echo the inputs.
	Extinction    float64
	Concentration float64
	PathLength    float64

	// Absorbance is A = ε·c·L.
	Absorbance float64

	// Transmittance is T = 10^(−A).
	Transmittance float64
}

// Evaluate runs the full single-component pipeline: validate, compute the
// absorbance, then derive the transmittance.
func Evaluate(extinction, concentration, pathLength float64) (Result, error) {
	p := New(extinction, concentration, pathLength)
	return EvaluateParams(p)
}

// EvaluateParams is Evaluate bound to a Params value.
func EvaluateParams(p Params) (Result, error) {
	a, err := p.Absorbance()
	if err != nil {
		return Result{}, err
	}
	return Result{
		Extinction:    p.Extinction,
		Concentration: p.Concentration,
		PathLength:    p.PathLength,
		Absorbance:    a,
		Transmittance: Transmittance(a),
	}, nil
}

// FromAbsorbance builds a Result directly from an already computed
// absorbance, skipping the product step but keeping the shared log base.
func FromAbsorbance(extinction, concentration, pathLength, absorbance float64) Result {
	return Result{
		Extinction:    extinction,
		Concentration: concentration,
		PathLength:    pathLength,
		Absorbance:    absorbance,
		Transmittance: Transmittance(absorbance),
	}
}

// FromTransmittance builds a Result from a transmittance by inverting the
// logarithmic relation; A and T always satisfy T = 10^(−A).
func FromTransmittance(extinction, concentration, pathLength, transmittance float64) Result {
	a := AbsorbanceFromTransmittance(transmittance)
	return Result{
		Extinction:    extinction,
		Concentration: concentration,
		PathLength:    pathLength,
		Absorbance:    a,
		Transmittance: transmittance,
	}
}

// (r Result) Params recovers the input parameter triple.
func (r Result) Params() Params {
	return Params{
		Extinction:    r.Extinction,
		Concentration: r.Concentration,
		PathLength:    r.PathLength,
	}
}

// (r Result) IsConsistent verifies that the closed-form relation
// T = 10^(−A) holds for the stored numbers within a small tolerance.
func (r Result) IsConsistent(tolerance float64) bool {
	expected := Transmittance(r.Absorbance)
	diff := expected - r.Transmittance
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
