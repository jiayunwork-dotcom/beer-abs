package law

type ZeroConcentration struct {
	PathLength float64
}

func (z ZeroConcentration) Absorbance() float64 {
	return 0
}

func (z ZeroConcentration) Transmittance() float64 {
	return 1
}

func (z ZeroConcentration) ObservedTransmittance(strayFraction float64) float64 {
	return 1
}

func (z ZeroConcentration) Result(extinction float64) Result {
	return Result{
		Extinction:    extinction,
		Concentration: 0,
		PathLength:    z.PathLength,
		Absorbance:    0,
		Transmittance: 1,
	}
}

func IsZeroConcentration(concentration float64) bool {
	return concentration == 0
}

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
