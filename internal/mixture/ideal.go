package mixture

import "math"

func (m Mixture) IsIdeal() bool {
	return m.StrayFraction == 0
}

type IdealResult struct {
	Absorbance float64

	Transmittance float64
}

func (m Mixture) Ideal() (IdealResult, error) {
	a, err := m.TotalAbsorbance()
	if err != nil {
		return IdealResult{}, err
	}
	t := math.Pow(10, -a)
	return IdealResult{Absorbance: a, Transmittance: t}, nil
}

type TransparentSample struct {
	Components []Component
}

func (t TransparentSample) Absorbance() float64 {
	return 0
}

func (t TransparentSample) Transmittance() float64 {
	return 1
}

func (t TransparentSample) ObservedTransmittance() float64 {
	return 1
}

func (m Mixture) IsTransparent() bool {
	total, err := m.TotalAbsorbance()
	return err == nil && total == 0
}
