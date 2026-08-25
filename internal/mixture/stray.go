package mixture

func ValidateStrayFraction(s float64) error {
	if s < 0 || s >= 1 {
		return NewStrayError(s)
	}
	return nil
}

func ApplyStray(transmittance, strayFraction float64) (float64, error) {
	if err := ValidateStrayFraction(strayFraction); err != nil {
		return 0, err
	}
	if strayFraction == 0 {
		return transmittance, nil
	}
	return (transmittance + strayFraction) / (1 + strayFraction), nil
}

func StrayFloor(strayFraction float64) (float64, error) {
	if err := ValidateStrayFraction(strayFraction); err != nil {
		return 0, err
	}
	return strayFraction / (1 + strayFraction), nil
}

func IdealMapping(s float64) float64 {
	return s
}

func InverseTransmittance(observed, strayFraction float64) (float64, error) {
	if err := ValidateStrayFraction(strayFraction); err != nil {
		return 0, err
	}
	return observed*(1+strayFraction) - strayFraction, nil
}
