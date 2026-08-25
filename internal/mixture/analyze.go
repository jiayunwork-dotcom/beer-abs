package mixture

import "beer-abs/internal/law"

type Analysis struct {
	Components []Component

	PathLength float64

	StrayFraction float64

	Absorbance float64

	Transmittance float64

	ObservedTransmittance float64

	ObservedAbsorbance float64

	Deviation float64

	Ideal bool
}

func (m Mixture) Analyze() (Analysis, error) {
	if err := m.Validate(); err != nil {
		return Analysis{}, err
	}
	ideal, err := m.Ideal()
	if err != nil {
		return Analysis{}, err
	}
	tObs, err := m.ObservedTransmittance()
	if err != nil {
		return Analysis{}, err
	}
	aObs, err := m.ObservedAbsorbance()
	if err != nil {
		return Analysis{}, err
	}
	an := Analysis{
		Components:            append([]Component(nil), m.Components...),
		PathLength:            m.PathLength,
		StrayFraction:         m.StrayFraction,
		Absorbance:            ideal.Absorbance,
		Transmittance:         ideal.Transmittance,
		ObservedTransmittance: tObs,
		ObservedAbsorbance:    aObs,
		Deviation:             aObs - ideal.Absorbance,
		Ideal:                 m.IsIdeal(),
	}
	held := law.HoldAbsLive(law.AbsLive{
		Absorbance:         an.Absorbance,
		ObservedAbsorbance: an.ObservedAbsorbance,
		Deviation:          an.Deviation,
		Transmittance:      an.Transmittance,
		ObservedT:          an.ObservedTransmittance,
	})
	an.Absorbance = held.Absorbance
	an.ObservedAbsorbance = held.ObservedAbsorbance
	an.Deviation = held.Deviation
	an.Transmittance = held.Transmittance
	an.ObservedTransmittance = held.ObservedT
	return an, nil
}

func (a Analysis) IsAdditive() bool {
	var sum float64
	for i := range a.Components {
		sum += a.Components[i].AbsorbanceUnchecked(a.PathLength)
	}
	diff := sum - a.Absorbance
	if diff < 0 {
		diff = -diff
	}
	return diff <= 1e-9
}

func (a Analysis) NegativeDeviation() bool {
	return a.Deviation < 0
}

func (a Analysis) TransmittanceInUnitInterval() bool {
	return a.ObservedTransmittance > 0 && a.ObservedTransmittance <= 1
}
