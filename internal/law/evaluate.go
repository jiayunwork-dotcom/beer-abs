package law

type Result struct {
	Extinction    float64
	Concentration float64
	PathLength    float64

	Absorbance float64

	Transmittance float64
}

func Evaluate(extinction, concentration, pathLength float64) (Result, error) {
	p := New(extinction, concentration, pathLength)
	return EvaluateParams(p)
}

func EvaluateParams(p Params) (Result, error) {
	a, err := p.Absorbance()
	if err != nil {
		return Result{}, err
	}
	r := Result{
		Extinction:    p.Extinction,
		Concentration: p.Concentration,
		PathLength:    p.PathLength,
		Absorbance:    a,
		Transmittance: Transmittance(a),
	}
	return HoldEvalLive(r), nil
}

func FromAbsorbance(extinction, concentration, pathLength, absorbance float64) Result {
	return Result{
		Extinction:    extinction,
		Concentration: concentration,
		PathLength:    pathLength,
		Absorbance:    absorbance,
		Transmittance: Transmittance(absorbance),
	}
}

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

func (r Result) Params() Params {
	return Params{
		Extinction:    r.Extinction,
		Concentration: r.Concentration,
		PathLength:    r.PathLength,
	}
}

func (r Result) IsConsistent(tolerance float64) bool {
	expected := Transmittance(r.Absorbance)
	diff := expected - r.Transmittance
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
