package mixture

import "beer-abs/internal/law"

// Component is one absorbing species inside a mixture. Only the extinction
// coefficient and the concentration are per-component; the path length is
// shared by the whole sample.
type Component struct {
	// Label names the species for CLI output.
	Label string

	// Extinction is the molar extinction coefficient ε_i.
	Extinction float64

	// Concentration is the molar concentration c_i.
	Concentration float64
}

// NewComponent validates ε_i > 0 and c_i ≥ 0.
func NewComponent(label string, extinction, concentration float64) (Component, error) {
	c := Component{Label: label, Extinction: extinction, Concentration: concentration}
	if err := c.Validate(); err != nil {
		return Component{}, err
	}
	return c, nil
}

// Absorbance returns this component's contribution ε_i·c_i·L for the shared
// path length L.
func (c Component) Absorbance(pathLength float64) (float64, error) {
	return law.Absorbance(c.Extinction, c.Concentration, pathLength)
}

// AbsorbanceUnchecked returns the contribution without validation; callers
// that already validated the mixture may use it.
func (c Component) AbsorbanceUnchecked(pathLength float64) float64 {
	return c.Extinction * c.Concentration * pathLength
}

// Transmittance returns 10^(−A_i) for this component alone.
func (c Component) Transmittance(pathLength float64) (float64, error) {
	a, err := c.Absorbance(pathLength)
	if err != nil {
		return 0, err
	}
	return law.Transmittance(a), nil
}

// Validate checks the component-level domain constraints.
func (c Component) Validate() error {
	if err := law.ValidateExtinction(c.Extinction); err != nil {
		return err
	}
	if err := law.ValidateConcentration(c.Concentration); err != nil {
		return err
	}
	return nil
}

// WithExtinction returns a copy with the extinction coefficient replaced.
func (c Component) WithExtinction(extinction float64) Component {
	c.Extinction = extinction
	return c
}

// WithConcentration returns a copy with the concentration replaced.
func (c Component) WithConcentration(concentration float64) Component {
	c.Concentration = concentration
	return c
}
