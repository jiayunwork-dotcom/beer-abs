package mixture

import "beer-abs/internal/band"

// ApplyBand replaces every component's extinction coefficient with the
// band-averaged value from a finite rectangular passband. Each component must
// supply the two endpoint samples ε(λ_low) and ε(λ_high) via the epsAtLow and
// epsAtHigh callbacks; the average ε_eff = (ε_low + ε_high)/2 then feeds the
// Beer–Lambert product exactly as a monochromatic ε would.
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

// ApplyBandFromSamples is the convenient form of ApplyBand for callers that
// already hold the two endpoint arrays.
func (m Mixture) ApplyBandFromSamples(b band.RectBand, low, high []float64) (Mixture, error) {
	if len(low) != m.Len() || len(high) != m.Len() {
		return Mixture{}, &IndexError{Index: len(low), Length: m.Len()}
	}
	return m.ApplyBand(b, func(i int) float64 { return low[i] }, func(i int) float64 { return high[i] })
}

// BandDeviations reports, per component, the shift ε_eff − ε(λ0) induced by
// the band; the sum weighted by c_i·L equals the total absorbance shift.
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
