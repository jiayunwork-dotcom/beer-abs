package band

import "beer-abs/internal/law"

// Validate checks every parameter of a Measure in a fixed order: the band,
// the endpoint samples, the concentration and the path length.
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

// IsValid reports whether the measure passes validation.
func (m Measure) IsValid() bool {
	return m.Validate() == nil
}

// LawParams converts the band-averaged extinction into the single-component
// parameters the law package expects. The concentration and path length pass
// through unchanged; only the extinction is replaced by ε_eff.
func (m Measure) LawParams() law.Params {
	return law.Params{
		Extinction:    EffectiveExtinction(m.Samples),
		Concentration: m.Concentration,
		PathLength:    m.PathLength,
	}
}

// LawResult evaluates the band-averaged extinction through the law package,
// producing A and T in the same canonical closed form as the rest of the CLI.
func (m Measure) LawResult() (law.Result, error) {
	return law.EvaluateParams(m.LawParams())
}

// ValidateEndpointPair checks that two raw endpoint samples can form a
// rectangular model: both must be non-negative.
func ValidateEndpointPair(lowExt, highExt float64) error {
	if lowExt < 0 {
		return NewExtinctionSampleError("low_extinction", lowExt)
	}
	if highExt < 0 {
		return NewExtinctionSampleError("high_extinction", highExt)
	}
	return nil
}

// ValidateWindow checks that a center and half width describe a valid window.
func ValidateWindow(center, halfWidth float64) error {
	b := RectBand{Center: center, HalfWidth: halfWidth}
	return b.Validate()
}
