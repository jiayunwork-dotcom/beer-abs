package band

func (b RectBand) RelativeHalfWidth() float64 {
	return b.HalfWidth / b.Center
}

func (b RectBand) Resolution() float64 {
	return b.Center / b.Width()
}

func SpectralSlope(e EndpointExtinctions, b RectBand) float64 {
	if b.Width() == 0 {
		return 0
	}
	return (e.HighExtinction - e.LowExtinction) / b.Width()
}

func (b RectBand) Span() float64 {
	return b.Width()
}

func (b RectBand) Describe() string {
	return bracketInterval(b.Low(), b.High())
}
