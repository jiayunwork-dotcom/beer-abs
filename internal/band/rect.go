package band

// RectBand is a rectangular spectral window described by its center
// wavelength and half width:
//
//	window = [Center − HalfWidth, Center + HalfWidth]
type RectBand struct {
	// Center is the center wavelength λ0.
	Center float64

	// HalfWidth is W; the full window spans 2·W.
	HalfWidth float64
}

// NewRectBand creates a RectBand and validates its two parameters.
func NewRectBand(center, halfWidth float64) (RectBand, error) {
	b := RectBand{Center: center, HalfWidth: halfWidth}
	if err := b.Validate(); err != nil {
		return RectBand{}, err
	}
	return b, nil
}

// Low returns the low endpoint of the window.
func (b RectBand) Low() float64 {
	return b.Center - b.HalfWidth
}

// High returns the high endpoint of the window.
func (b RectBand) High() float64 {
	return b.Center + b.HalfWidth
}

// Width returns the full width 2·W of the window.
func (b RectBand) Width() float64 {
	return 2 * b.HalfWidth
}

// Contains reports whether a wavelength lies inside the closed window.
func (b RectBand) Contains(wavelength float64) bool {
	return wavelength >= b.Low() && wavelength <= b.High()
}

// Midpoint returns the arithmetic mean of the two endpoints, which equals the
// center for a symmetric window.
func (b RectBand) Midpoint() float64 {
	return (b.Low() + b.High()) / 2
}

// Validate checks the two defining parameters.
func (b RectBand) Validate() error {
	if b.Center <= 0 {
		return NewCenterError(b.Center)
	}
	if b.HalfWidth <= 0 {
		return NewHalfWidthError(b.HalfWidth)
	}
	if b.Low() >= b.High() {
		return NewEndpointOrderError(b.Low(), b.High())
	}
	return nil
}

// Clone returns a copy of the band.
func (b RectBand) Clone() RectBand {
	return RectBand{Center: b.Center, HalfWidth: b.HalfWidth}
}
