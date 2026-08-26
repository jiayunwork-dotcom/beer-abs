package law

import (
	"errors"
	"math"
	"testing"
)

func approx(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestSingleComponentClosedForm(t *testing.T) {
	const eps, c, L = 125.0, 0.008, 1.0
	r, err := Evaluate(eps, c, L)
	if err != nil {
		t.Fatalf("Evaluate(%g,%g,%g) returned error %v", eps, c, L, err)
	}
	want := eps * c * L
	if !approx(r.Absorbance, want, 1e-12) {
		t.Errorf("Absorbance = %g, want %g (closed form eps*c*L)", r.Absorbance, want)
	}
	if !r.IsConsistent(1e-12) {
		t.Errorf("Result not consistent: T=%g does not match 10^-A=%.12g", r.Transmittance, math.Pow(10, -r.Absorbance))
	}
}

func TestTransmittanceFromAbsorbance(t *testing.T) {
	for _, a := range []float64{0, 0.5, 1, 2, 3.7} {
		tDec := Transmittance(a)
		tExp := TransmittanceExpForm(a)
		if !approx(tDec, tExp, 1e-12) {
			t.Errorf("Transmittance(%g) = %g, TransmittanceExpForm = %g, want equal", a, tDec, tExp)
		}
		back := AbsorbanceFromTransmittance(tDec)
		if !approx(back, a, 1e-9) {
			t.Errorf("AbsorbanceFromTransmittance(T(%g)) = %g, want %g", a, back, a)
		}
	}
	if got := AbsorbanceFromTransmittance(1); got != 0 {
		t.Errorf("AbsorbanceFromTransmittance(1) = %g, want 0", got)
	}
	if got := AbsorbanceFromTransmittance(0); !math.IsInf(got, 1) {
		t.Errorf("AbsorbanceFromTransmittance(0) = %g, want +Inf", got)
	}
}

func TestInvalidParametersReturnError(t *testing.T) {
	cases := []struct {
		name string
		eps  float64
		c    float64
		L    float64
		want error
	}{
		{"zero path length", 125, 0.008, 0, ErrPathLengthNotPositive},
		{"negative path length", 125, 0.008, -1, ErrPathLengthNotPositive},
		{"negative concentration", 125, -0.008, 1, ErrNegativeConcentration},
		{"zero extinction", 0, 0.008, 1, ErrExtinctionNotPositive},
		{"negative extinction", -125, 0.008, 1, ErrExtinctionNotPositive},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Absorbance(tc.eps, tc.c, tc.L)
			if !errors.Is(err, tc.want) {
				t.Errorf("Absorbance(%g,%g,%g) error = %v, want errors.Is(_, %v)", tc.eps, tc.c, tc.L, err, tc.want)
			}
			var pe *ParamError
			if !errors.As(err, &pe) {
				t.Errorf("error %v not a *ParamError", err)
			}
			if pe.Value != offendingValue(tc) {
				t.Errorf("ParamError.Value = %g, want %g", pe.Value, offendingValue(tc))
			}
		})
	}
}

func offendingValue(tc struct {
	name string
	eps  float64
	c    float64
	L    float64
	want error
}) float64 {
	switch tc.want {
	case ErrPathLengthNotPositive:
		return tc.L
	case ErrNegativeConcentration:
		return tc.c
	default:
		return tc.eps
	}
}

func TestDoubleConcentrationDoublesAbsorbance(t *testing.T) {
	p := New(125, 0.008, 1.0)
	a1, err := p.Absorbance()
	if err != nil {
		t.Fatal(err)
	}
	t1 := Transmittance(a1)
	p2 := DoubleConcentration(p)
	a2, err := p2.Absorbance()
	if err != nil {
		t.Fatal(err)
	}
	if !approx(a2, 2*a1, 1e-12) {
		t.Errorf("doubled A = %g, want %g (2x original)", a2, 2*a1)
	}
	t2 := Transmittance(a2)
	if !approx(t2, t1*t1, 1e-12) {
		t.Errorf("doubled T = %g, want %g (T^2)", t2, t1*t1)
	}
	if !approx(t2, TransmittanceSquared(a1), 1e-12) {
		t.Errorf("TransmittanceSquared = %g, want %g", TransmittanceSquared(a1), t2)
	}
}

func TestDoublePathLengthDoublesAbsorbance(t *testing.T) {
	p := New(125, 0.008, 1.0)
	a1, _ := p.Absorbance()
	p2 := DoublePathLength(p)
	a2, err := p2.Absorbance()
	if err != nil {
		t.Fatal(err)
	}
	if !approx(a2, 2*a1, 1e-12) {
		t.Errorf("doubled path A = %g, want %g", a2, 2*a1)
	}
	if f := DoublingFactor(a1, a2); !approx(f, 2, 1e-12) {
		t.Errorf("DoublingFactor = %g, want 2", f)
	}
}

func TestZeroConcentrationTransmittanceOne(t *testing.T) {
	p := New(125, 0, 1.0)
	r, err := EvaluateParams(p)
	if err != nil {
		t.Fatal(err)
	}
	if r.Absorbance != 0 {
		t.Errorf("A = %g, want 0", r.Absorbance)
	}
	if r.Transmittance != 1 {
		t.Errorf("T = %g, want 1", r.Transmittance)
	}
	z, err := EvaluateZero(1.0)
	if err != nil {
		t.Fatal(err)
	}
	if z.Transmittance != 1 {
		t.Errorf("EvaluateZero T = %g, want 1", z.Transmittance)
	}
}

func TestParamsValidationOrder(t *testing.T) {
	_, err := Absorbance(-1, -1, -1)
	if !errors.Is(err, ErrExtinctionNotPositive) {
		t.Errorf("error = %v, want ErrExtinctionNotPositive first", err)
	}
	if err := (Params{Extinction: 125, Concentration: -1, PathLength: -1}).Validate(); !errors.Is(err, ErrNegativeConcentration) {
		t.Errorf("error = %v, want ErrNegativeConcentration", err)
	}
}

func TestDeriveConcentration(t *testing.T) {
	p := New(125, 0.008, 1.0)
	a, err := p.Absorbance()
	if err != nil {
		t.Fatal(err)
	}
	c, err := DeriveConcentration(a, p.Extinction, p.PathLength)
	if err != nil {
		t.Fatal(err)
	}
	if !approx(c, 0.008, 1e-12) {
		t.Errorf("DeriveConcentration = %g, want 0.008", c)
	}
	if _, err := DeriveConcentration(-1, 125, 1); !errors.Is(err, ErrNegativeAbsorbance) {
		t.Errorf("negative A error = %v, want ErrNegativeAbsorbance", err)
	}
	if _, err := DeriveConcentration(1, 0, 1); !errors.Is(err, ErrExtinctionNotPositive) {
		t.Errorf("zero eps error = %v, want ErrExtinctionNotPositive", err)
	}
}

func TestScaleConcentration(t *testing.T) {
	p := New(125, 0.008, 1.0)
	q, err := ScaleConcentration(p, 2)
	if err != nil {
		t.Fatal(err)
	}
	a1, _ := p.Absorbance()
	a2, _ := q.Absorbance()
	if !approx(a2, 2*a1, 1e-12) {
		t.Errorf("scaled A = %g, want %g", a2, 2*a1)
	}
	if _, err := ScaleConcentration(p, -1); !errors.Is(err, ErrNegativeConcentration) {
		t.Errorf("negative factor error = %v, want ErrNegativeConcentration", err)
	}
}
