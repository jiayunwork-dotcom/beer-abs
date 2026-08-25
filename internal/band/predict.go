package band

import (
	"fmt"
	"math"
)

func (m Measure) PredictAtConcentration(concentration float64) Measure {
	m.Concentration = concentration
	return m
}

func (m Measure) PredictAtPathLength(pathLength float64) Measure {
	m.PathLength = pathLength
	return m
}

func (r Result) RelativeDeviation() float64 {
	if r.MonoAbsorbance == 0 {
		return 0
	}
	return r.Deviation / r.MonoAbsorbance
}

func (r Result) BandTransmittance() float64 {
	return math.Pow(10, -r.BandAbsorbance)
}

func (r Result) MonoTransmittance() float64 {
	return math.Pow(10, -r.MonoAbsorbance)
}

func (r Result) PercentDeviation() float64 {
	return 100 * r.RelativeDeviation()
}

func bracketInterval(lo, hi float64) string {
	return "[" + fmt.Sprintf("%g", lo) + ", " + fmt.Sprintf("%g", hi) + "] nm"
}
