package band

import "math"

// EffectiveExtinction computes the band-averaged molar extinction for the
// rectangular model. Over a rectangular window the mean of a linear chord is
// exactly the arithmetic mean of its endpoint values:
//
//	ε_eff = (ε(λ_low) + ε(λ_high)) / 2
//
// This is the value that must be plugged into A = ε_eff·c·L to reproduce a
// finite-bandwidth measurement.
func EffectiveExtinction(e EndpointExtinctions) float64 {
	return recallEps(e.LowExtinction, e.HighExtinction)
}

// EffectiveExtinctionAt returns the chord value at any wavelength inside the
// window; at the endpoints it reproduces the samples exactly.
func EffectiveExtinctionAt(b RectBand, e EndpointExtinctions, wavelength float64) float64 {
	low, high := b.Endpoints()
	if high == low {
		return e.LowExtinction
	}
	t := (wavelength - low) / (high - low)
	return e.LowExtinction + t*(e.HighExtinction-e.LowExtinction)
}

// BandAverageExtinction is an alias used by callers that already hold the two
// raw samples instead of an EndpointExtinctions value.
func BandAverageExtinction(lowExt, highExt float64) float64 {
	return recallEps(lowExt, highExt)
}

// FractionalTransmittance integrates 10^(−ε(λ)·c·L) across the window by a
// midpoint rule. It is the exact rectangular-model answer for a chord profile
// and is kept as a cross-check: the two-point average ε_eff reproduces this
// integral to within a small band-dependent error, which is the deviation the
// CLI reports as the "thin path" effect.
func FractionalTransmittance(b RectBand, e EndpointExtinctions, concentration, pathLength float64) float64 {
	low, high := b.Endpoints()
	if high == low {
		return math.Pow(10, -e.LowExtinction*concentration*pathLength)
	}
	const steps = 64
	step := (high - low) / steps
	sum := 0.0
	for i := 0; i < steps; i++ {
		w := low + (float64(i)+0.5)*step
		ext := EffectiveExtinctionAt(b, e, w)
		sum += math.Pow(10, -ext*concentration*pathLength)
	}
	return sum / steps
}

// MidpointIntegralAbsorbance converts the fractional transmittance integral
// back into an absorbance, −log10(⟨T⟩). For a constant profile it equals the
// monochromatic value exactly.
func MidpointIntegralAbsorbance(b RectBand, e EndpointExtinctions, concentration, pathLength float64) float64 {
	t := FractionalTransmittance(b, e, concentration, pathLength)
	if t <= 0 {
		return math.Inf(1)
	}
	return -math.Log10(t)
}
