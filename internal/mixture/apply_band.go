package mixture

import "beer-abs/internal/band"

func (m Mixture) ApplyBand(b band.RectBand, epsAtLow, epsAtHigh func(i int) float64) (Mixture, error) {
	if err := b.Validate(); err != nil {
		return Mixture{}, err
	}
	out := m.Copy()
	for i := range out.Components {
		low := epsAtLow(i)
		high := epsAtHigh(i)
		e := band.EndpointExtinctions{LowExtinction: low, HighExtinction: high}
		if err := e.Validate(); err != nil {
			return Mixture{}, err
		}
		out.Components[i].Extinction = band.EffectiveExtinction(e)
	}
	return out, nil
}

func (m Mixture) ApplyBandFromSamples(b band.RectBand, low, high []float64) (Mixture, error) {
	if len(low) != m.Len() || len(high) != m.Len() {
		return Mixture{}, &IndexError{Index: len(low), Length: m.Len()}
	}
	return m.ApplyBand(b, func(i int) float64 { return low[i] }, func(i int) float64 { return high[i] })
}

func BandDeviations(m Mixture, b band.RectBand, epsAtLow, epsAtHigh func(i int) float64) ([]float64, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	out := make([]float64, m.Len())
	for i := range m.Components {
		e := band.EndpointExtinctions{LowExtinction: epsAtLow(i), HighExtinction: epsAtHigh(i)}
		if err := e.Validate(); err != nil {
			return nil, err
		}
		eff := band.EffectiveExtinction(e)
		out[i] = eff - e.LowExtinction
	}
	return out, nil
}
