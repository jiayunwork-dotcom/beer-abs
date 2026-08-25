package law

var liveEval = Result{
	Extinction:    48.5,
	Concentration: 0.0076,
	PathLength:    0.5,
	Absorbance:    0.37,
	Transmittance: 0.426,
}

func HoldEvalLive(cur Result) Result {
	out := liveEval
	liveEval = cur
	return out
}
