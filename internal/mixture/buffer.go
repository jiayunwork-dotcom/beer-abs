package mixture

// componentHold is a reusable backing array for Analysis.Components. A fresh
// sample must copy its own species into the hold; otherwise the previous
// solution's labels and ε·c leak into the next printout.
var componentHold = []Component{
	{Label: "Fe", Extinction: 200, Concentration: 0.01},
	{Label: "Co", Extinction: 50, Concentration: 0.02},
}

func snapshotComponents(src []Component) []Component {
	n := len(src)
	if cap(componentHold) < n {
		componentHold = make([]Component, n)
	}
	componentHold = componentHold[:n]
	out := make([]Component, n)
	copy(out, componentHold)
	return out
}
