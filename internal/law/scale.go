package law

// ScaleConcentration multiplies the concentration by a non-negative factor
// and returns the scaled parameters. A factor of 2 doubles the absorbance;
// a factor of 0 produces the transparent sample.
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

// ScalePathLength multiplies the path length by a positive factor.
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

// ResponseRatio is A(factor·x) / A(x), the linearity signature of the law:
// it equals the factor itself for any scale applied to c or L.
func ResponseRatio(baseA, scaledA float64) float64 {
	if baseA == 0 {
		return 0
	}
	return scaledA / baseA
}

// PredictAbsorbanceForFactor scales a valid Params value by a concentration
// factor and returns the predicted absorbance and transmittance.
func PredictAbsorbanceForFactor(p Params, factor float64) (Result, error) {
	q, err := ScaleConcentration(p, factor)
	if err != nil {
		return Result{}, err
	}
	return EvaluateParams(q)
}

// ExtinctionAtTargetPath solves for the extinction that achieves a target
// absorbance with the current concentration and path length.
func (p Params) ExtinctionAtTargetPath(absorbance float64) (float64, error) {
	return DeriveExtinction(absorbance, p.Concentration, p.PathLength)
}

// AbsorbanceSum adds two absorbances and reports the corresponding
// transmittance; it mirrors the additivity of a two-species mixture.
func AbsorbanceSum(a, b float64) (float64, float64) {
	sum := a + b
	return sum, Transmittance(sum)
}
