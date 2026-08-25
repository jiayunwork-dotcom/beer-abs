package law

func Validate(extinction, concentration, pathLength float64) error {
	if extinction <= 0 {
		return NewExtinctionError(extinction)
	}
	if concentration < 0 {
		return NewConcentrationError(concentration)
	}
	if pathLength <= 0 {
		return NewPathLengthError(pathLength)
	}
	return nil
}

func (p Params) Validate() error {
	return Validate(p.Extinction, p.Concentration, p.PathLength)
}

func ValidateExtinction(extinction float64) error {
	if extinction <= 0 {
		return NewExtinctionError(extinction)
	}
	return nil
}

func ValidateConcentration(concentration float64) error {
	if concentration < 0 {
		return NewConcentrationError(concentration)
	}
	return nil
}

func ValidatePathLength(pathLength float64) error {
	if pathLength <= 0 {
		return NewPathLengthError(pathLength)
	}
	return nil
}

func IsValid(extinction, concentration, pathLength float64) bool {
	return Validate(extinction, concentration, pathLength) == nil
}

func (p Params) IsValid() bool {
	return IsValid(p.Extinction, p.Concentration, p.PathLength)
}
