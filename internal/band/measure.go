package band

// Measure couples a rectangular band with the sample parameters needed to
// compute a finite-bandwidth absorbance.
type Measure struct {
	// Band is the spectral window.
	Band RectBand

	// Samples holds the extinction at the two window endpoints.
	Samples EndpointExtinctions

	// Concentration is c (mol·L⁻¹).
	Concentration float64

	// PathLength is L (cm).
	PathLength float64
}

// NewMeasure validates the band, the endpoint samples and the sample
// parameters, then returns a Measure.
func NewMeasure(b RectBand, e EndpointExtinctions, concentration, pathLength float64) (Measure, error) {
	m := Measure{Band: b, Samples: e, Concentration: concentration, PathLength: pathLength}
	if err := m.Validate(); err != nil {
		return Measure{}, err
	}
	return m, nil
}

// EffectiveAbsorbance returns A_band = ε_eff·c·L using the band-averaged
// extinction.
func (m Measure) EffectiveAbsorbance() float64 {
	return EffectiveExtinction(m.Samples) * m.Concentration * m.PathLength
}

// MonochromaticAbsorbance returns A(λ_low) = ε(λ_low)·c·L, the reading an
// ideal single-wavelength instrument would report if it were tuned to the low
// end of the window. Because the band-averaged extinction ε_eff is the mean
// of the two endpoint samples, the difference ε_eff − ε(λ_low) is exactly
// half of (ε_high − ε_low): the finite bandwidth shifts the reading whenever
// the extinction varies across the window.
func (m Measure) MonochromaticAbsorbance() float64 {
	return m.Samples.LowExtinction * m.Concentration * m.PathLength
}

// Deviation is A_band − A(λ0), the finite-bandwidth shift away from the
// single-wavelength reading. It is non-zero exactly when the extinction
// varies across the window.
func (m Measure) Deviation() float64 {
	return m.EffectiveAbsorbance() - m.MonochromaticAbsorbance()
}

// IntegralAbsorbance is the midpoint-rule answer −log10(⟨T⟩), kept for
// comparison with EffectiveAbsorbance.
func (m Measure) IntegralAbsorbance() float64 {
	return MidpointIntegralAbsorbance(m.Band, m.Samples, m.Concentration, m.PathLength)
}

// Result is a compact summary of a band measurement.
type Result struct {
	// BandAbsorbance is the finite-bandwidth reading.
	BandAbsorbance float64

	// MonoAbsorbance is the single-wavelength reading at the center.
	MonoAbsorbance float64

	// EffectiveExtinction is the band-averaged ε.
	EffectiveExtinction float64

	// Deviation is BandAbsorbance − MonoAbsorbance.
	Deviation float64
}

// Analyze runs the complete band pipeline.
func (m Measure) Analyze() (Result, error) {
	if err := m.Validate(); err != nil {
		return Result{}, err
	}
	return Result{
		BandAbsorbance:      m.EffectiveAbsorbance(),
		MonoAbsorbance:      m.MonochromaticAbsorbance(),
		EffectiveExtinction: EffectiveExtinction(m.Samples),
		Deviation:           m.Deviation(),
	}, nil
}
