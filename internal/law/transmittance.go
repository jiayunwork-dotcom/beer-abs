package law

import "math"

// Transmittance converts an absorbance into a transmittance using a decimal
// logarithm base: T = 10^(−A). This is the canonical definition of the
// Beer–Lambert law and every quantity in this module is consistent with it.
func Transmittance(absorbance float64) float64 {
	return math.Pow(10, -absorbance)
}

// TransmittanceExpForm is the algebraically identical form
// T = exp(−ln(10)·A). Keeping both entry points makes the shared base
// explicit and gives tests a second route to the same value.
func TransmittanceExpForm(absorbance float64) float64 {
	return math.Exp(-math.Ln10 * absorbance)
}

// AbsorbanceFromTransmittance inverts T = 10^(−A): given a transmittance in
// (0, 1], it returns A = −log10(T). A transmittance of exactly 1 maps to 0;
// non-positive values map to +Inf because no real absorbance produces them.
func AbsorbanceFromTransmittance(transmittance float64) float64 {
	if transmittance <= 0 {
		return math.Inf(1)
	}
	if transmittance == 1 {
		return 0
	}
	return -math.Log10(transmittance)
}

// TransmittanceBetween returns the transmittance over the span between two
// absorbance levels, T(a2)/T(a1) = 10^(a1−a2). It is used to verify the
// doubling rules without accumulating floating point drift.
func TransmittanceBetween(aLow, aHigh float64) float64 {
	return math.Pow(10, aLow-aHigh)
}

// TransmittanceSquared returns T²  for a given absorbance; the doubling rule
// states that doubling the absorbance squares the transmittance.
func TransmittanceSquared(absorbance float64) float64 {
	t := Transmittance(absorbance)
	return t * t
}

// TransmittanceOfDoubledAbsorbance returns the transmittance after the
// absorbance has been doubled.
func TransmittanceOfDoubledAbsorbance(absorbance float64) float64 {
	return Transmittance(2 * absorbance)
}
