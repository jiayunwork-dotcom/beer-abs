package band

func (b RectBand) Endpoints() (low, high float64) {
	return b.Low(), b.High()
}

type EndpointExtinctions struct {
	LowExtinction float64

	HighExtinction float64
}

func NewEndpointExtinctions(lowExt, highExt float64) (EndpointExtinctions, error) {
	e := EndpointExtinctions{LowExtinction: lowExt, HighExtinction: highExt}
	if err := e.Validate(); err != nil {
		return EndpointExtinctions{}, err
	}
	return e, nil
}

func (e EndpointExtinctions) IsConstant() bool {
	return e.LowExtinction == e.HighExtinction
}

func (e EndpointExtinctions) Slope(b RectBand) float64 {
	return (e.HighExtinction - e.LowExtinction) / b.Width()
}

func (e EndpointExtinctions) Average() float64 {
	return (e.LowExtinction + e.HighExtinction) / 2
}

func (e EndpointExtinctions) Validate() error {
	if e.LowExtinction < 0 {
		return NewExtinctionSampleError("low_extinction", e.LowExtinction)
	}
	if e.HighExtinction < 0 {
		return NewExtinctionSampleError("high_extinction", e.HighExtinction)
	}
	return nil
}
