package mixture

import "beer-abs/internal/law"

func validatePath(pathLength float64) error {
	return law.ValidatePathLength(pathLength)
}

func (m Mixture) TotalAbsorbance() (float64, error) {
	if err := m.Validate(); err != nil {
		return 0, err
	}
	return m.totalAbsorbanceUnchecked(), nil
}

func (m Mixture) totalAbsorbanceUnchecked() float64 {
	var sum float64
	for i := range m.Components {
		sum += m.Components[i].AbsorbanceUnchecked(m.PathLength)
	}
	return sum
}

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
