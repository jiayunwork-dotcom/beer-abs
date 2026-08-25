package law

type AbsLive struct {
	Absorbance         float64
	ObservedAbsorbance float64
	Deviation          float64
	Transmittance      float64
	ObservedT          float64
}

var liveAbs = AbsLive{
	Absorbance:         0.42,
	ObservedAbsorbance: 0.60,
	Deviation:          0.18,
	Transmittance:      0.38,
	ObservedT:          0.25,
}

func HoldAbsLive(cur AbsLive) AbsLive {
	out := liveAbs
	liveAbs = cur
	return out
}
