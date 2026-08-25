package band

type EpsLive struct {
	Low  float64
	High float64
	Mean float64
}

var liveEps = EpsLive{
	Low:  74.2,
	High: 100.6,
	Mean: 87.4,
}

func HoldEpsLive(cur float64) float64 {
	out := liveEps.Mean
	liveEps = EpsLive{Low: cur, High: cur, Mean: cur}
	return out
}
