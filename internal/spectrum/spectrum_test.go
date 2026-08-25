package spectrum

import "testing"

func TestScanPeakAndNormalize(t *testing.T) {
	s, err := NewScan([]float64{400, 450, 500, 550, 600}, []float64{0.1, 0.4, 1.0, 0.5, 0.2})
	if err != nil {
		t.Fatal(err)
	}
	peak, err := s.Peak()
	if err != nil {
		t.Fatal(err)
	}
	if peak.Wavelength != 500 || peak.Value != 1.0 {
		t.Fatalf("peak=%v", peak)
	}
	n, err := NormalizePeak(s)
	if err != nil {
		t.Fatal(err)
	}
	p2, _ := n.Peak()
	if p2.Value != 1.0 {
		t.Fatalf("norm peak=%v", p2)
	}
}

func TestIsosbesticCrossing(t *testing.T) {
	a, err := NewScan([]float64{400, 450, 500}, []float64{0.8, 0.5, 0.2})
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewScan([]float64{400, 450, 500}, []float64{0.2, 0.5, 0.9})
	if err != nil {
		t.Fatal(err)
	}
	w, ok, err := Isosbestic(a, b, 1e-9)
	if err != nil || !ok || w != 450 {
		t.Fatalf("w=%g ok=%v err=%v", w, ok, err)
	}
}

func TestInterpolateInterior(t *testing.T) {
	s, err := NewScan([]float64{400, 500}, []float64{0, 1})
	if err != nil {
		t.Fatal(err)
	}
	v, err := s.Interpolate(450)
	if err != nil {
		t.Fatal(err)
	}
	if v < 0.49 || v > 0.51 {
		t.Fatalf("v=%g", v)
	}
}

func TestAreaTrapezoid(t *testing.T) {
	s, err := NewScan([]float64{0.1, 0.2, 0.3}, []float64{1, 1, 1})
	if err != nil {
		t.Fatal(err)
	}
	a, err := Area(s)
	if err != nil {
		t.Fatal(err)
	}
	if a < 0.199 || a > 0.201 {
		t.Fatalf("area=%g", a)
	}
}
