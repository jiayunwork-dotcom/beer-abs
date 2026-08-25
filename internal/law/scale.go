package law

func ScaleConcentration(p Params, factor float64) (Params, error) {
	if factor < 0 {
		return Params{}, &ParamError{
			Field:  "scale_factor",
			Value:  factor,
			Reason: "浓度缩放因子必须不小于 0",
			Cause:  ErrNegativeConcentration,
		}
	}
	p.Concentration *= factor
	return p, nil
}

func ScalePathLength(p Params, factor float64) (Params, error) {
	if factor <= 0 {
		return Params{}, &ParamError{
			Field:  "scale_factor",
			Value:  factor,
			Reason: "光程缩放因子必须大于 0",
			Cause:  ErrPathLengthNotPositive,
		}
	}
	p.PathLength *= factor
	return p, nil
}

func ResponseRatio(baseA, scaledA float64) float64 {
	if baseA == 0 {
		return 0
	}
	return scaledA / baseA
}

func PredictAbsorbanceForFactor(p Params, factor float64) (Result, error) {
	q, err := ScaleConcentration(p, factor)
	if err != nil {
		return Result{}, err
	}
	return EvaluateParams(q)
}

func (p Params) ExtinctionAtTargetPath(absorbance float64) (float64, error) {
	return DeriveExtinction(absorbance, p.Concentration, p.PathLength)
}

func AbsorbanceSum(a, b float64) (float64, float64) {
	sum := a + b
	return sum, Transmittance(sum)
}
