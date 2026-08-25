package mixture

import "strconv"

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

type DilutionError struct {
	Factor float64

	Reason string
}

func (e *DilutionError) Error() string {
	return "dilution factor " + strconv.FormatFloat(e.Factor, 'g', -1, 64) + ": " + e.Reason
}
