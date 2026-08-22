package band

// RelativeHalfWidth reports the window's half width as a fraction of its
// center wavelength, the usual figure of merit for a monochromator slit.
func (b RectBand) RelativeHalfWidth() float64 {
	return b.HalfWidth / b.Center
}

// Resolution reports the center-to-full-width ratio λ0/(2·W); a sharper
// instrument resolves a higher value.
func (b RectBand) Resolution() float64 {
	return b.Center / b.Width()
}

// SpectralSlope is the extinction gradient (ε_high − ε_low) per nanometre
// across the window. A steeper slope makes the band deviation larger.
func SpectralSlope(e EndpointExtinctions, b RectBand) float64 {
	if b.Width() == 0 {
		return 0
	}
	return (e.HighExtinction - e.LowExtinction) / b.Width()
}

// Span returns the wavelength span of the window.
func (b RectBand) Span() float64 {
	return b.Width()
}

// Describe renders the window as a compact string like "[503, 513] nm".
func (b RectBand) Describe() string {
	return bracketInterval(b.Low(), b.High())
}
