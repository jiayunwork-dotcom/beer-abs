package mixture

import "context"

// Analysis is the full output of a mixture evaluation: the ideal Beer–Lambert
// reading, the stray-light-corrected observed reading, and the deviation
// between them.
type Analysis struct {
	// Components echoes the component list.
	Components []Component

	// PathLength is the shared optical path.
	PathLength float64

	// StrayFraction is the stray-light fraction in force.
	StrayFraction float64

	// Absorbance is the ideal total absorbance A_tot.
	Absorbance float64

	// Transmittance is the ideal transmittance T = 10^(−A_tot).
	Transmittance float64

	// ObservedTransmittance is T_obs = (T + s)/(1 + s).
	ObservedTransmittance float64

	// ObservedAbsorbance is A_obs = −log10(T_obs).
	ObservedAbsorbance float64

	// Deviation is A_obs − A_ideal; negative when stray light suppresses
	// the high-absorbance reading.
	Deviation float64

	// Ideal marks whether the analysis ran without stray light.
	Ideal bool
}

// Analyze computes the complete ideal + observed picture for the mixture.
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
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	out := Analysis{
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
	return finishAnalysis(ctx, out), nil
}

// IsAdditive reports whether the mixture's ideal absorbance equals the sum of
// its components' individual absorbances, the defining invariant of the
// non-interacting model.
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

// NegativeDeviation reports whether stray light pushed the observed
// absorbance below the ideal one.
func (a Analysis) NegativeDeviation() bool {
	return a.Deviation < 0
}

// TransmittanceInUnitInterval reports whether the observed transmittance
// stayed within (0, 1], the physically meaningful range.
func (a Analysis) TransmittanceInUnitInterval() bool {
	return a.ObservedTransmittance > 0 && a.ObservedTransmittance <= 1
}
