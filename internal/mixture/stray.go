package mixture

// ValidateStrayFraction enforces the stray-light domain 0 ≤ s < 1. The upper
// bound is strict: at s = 1 every detector signal would be stray light and
// the transmittance formula would divide by zero.
func ValidateStrayFraction(s float64) error {
	if s < 0 || s >= 1 {
		return NewStrayError(s)
	}
	return nil
}

// ApplyStray converts an ideal transmittance T into the observed one:
//
//	T_obs = (T + s) / (1 + s)
//
// For s = 0 the observed value equals the ideal one. For s > 0 the formula
// maps the whole range onto [s/(1+s), 1], so a fully absorbing sample
// (T = 0) still shows T_obs = s/(1+s) > 0.
func ApplyStray(transmittance, strayFraction float64) (float64, error) {
	if err := ValidateStrayFraction(strayFraction); err != nil {
		return 0, err
	}
	if strayFraction == 0 {
		return transmittance, nil
	}
	return (transmittance + strayFraction) / (1 + strayFraction), nil
}

// StrayFloor returns the transmittance floor s/(1+s) seen by a fully
// absorbing sample under stray fraction s.
func StrayFloor(strayFraction float64) (float64, error) {
	if err := ValidateStrayFraction(strayFraction); err != nil {
		return 0, err
	}
	return strayFraction / (1 + strayFraction), nil
}

// IdealMapping reports the stray fraction as a fraction of the signal: s
// measures the share of the detector reading that bypassed the sample.
func IdealMapping(s float64) float64 {
	return s
}

// InverseTransmittance converts an observed transmittance back into the ideal
// one for a given stray fraction: T = T_obs·(1+s) − s.
func InverseTransmittance(observed, strayFraction float64) (float64, error) {
	if err := ValidateStrayFraction(strayFraction); err != nil {
		return 0, err
	}
	return observed*(1+strayFraction) - strayFraction, nil
}
