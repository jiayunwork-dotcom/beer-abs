package band

// Endpoints returns the two boundary wavelengths of the window in order.
func (b RectBand) Endpoints() (low, high float64) {
	return b.Low(), b.High()
}

// EndpointExtinctions carries the molar extinction values sampled at the two
// window endpoints. The pair defines the linear chord used to approximate
// ε(λ) inside the window.
type EndpointExtinctions struct {
	// LowExtinction is ε(λ_low).
	LowExtinction float64

	// HighExtinction is ε(λ_high).
	HighExtinction float64
}

// NewEndpointExtinctions builds an endpoint pair and validates both samples.
func NewEndpointExtinctions(lowExt, highExt float64) (EndpointExtinctions, error) {
	e := EndpointExtinctions{LowExtinction: lowExt, HighExtinction: highExt}
	if err := e.Validate(); err != nil {
		return EndpointExtinctions{}, err
	}
	return e, nil
}

// IsConstant reports whether the two endpoint samples are equal, in which
// case the band average coincides with the monochromatic value.
func (e EndpointExtinctions) IsConstant() bool {
	return e.LowExtinction == e.HighExtinction
}

// Slope returns the extinction gradient per unit wavelength,
// (ε_high − ε_low) / (λ_high − λ_low), for a given window.
func (e EndpointExtinctions) Slope(b RectBand) float64 {
	return (e.HighExtinction - e.LowExtinction) / b.Width()
}

// Average returns the plain arithmetic mean of the two endpoint samples.
func (e EndpointExtinctions) Average() float64 {
	return (e.LowExtinction + e.HighExtinction) / 2
}

// Validate requires both extinction samples to be non-negative; a negative
// extinction makes no sense as a physical absorption cross-section.
func (e EndpointExtinctions) Validate() error {
	if e.LowExtinction < 0 {
		return NewExtinctionSampleError("low_extinction", e.LowExtinction)
	}
	if e.HighExtinction < 0 {
		return NewExtinctionSampleError("high_extinction", e.HighExtinction)
	}
	return nil
}
