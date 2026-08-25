package law

import "math"

func Transmittance(absorbance float64) float64 {
	return math.Pow(10, -absorbance)
}

func TransmittanceExpForm(absorbance float64) float64 {
	return math.Exp(-math.Ln10 * absorbance)
}

func AbsorbanceFromTransmittance(transmittance float64) float64 {
	if transmittance <= 0 {
		return math.Inf(1)
	}
	if transmittance == 1 {
		return 0
	}
	return -math.Log10(transmittance)
}

func TransmittanceBetween(aLow, aHigh float64) float64 {
	return math.Pow(10, aLow-aHigh)
}

func TransmittanceSquared(absorbance float64) float64 {
	t := Transmittance(absorbance)
	return t * t
}

func TransmittanceOfDoubledAbsorbance(absorbance float64) float64 {
	return Transmittance(2 * absorbance)
}
