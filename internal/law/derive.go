package law

func DeriveConcentration(absorbance, extinction, pathLength float64) (float64, error) {
	if err := ValidateExtinction(extinction); err != nil {
		return 0, err
	}
	if err := ValidatePathLength(pathLength); err != nil {
		return 0, err
	}
	if absorbance < 0 {
		return 0, NewAbsorbanceError(absorbance)
	}
	if extinction == 0 {
		return 0, NewExtinctionError(0)
	}
	return absorbance / (extinction * pathLength), nil
}

func DerivePathLength(absorbance, extinction, concentration float64) (float64, error) {
	if err := ValidateExtinction(extinction); err != nil {
		return 0, err
	}
	if err := ValidateConcentration(concentration); err != nil {
		return 0, err
	}
	if absorbance < 0 {
		return 0, NewAbsorbanceError(absorbance)
	}
	if concentration == 0 {
		return 0, NewConcentrationError(0)
	}
	return absorbance / (extinction * concentration), nil
}

func DeriveExtinction(absorbance, concentration, pathLength float64) (float64, error) {
	if err := ValidateConcentration(concentration); err != nil {
		return 0, err
	}
	if err := ValidatePathLength(pathLength); err != nil {
		return 0, err
	}
	if absorbance < 0 {
		return 0, NewAbsorbanceError(absorbance)
	}
	if concentration == 0 {
		return 0, NewConcentrationError(0)
	}
	return absorbance / (concentration * pathLength), nil
}

func (p Params) ConcentrationFor(absorbance float64) (float64, error) {
	return DeriveConcentration(absorbance, p.Extinction, p.PathLength)
}

func (p Params) AbsorbanceAtConcentration(concentration float64) (Result, error) {
	q := p.WithConcentration(concentration)
	return EvaluateParams(q)
}

func TransmittanceTarget(absorbance float64) float64 {
	return Transmittance(absorbance)
}
