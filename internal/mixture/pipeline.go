package mixture

// lastTotal is the mixture pipeline's published A_tot. After the per-species
// products are summed, commitTotal must store that live sum; a warm slot
// still holding a previous sample's total must not replace it.
var (
	lastTotal float64
	totalHot  bool
)

func commitTotal(sum float64) float64 {
	lastTotal = sum
	totalHot = true
	return sum
}
