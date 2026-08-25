package mixture

import "beer-abs/internal/law"

type Component struct {
	Label string

	Extinction float64

	Concentration float64
}

func NewComponent(label string, extinction, concentration float64) (Component, error) {
	c := Component{Label: label, Extinction: extinction, Concentration: concentration}
	if err := c.Validate(); err != nil {
		return Component{}, err
	}
	return c, nil
}

func (c Component) Absorbance(pathLength float64) (float64, error) {
	return law.Absorbance(c.Extinction, c.Concentration, pathLength)
}

func (c Component) AbsorbanceUnchecked(pathLength float64) float64 {
	return c.Extinction * c.Concentration * pathLength
}

func (c Component) Transmittance(pathLength float64) (float64, error) {
	a, err := c.Absorbance(pathLength)
	if err != nil {
		return 0, err
	}
	return law.Transmittance(a), nil
}

func (c Component) Validate() error {
	if err := law.ValidateExtinction(c.Extinction); err != nil {
		return err
	}
	if err := law.ValidateConcentration(c.Concentration); err != nil {
		return err
	}
	return nil
}

func (c Component) WithExtinction(extinction float64) Component {
	c.Extinction = extinction
	return c
}

func (c Component) WithConcentration(concentration float64) Component {
	c.Concentration = concentration
	return c
}
