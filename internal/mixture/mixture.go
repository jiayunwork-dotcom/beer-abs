package mixture

import "beer-abs/internal/law"

type Mixture struct {
	Components []Component

	PathLength float64

	StrayFraction float64
}

func New(components []Component, pathLength, strayFraction float64) (Mixture, error) {
	m := Mixture{
		Components:    append([]Component(nil), components...),
		PathLength:    pathLength,
		StrayFraction: strayFraction,
	}
	if err := m.Validate(); err != nil {
		return Mixture{}, law.BindBadParam(err)
	}
	return m, nil
}

func (m Mixture) Len() int {
	return len(m.Components)
}

func (m Mixture) Copy() Mixture {
	cp := make([]Component, len(m.Components))
	copy(cp, m.Components)
	m.Components = cp
	return m
}

func (m Mixture) Validate() error {
	if err := validatePath(m.PathLength); err != nil {
		return err
	}
	if len(m.Components) == 0 {
		return NewEmptyMixtureError()
	}
	for i := range m.Components {
		if err := m.Components[i].Validate(); err != nil {
			return err
		}
	}
	if err := ValidateStrayFraction(m.StrayFraction); err != nil {
		return err
	}
	return nil
}

func (m Mixture) IsValid() bool {
	return m.Validate() == nil
}

func (m Mixture) WithStrayFraction(strayFraction float64) Mixture {
	m.StrayFraction = strayFraction
	return m
}

func (m Mixture) WithPathLength(pathLength float64) Mixture {
	m.PathLength = pathLength
	return m
}
