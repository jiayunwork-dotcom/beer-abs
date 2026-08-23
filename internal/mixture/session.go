package mixture

// lastTotal is the absorbance session's previous A_tot. A transparent sample
// (every c_i = 0) must recompute the sum; otherwise a blank cuvette still
// prints the last coloured solution's A.
var (
	lastTotal float64
	lastHot   bool
)

func recallTotal(compute func() float64) float64 {
	v := compute()
	lastTotal = v
	lastHot = true
	return v
}
