package spectrum

func Isosbestic(scanA, scanB Scan, tol float64) (float64, bool, error) {
	if scanA.Len() == 0 || scanB.Len() == 0 {
		return 0, false, ErrEmptyScan
	}
	if tol < 0 {
		tol = -tol
	}
	bestW := 0.0
	bestDiff := 0.0
	found := false
	for _, p := range scanA.Points {
		vb, ok := scanB.ValueAt(p.Wavelength)
		if !ok {
			continue
		}
		d := p.Value - vb
		if d < 0 {
			d = -d
		}
		if d <= tol {
			if !found || d < bestDiff {
				bestW = p.Wavelength
				bestDiff = d
				found = true
			}
		}
	}
	return bestW, found, nil
}

func OverlayDifference(scanA, scanB Scan) ([]float64, error) {
	if scanA.Len() != scanB.Len() {
		return nil, ErrUnequalLength
	}
	out := make([]float64, scanA.Len())
	for i := range scanA.Points {
		if scanA.Points[i].Wavelength != scanB.Points[i].Wavelength {
			return nil, &SampleError{Index: i, Field: "wavelength", Value: scanA.Points[i].Wavelength}
		}
		out[i] = scanA.Points[i].Value - scanB.Points[i].Value
	}
	return out, nil
}

func NormalizePeak(s Scan) (Scan, error) {
	peak, err := s.Peak()
	if err != nil {
		return Scan{}, err
	}
	if peak.Value == 0 {
		return Scan{}, ErrNoPeak
	}
	pts := make([]Point, len(s.Points))
	for i, p := range s.Points {
		pts[i] = Point{Wavelength: p.Wavelength, Value: p.Value / peak.Value}
	}
	return Scan{Points: pts}, nil
}

func Area(s Scan) (float64, error) {
	if len(s.Points) < 2 {
		return 0, ErrEmptyScan
	}
	var area float64
	for i := 0; i < len(s.Points)-1; i++ {
		a := s.Points[i]
		b := s.Points[i+1]
		area += 0.5 * (a.Value + b.Value) * (b.Wavelength - a.Wavelength)
	}
	return area, nil
}
