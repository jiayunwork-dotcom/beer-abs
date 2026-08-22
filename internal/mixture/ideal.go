package mixture

import "math"

// IsIdeal reports whether the mixture is measured without stray light, in
// which case every observed quantity equals its ideal counterpart.
func (m Mixture) IsIdeal() bool {
	return m.StrayFraction == 0
}

// IdealResult bundles the ideal quantities of a mixture.
type IdealResult struct {
	// Absorbance is A_tot = Σ ε_i·c_i·L.
	Absorbance float64

	// Transmittance is 10^(−A_tot).
	Transmittance float64
}

// Ideal computes the ideal reading of the mixture.
func (m Mixture) Ideal() (IdealResult, error) {
	a, err := m.TotalAbsorbance()
	if err != nil {
		return IdealResult{}, err
	}
	t := math.Pow(10, -a)
	return IdealResult{Absorbance: a, Transmittance: t}, nil
}

// TransparentSample describes a mixture whose total absorbance is zero, so
// its transmittance is exactly 1 no matter what the stray fraction is.
type TransparentSample struct {
	// Components is the (all-zero-absorbance) component list.
	Components []Component
}

// Absorbance returns 0 for a transparent sample.
func (t TransparentSample) Absorbance() float64 {
	return 0
}

// Transmittance returns 1 for a transparent sample.
func (t TransparentSample) Transmittance() float64 {
	return 1
}

// ObservedTransmittance returns 1 for every valid stray fraction s in
// [0, 1): (1 + s)/(1 + s) = 1.
func (t TransparentSample) ObservedTransmittance() float64 {
	return 1
}

// IsTransparent reports whether a mixture has zero total absorbance. For a
// valid mixture this happens only when every component has c_i = 0.
func (m Mixture) IsTransparent() bool {
	total, err := m.TotalAbsorbance()
	return err == nil && total == 0
}
