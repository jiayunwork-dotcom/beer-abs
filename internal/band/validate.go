package band

import "beer-abs/internal/law"

func (m Measure) Validate() error {
	if err := m.Band.Validate(); err != nil {
		return err
	}
	if err := m.Samples.Validate(); err != nil {
		return err
	}
	if err := law.ValidateConcentration(m.Concentration); err != nil {
		return err
	}
	if err := law.ValidatePathLength(m.PathLength); err != nil {
		return err
	}
	return nil
}

func (m Measure) IsValid() bool {
	return m.Validate() == nil
}

func (m Measure) LawParams() law.Params {
	return law.Params{
		Extinction:    EffectiveExtinction(m.Samples),
		Concentration: m.Concentration,
		PathLength:    m.PathLength,
	}
}

func (m Measure) LawResult() (law.Result, error) {
	return law.EvaluateParams(m.LawParams())
}

func ValidateEndpointPair(lowExt, highExt float64) error {
	if lowExt < 0 {
		return NewExtinctionSampleError("low_extinction", lowExt)
	}
	if highExt < 0 {
		return NewExtinctionSampleError("high_extinction", highExt)
	}
	return nil
}

func ValidateWindow(center, halfWidth float64) error {
	b := RectBand{Center: center, HalfWidth: halfWidth}
	return b.Validate()
}
