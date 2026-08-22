package law

// Params bundles the three physical inputs of a single-component
// Beer–Lambert measurement.
type Params struct {
	// Extinction is the molar extinction coefficient ε (L·mol⁻¹·cm⁻¹).
	Extinction float64

	// Concentration is the molar concentration c (mol·L⁻¹).
	Concentration float64

	// PathLength is the optical path length L (cm).
	PathLength float64
}

// New creates a Params value without validating anything; validation happens
// when a computed quantity is requested. Prefer NewChecked when the inputs
// must be rejected as early as possible.
func New(extinction, concentration, pathLength float64) Params {
	return Params{
		Extinction:    extinction,
		Concentration: concentration,
		PathLength:    pathLength,
	}
}

// NewChecked validates the three inputs and returns a Params value only when
// all of them lie in the physical domain (ε > 0, c ≥ 0, L > 0).
func NewChecked(extinction, concentration, pathLength float64) (Params, error) {
	p := New(extinction, concentration, pathLength)
	if err := p.Validate(); err != nil {
		return Params{}, err
	}
	return p, nil
}

// WithExtinction returns a copy of p with the extinction coefficient replaced.
func (p Params) WithExtinction(extinction float64) Params {
	p.Extinction = extinction
	return p
}

// WithConcentration returns a copy of p with the concentration replaced.
func (p Params) WithConcentration(concentration float64) Params {
	p.Concentration = concentration
	return p
}

// WithPathLength returns a copy of p with the path length replaced.
func (p Params) WithPathLength(pathLength float64) Params {
	p.PathLength = pathLength
	return p
}

// Clone returns a deep copy of p.
func (p Params) Clone() Params {
	return Params{
		Extinction:    p.Extinction,
		Concentration: p.Concentration,
		PathLength:    p.PathLength,
	}
}

// Fields returns the three scalars in canonical order ε, c, L.
func (p Params) Fields() (extinction, concentration, pathLength float64) {
	return p.Extinction, p.Concentration, p.PathLength
}
