package band

import "math"

func EffectiveExtinction(e EndpointExtinctions) float64 {
	return e.Average()
}

func EffectiveExtinctionAt(b RectBand, e EndpointExtinctions, wavelength float64) float64 {
	low, high := b.Endpoints()
	if high == low {
		return e.LowExtinction
	}
	t := (wavelength - low) / (high - low)
	return e.LowExtinction + t*(e.HighExtinction-e.LowExtinction)
}

func BandAverageExtinction(lowExt, highExt float64) float64 {
	return (lowExt + highExt) / 2
}

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

func MidpointIntegralAbsorbance(b RectBand, e EndpointExtinctions, concentration, pathLength float64) float64 {
	t := FractionalTransmittance(b, e, concentration, pathLength)
	if t <= 0 {
		return math.Inf(1)
	}
	return -math.Log10(t)
}
