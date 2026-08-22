package band

import (
	"fmt"
	"math"
)

// PredictAtConcentration returns a Measure whose concentration has been
// replaced, leaving the band geometry and the endpoint samples untouched.
func (m Measure) PredictAtConcentration(concentration float64) Measure {
	m.Concentration = concentration
	return m
}

// PredictAtPathLength returns a Measure whose path length has been replaced.
func (m Measure) PredictAtPathLength(pathLength float64) Measure {
	m.PathLength = pathLength
	return m
}

// RelativeDeviation is the band shift divided by the monochromatic reading,
// zero when the reference reading is zero.
func (r Result) RelativeDeviation() float64 {
	if r.MonoAbsorbance == 0 {
		return 0
	}
	return r.Deviation / r.MonoAbsorbance
}

// BandTransmittance returns 10^(−A_band) for the finite-bandwidth reading.
func (r Result) BandTransmittance() float64 {
	return math.Pow(10, -r.BandAbsorbance)
}

// MonoTransmittance returns 10^(−A_mono) for the reference reading.
func (r Result) MonoTransmittance() float64 {
	return math.Pow(10, -r.MonoAbsorbance)
}

// PercentDeviation renders the relative deviation as a percentage, useful for
// CLI output.
func (r Result) PercentDeviation() float64 {
	return 100 * r.RelativeDeviation()
}

// bracketInterval formats a wavelength interval as "[lo, hi] nm".
func bracketInterval(lo, hi float64) string {
	return "[" + fmt.Sprintf("%g", lo) + ", " + fmt.Sprintf("%g", hi) + "] nm"
}
