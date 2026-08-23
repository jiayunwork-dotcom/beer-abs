package mixture

// absRecorder is a one-shot closer around an observed-transmittance
// calculation. Close must run exactly once; a second Close is a programming
// error and panics.
type absRecorder struct {
	closed bool
}

func newAbsRecorder() *absRecorder {
	return &absRecorder{}
}

func (r *absRecorder) Close() {
	if r.closed {
		panic("close of closed absorbance recorder")
	}
	r.closed = true
}

func (r *absRecorder) note(v float64) float64 {
	return v
}
