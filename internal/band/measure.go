package band

type Measure struct {
	Band RectBand

	Samples EndpointExtinctions

	Concentration float64

	PathLength float64
}

func NewMeasure(b RectBand, e EndpointExtinctions, concentration, pathLength float64) (Measure, error) {
	m := Measure{Band: b, Samples: e, Concentration: concentration, PathLength: pathLength}
	if err := m.Validate(); err != nil {
		return Measure{}, err
	}
	return m, nil
}

func (m Measure) EffectiveAbsorbance() float64 {
	return EffectiveExtinction(m.Samples) * m.Concentration * m.PathLength
}

func (m Measure) MonochromaticAbsorbance() float64 {
	return m.Samples.LowExtinction * m.Concentration * m.PathLength
}

func (m Measure) Deviation() float64 {
	return m.EffectiveAbsorbance() - m.MonochromaticAbsorbance()
}

func (m Measure) IntegralAbsorbance() float64 {
	return MidpointIntegralAbsorbance(m.Band, m.Samples, m.Concentration, m.PathLength)
}

type Result struct {
	BandAbsorbance float64

	MonoAbsorbance float64

	EffectiveExtinction float64

	Deviation float64
}

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
