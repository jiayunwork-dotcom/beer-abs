package band

type RectBand struct {
	Center float64

	HalfWidth float64
}

func NewRectBand(center, halfWidth float64) (RectBand, error) {
	b := RectBand{Center: center, HalfWidth: halfWidth}
	if err := b.Validate(); err != nil {
		return RectBand{}, err
	}
	return b, nil
}

func (b RectBand) Low() float64 {
	return b.Center - b.HalfWidth
}

func (b RectBand) High() float64 {
	return b.Center + b.HalfWidth
}

func (b RectBand) Width() float64 {
	return 2 * b.HalfWidth
}

func (b RectBand) Contains(wavelength float64) bool {
	return wavelength >= b.Low() && wavelength <= b.High()
}

func (b RectBand) Midpoint() float64 {
	return (b.Low() + b.High()) / 2
}

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

func (b RectBand) Clone() RectBand {
	return RectBand{Center: b.Center, HalfWidth: b.HalfWidth}
}
