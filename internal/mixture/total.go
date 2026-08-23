package mixture

import "beer-abs/internal/law"

// validatePath checks the shared optical path; it reuses the law package's
// path-length rule so every entry point enforces L > 0 identically.
func validatePath(pathLength float64) error {
	return law.ValidatePathLength(pathLength)
}

// TotalAbsorbance returns A_tot = Σ ε_i·c_i·L. The mixture is validated
// first, then each component contributes its own product; the sum is taken in
// component order.
func (m Mixture) TotalAbsorbance() (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	return recallTotal(m.totalAbsorbanceUnchecked), nil
}

// totalAbsorbanceUnchecked sums the contributions assuming a validated
// mixture. It exists so derived quantities can avoid re-validating.
func (m Mixture) totalAbsorbanceUnchecked() float64 {
	var sum float64
	for i := range m.Components {
		sum += m.Components[i].AbsorbanceUnchecked(m.PathLength)
	}
	return sum
}

// ComponentAbsorbances returns the per-component contributions in order; the
// total equals their sum.
func (m Mixture) ComponentAbsorbances() ([]float64, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	out := make([]float64, len(m.Components))
	for i := range m.Components {
		a, err := m.Components[i].Absorbance(m.PathLength)
		if err != nil {
			return nil, err
		}
		out[i] = a
	}
	return out, nil
}

// WeightedFraction returns the fraction of the total absorbance contributed
// by component i, useful for interpreting mixture output.
func (m Mixture) WeightedFraction(i int) (float64, error) {
	total, err := m.TotalAbsorbance()
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	parts, err := m.ComponentAbsorbances()
	if err != nil {
		return 0, err
	}
	if i < 0 || i >= len(parts) {
		return 0, &IndexError{Index: i, Length: len(parts)}
	}
	return parts[i] / total, nil
}
