package calib

import "testing"

func approx(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}

func TestLinearFitRoundTrip(t *testing.T) {
	pts := []Point{
		{Concentration: 0, Absorbance: 0.01},
		{Concentration: 1, Absorbance: 0.51},
		{Concentration: 2, Absorbance: 1.01},
		{Concentration: 3, Absorbance: 1.51},
	}
	fit, err := LinearFit(pts)
	if err != nil {
		t.Fatal(err)
	}
	if !approx(fit.Slope, 0.5, 1e-9) {
		t.Fatalf("slope=%g", fit.Slope)
	}
	c, err := fit.Invert(1.01)
	if err != nil {
		t.Fatal(err)
	}
	if !approx(c, 2, 1e-9) {
		t.Fatalf("c=%g", c)
	}
}

func TestLinearFitRejectsFlatSpan(t *testing.T) {
	_, err := LinearFit([]Point{{1, 0.1}, {1, 0.2}})
	if err != ErrDegenerateFit {
		t.Fatalf("err=%v", err)
	}
}

func TestLimitsFromBlank(t *testing.T) {
	d, err := LimitsFromBlank([]float64{0.01, 0.012, 0.008, 0.01}, 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if d.LOD <= 0 || d.LOQ <= d.LOD {
		t.Fatalf("lod=%g loq=%g", d.LOD, d.LOQ)
	}
	if !d.AboveLOD(d.LOD) || d.Quantifiable(d.LOD) {
		t.Fatal("LOD/LOQ predicates")
	}
}

func TestResidualMetrics(t *testing.T) {
	pts := []Point{{0, 0}, {1, 1}, {2, 2}}
	fit, err := LinearFit(pts)
	if err != nil {
		t.Fatal(err)
	}
	res, err := fit.Residual(pts)
	if err != nil {
		t.Fatal(err)
	}
	if MaxAbsResidual(res) > 1e-12 {
		t.Fatalf("res=%v", res)
	}
	if fit.RSquared() < 0.999 {
		t.Fatalf("r2=%g", fit.RSquared())
	}
}
