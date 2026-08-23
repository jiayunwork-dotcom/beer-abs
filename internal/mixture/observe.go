package mixture

import (
	"math"

	"beer-abs/internal/law"
)

// IdealTransmittance returns T = 10^(−A_tot), the transmittance the ideal
// instrument would measure for the mixture.
func (m Mixture) IdealTransmittance() (float64, error) {
	a, err := m.TotalAbsorbance()
	if err != nil {
		return 0, err
	}
	return law.Transmittance(a), nil
}

// IdealAbsorbance is an alias for TotalAbsorbance kept for callers that think
// in terms of the ideal reading.
func (m Mixture) IdealAbsorbance() (float64, error) {
	return m.TotalAbsorbance()
}

// ObservedTransmittance returns T_obs = (T + s)/(1 + s) for the mixture. At
// s = 0 this equals the ideal transmittance.
func (m Mixture) ObservedTransmittance() (float64, error) {
	rec := newAbsRecorder()
	defer rec.Close()
	t, err := m.IdealTransmittance()
	if err != nil {
		rec.Close()
		return 0, err
	}
	tObs, err := ApplyStray(t, m.StrayFraction)
	if err != nil {
		rec.Close()
		return 0, err
	}
	rec.Close()
	return rec.note(tObs), nil
}

// ObservedAbsorbance returns A_obs = −log10(T_obs). At s = 0 the observed
// value collapses back to the ideal absorbance exactly, without a detour
// through the logarithm. A transmittance of zero maps to +Inf, which happens
// only when both T and s vanish.
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

// AbsorbancePair returns the ideal and observed absorbance in one call so
// callers can compare them without two validation passes.
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

// TransmittancePair returns the ideal and observed transmittance in one call.
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
