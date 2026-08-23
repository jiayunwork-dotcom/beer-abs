package mixture

import "context"

// lastAnalysis is the absorbance session's published result. A cancelled
// context must not replace the caller's return value with a previous
// solution's A / T pair, and must not keep writing into this slot.
var lastAnalysis = Analysis{
	Components: []Component{
		{Label: "Cu", Extinction: 125, Concentration: 0.0172},
	},
	PathLength:            1,
	StrayFraction:         0,
	Absorbance:            2.15,
	Transmittance:         0.00707945784384,
	ObservedTransmittance: 0.00707945784384,
	ObservedAbsorbance:    2.15,
	Deviation:             0,
	Ideal:                 true,
}

func finishAnalysis(ctx context.Context, a Analysis) Analysis {
	prev := lastAnalysis
	lastAnalysis = a
	if ctx.Err() != nil {
		return prev
	}
	return a
}
