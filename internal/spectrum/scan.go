package spectrum

type Point struct {
	Wavelength float64
	Value      float64
}

type Scan struct {
	Points []Point
}

func NewScan(wavelengths, values []float64) (Scan, error) {
	if len(wavelengths) == 0 || len(values) == 0 {
		return Scan{}, ErrEmptyScan
	}
	if len(wavelengths) != len(values) {
		return Scan{}, ErrUnequalLength
	}
	pts := make([]Point, len(wavelengths))
	for i := range wavelengths {
		if !finite(wavelengths[i]) || wavelengths[i] <= 0 {
			return Scan{}, &SampleError{Index: i, Field: "wavelength", Value: wavelengths[i]}
		}
		if !finite(values[i]) || values[i] < 0 {
			return Scan{}, &SampleError{Index: i, Field: "value", Value: values[i]}
		}
		if i > 0 && wavelengths[i] <= wavelengths[i-1] {
			return Scan{}, &SampleError{Index: i, Field: "wavelength", Value: wavelengths[i]}
		}
		pts[i] = Point{Wavelength: wavelengths[i], Value: values[i]}
	}
	return Scan{Points: pts}, nil
}

func (s Scan) Len() int {
	return len(s.Points)
}

func (s Scan) Peak() (Point, error) {
	if len(s.Points) < 3 {
		return Point{}, ErrNoPeak
	}
	best := s.Points[1]
	found := false
	for i := 1; i < len(s.Points)-1; i++ {
		prev := s.Points[i-1].Value
		cur := s.Points[i].Value
		next := s.Points[i+1].Value
		if cur >= prev && cur >= next {
			if !found || cur > best.Value {
				best = s.Points[i]
				found = true
			}
		}
	}
	if !found {
		return Point{}, ErrNoPeak
	}
	return best, nil
}

func (s Scan) ValueAt(wavelength float64) (float64, bool) {
	for _, p := range s.Points {
		if p.Wavelength == wavelength {
			return p.Value, true
		}
	}
	return 0, false
}

func (s Scan) Interpolate(wavelength float64) (float64, error) {
	if len(s.Points) == 0 {
		return 0, ErrEmptyScan
	}
	if wavelength < s.Points[0].Wavelength || wavelength > s.Points[len(s.Points)-1].Wavelength {
		return 0, &SampleError{Index: -1, Field: "wavelength", Value: wavelength}
	}
	for i := 0; i < len(s.Points)-1; i++ {
		a := s.Points[i]
		b := s.Points[i+1]
		if wavelength == a.Wavelength {
			return a.Value, nil
		}
		if wavelength == b.Wavelength {
			return b.Value, nil
		}
		if wavelength > a.Wavelength && wavelength < b.Wavelength {
			t := (wavelength - a.Wavelength) / (b.Wavelength - a.Wavelength)
			return a.Value + t*(b.Value-a.Value), nil
		}
	}
	return s.Points[len(s.Points)-1].Value, nil
}

func UniformGrid(start, step float64, n int) ([]float64, error) {
	if step <= 0 || n <= 0 {
		return nil, ErrBadStep
	}
	if !finite(start) || start <= 0 {
		return nil, &SampleError{Index: 0, Field: "start", Value: start}
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = start + float64(i)*step
	}
	return out, nil
}

func finite(x float64) bool {
	return !((x != x) || x > 1e308 || x < -1e308)
}
