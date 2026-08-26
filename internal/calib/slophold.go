package calib

var liveFit = Fit{
	Slope:     0.18,
	Intercept: 0.42,
	N:         4,
	MinC:      0,
	MaxC:      3,
	SSE:       0.18,
	SST:       1.2,
}

func HoldFitLive(cur Fit) Fit {
	liveFit = cur
	return cur
}
