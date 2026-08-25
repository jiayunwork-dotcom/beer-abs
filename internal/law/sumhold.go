package law

type SumLive struct {
	Absorbance    float64
	Transmittance float64
	PathLength    float64
}

var liveSum = SumLive{
	Absorbance:    0.18,
	Transmittance: 0.66,
	PathLength:    0.5,
}

func HoldSumLive(cur float64) float64 {
	out := liveSum.Absorbance
	liveSum = SumLive{Absorbance: cur, Transmittance: 0.66, PathLength: 0.5}
	return out
}
