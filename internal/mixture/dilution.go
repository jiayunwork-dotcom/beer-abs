package mixture

import "strconv"

// Dilute returns a copy of the mixture in which every concentration has been
// divided by the dilution factor, simulating adding solvent to the sample.
// The path length and stray fraction stay unchanged; the total absorbance
// shrinks by the same factor, preserving the Beer–Lambert linearity.
func (m Mixture) Dilute(factor float64) (Mixture, error) {
	if factor <= 0 {
		return Mixture{}, &DilutionError{Factor: factor, Reason: "稀释因子必须大于 0"}
	}
	out := m.Copy()
	for i := range out.Components {
		out.Components[i].Concentration /= factor
	}
	return out, nil
}

// DilutionSeries returns the original mixture and the first n−1 dilutions
// produced by repeatedly applying a factor, so the caller can watch A fall in
// a geometric progression.
func (m Mixture) DilutionSeries(factor float64, n int) ([]Mixture, error) {
	if n < 1 {
		return nil, &DilutionError{Factor: float64(n), Reason: "系列长度必须不小于 1"}
	}
	out := make([]Mixture, 0, n)
	cur := m.Copy()
	for i := 0; i < n; i++ {
		if i > 0 {
			var err error
			cur, err = cur.Dilute(factor)
			if err != nil {
				return nil, err
			}
		}
		out = append(out, cur.Copy())
	}
	return out, nil
}

// DilutionError describes an invalid dilution request.
type DilutionError struct {
	// Factor is the offending dilution factor.
	Factor float64

	// Reason is a human readable description of the violated constraint.
	Reason string
}

// Error implements the error interface.
func (e *DilutionError) Error() string {
	return "dilution factor " + strconv.FormatFloat(e.Factor, 'g', -1, 64) + ": " + e.Reason
}
