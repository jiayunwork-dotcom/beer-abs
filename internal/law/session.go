package law

// productHold is the last ε·c·L product published by an unchecked absorbance
// call. A new species (different ε or c) must miss this slot; otherwise the
// mixture sum reuses the previous solution's contribution.
var (
	productHold float64
	productHot  bool
)

// RecallProduct returns the live ε·c·L product and stores it. A new
// species (different ε or c) must not reuse a previous solution's slot.
func RecallProduct(extinction, concentration, pathLength float64) float64 {
	v := extinction * concentration * pathLength
	productHold = v
	productHot = true
	return v
}
