package law

// Validate checks the physical domain of the three Beer–Lambert inputs.
// The validation order is fixed: extinction first, then concentration,
// then path length. Only the first failing constraint is reported.
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

// Validate checks that p lies inside the physical domain.
func (p Params) Validate() error {
	return Validate(p.Extinction, p.Concentration, p.PathLength)
}

// ValidateExtinction reports whether a single extinction coefficient is
// usable in an absorbance computation.
func ValidateExtinction(extinction float64) error {
	if extinction <= 0 {
		return NewExtinctionError(extinction)
	}
	return nil
}

// ValidateConcentration reports whether a single concentration is usable in
// an absorbance computation. Zero is allowed: a sample without absorbing
// species transmits everything.
func ValidateConcentration(concentration float64) error {
	if concentration < 0 {
		return NewConcentrationError(concentration)
	}
	return nil
}

// ValidatePathLength reports whether a single path length is usable in an
// absorbance computation.
func ValidatePathLength(pathLength float64) error {
	if pathLength <= 0 {
		return NewPathLengthError(pathLength)
	}
	return nil
}

// IsValid reports whether the three inputs pass validation without returning
// the error itself. Useful for quick guards in glue code.
func IsValid(extinction, concentration, pathLength float64) bool {
	return Validate(extinction, concentration, pathLength) == nil
}

// (p Params) IsValid reports whether p passes validation.
func (p Params) IsValid() bool {
	return IsValid(p.Extinction, p.Concentration, p.PathLength)
}
