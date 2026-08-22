package law

// DeriveConcentration solves the Beer–Lambert product for the concentration
// that would produce a given absorbance:
//
//	c = A / (ε·L)
//
// The extinction and the path length must be valid, and the absorbance must
// not be negative.
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

// DerivePathLength solves L = A / (ε·c), with c > 0 required.
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

// DeriveExtinction solves ε = A / (c·L), with c > 0 required. A zero path
// length is rejected before the division.
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

// ConcentrationFor is the Params form of DeriveConcentration: it returns the
// concentration needed to reach an absorbance target with the current
// extinction and path length.
func (p Params) ConcentrationFor(absorbance float64) (float64, error) {
	return DeriveConcentration(absorbance, p.Extinction, p.PathLength)
}

// AbsorbanceAtConcentration evaluates A = ε·c·L at a given concentration and
// returns the transmittance alongside it.
func (p Params) AbsorbanceAtConcentration(concentration float64) (Result, error) {
	q := p.WithConcentration(concentration)
	return EvaluateParams(q)
}

// TransmittanceTarget reports the transmittance whose decimal logarithm is
// exactly −A; it is provided so derived quantities share the same base.
func TransmittanceTarget(absorbance float64) float64 {
	return Transmittance(absorbance)
}
