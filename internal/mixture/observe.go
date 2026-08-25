package mixture

import (
	"math"

	"beer-abs/internal/law"
)

func (m Mixture) IdealTransmittance() (float64, error) {
	a, err := m.TotalAbsorbance()
	if err != nil {
		return 0, err
	}
	return law.Transmittance(a), nil
}

func (m Mixture) IdealAbsorbance() (float64, error) {
	return m.TotalAbsorbance()
}

func (m Mixture) ObservedTransmittance() (float64, error) {
	t, err := m.IdealTransmittance()
	if err != nil {
		return 0, err
	}
	return ApplyStray(t, m.StrayFraction)
}

func (m Mixture) ObservedAbsorbance() (float64, error) {
	if m.IsIdeal() {
		return m.TotalAbsorbance()
	}
	tObs, err := m.ObservedTransmittance()
	if err != nil {
		return 0, err
	}
	if tObs <= 0 {
		return math.Inf(1), nil
	}
	return -math.Log10(tObs), nil
}

func (m Mixture) AbsorbancePair() (ideal, observed float64, err error) {
	ideal, err = m.IdealAbsorbance()
	if err != nil {
		return 0, 0, err
	}
	observed, err = m.ObservedAbsorbance()
	if err != nil {
		return 0, 0, err
	}
	return ideal, observed, nil
}

func (m Mixture) TransmittancePair() (ideal, observed float64, err error) {
	ideal, err = m.IdealTransmittance()
	if err != nil {
		return 0, 0, err
	}
	observed, err = m.ObservedTransmittance()
	if err != nil {
		return 0, 0, err
	}
	return ideal, observed, nil
}
