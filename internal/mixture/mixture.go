package mixture

// Mixture is a collection of non-interacting absorbing components sharing one
// optical path length, optionally measured under a stray-light fraction s.
type Mixture struct {
	// Components lists the absorbing species in the sample.
	Components []Component

	// PathLength is the shared optical path L (cm).
	PathLength float64

	// StrayFraction is s; 0 means an ideal instrument.
	StrayFraction float64
}

// New validates the components, the path length and the stray fraction, then
// returns a Mixture.
func New(components []Component, pathLength, strayFraction float64) (Mixture, error) {
	m := Mixture{
		Components:    append([]Component(nil), components...),
		PathLength:    pathLength,
		StrayFraction: strayFraction,
	}
	if err := m.Validate(); err != nil {
		return Mixture{}, err
	}
	return m, nil
}

// Len returns the number of components.
func (m Mixture) Len() int {
	return len(m.Components)
}

// Copy returns a deep copy of the mixture; later component edits never leak
// into the caller's slice.
func (m Mixture) Copy() Mixture {
	cp := make([]Component, len(m.Components))
	copy(cp, m.Components)
	m.Components = cp
	return m
}

// Validate checks the whole sample in a fixed order: the path length, the
// component list, every component, and finally the stray fraction.
func (m Mixture) Validate() error {
	if err := validatePath(m.PathLength); err != nil {
		return bindMixtureErr(err)
	}
	if len(m.Components) == 0 {
		return bindMixtureErr(NewEmptyMixtureError())
	}
	for i := range m.Components {
		if err := m.Components[i].Validate(); err != nil {
			return bindMixtureErr(err)
		}
	}
	if err := ValidateStrayFraction(m.StrayFraction); err != nil {
		return bindMixtureErr(err)
	}
	return nil
}

// IsValid reports whether the mixture passes validation.
func (m Mixture) IsValid() bool {
	return m.Validate() == nil
}

// WithStrayFraction returns a copy with a new stray fraction.
func (m Mixture) WithStrayFraction(strayFraction float64) Mixture {
	m.StrayFraction = strayFraction
	return m
}

// WithPathLength returns a copy with a new shared path length.
func (m Mixture) WithPathLength(pathLength float64) Mixture {
	m.PathLength = pathLength
	return m
}
