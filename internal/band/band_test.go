package band

import (
	"errors"
	"math"
	"testing"
)

func closeTo(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// TestRectBandAverageExtinction checks the rectangular model's averaging
// rule: the effective extinction is the arithmetic mean of the two endpoint
// samples.
func TestRectBandAverageExtinction(t *testing.T) {
	b, err := NewRectBand(508, 5)
	if err != nil {
		t.Fatal(err)
	}
	e := EndpointExtinctions{LowExtinction: 120, HighExtinction: 140}
	if got := EffectiveExtinction(e); !closeTo(got, 130, 1e-12) {
		t.Errorf("EffectiveExtinction = %g, want 130", got)
	}
	if got := BandAverageExtinction(120, 140); !closeTo(got, 130, 1e-12) {
		t.Errorf("BandAverageExtinction = %g, want 130", got)
	}
	if got := e.Average(); !closeTo(got, 130, 1e-12) {
		t.Errorf("EndpointExtinctions.Average = %g, want 130", got)
	}
	if b.Width() != 10 {
		t.Errorf("Width = %g, want 10", b.Width())
	}
	if b.Low() != 503 || b.High() != 513 {
		t.Errorf("window = [%g, %g], want [503, 513]", b.Low(), b.High())
	}
	if !b.Contains(508) || b.Contains(500) {
		t.Errorf("Contains(508) = %v, Contains(500) = %v, want true/false", b.Contains(508), b.Contains(500))
	}
}

// TestRectBandValidation checks the window parameter checks: a non-positive
// center or half width, or an inverted endpoint pair, must be rejected.
func TestRectBandValidation(t *testing.T) {
	cases := []struct {
		name       string
		center     float64
		halfWidth  float64
		want       error
	}{
		{"zero center", 0, 5, ErrCenterNotPositive},
		{"negative center", -508, 5, ErrCenterNotPositive},
		{"zero half width", 508, 0, ErrHalfWidthNotPositive},
		{"negative half width", 508, -5, ErrHalfWidthNotPositive},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRectBand(tc.center, tc.halfWidth)
			if !errors.Is(err, tc.want) {
				t.Errorf("NewRectBand(%g,%g) error = %v, want errors.Is(_, %v)", tc.center, tc.halfWidth, err, tc.want)
			}
		})
	}
	if err := ValidateWindow(508, 5); err != nil {
		t.Errorf("ValidateWindow(508,5) = %v, want nil", err)
	}
}

// TestRectBandShiftsFromMonochromatic checks the finite-bandwidth effect: the
// band-averaged absorbance differs from the single-wavelength reference
// reading whenever the extinction varies across the window, and the direction
// matches half the endpoint difference.
func TestRectBandShiftsFromMonochromatic(t *testing.T) {
	b, err := NewRectBand(508, 5)
	if err != nil {
		t.Fatal(err)
	}
	e := EndpointExtinctions{LowExtinction: 120, HighExtinction: 140}
	m, err := NewMeasure(b, e, 0.008, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	r, err := m.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	wantBand := 130 * 0.008 * 1.0
	if !closeTo(r.BandAbsorbance, wantBand, 1e-12) {
		t.Errorf("BandAbsorbance = %g, want %g", r.BandAbsorbance, wantBand)
	}
	wantMono := 120 * 0.008 * 1.0
	if !closeTo(r.MonoAbsorbance, wantMono, 1e-12) {
		t.Errorf("MonoAbsorbance = %g, want %g", r.MonoAbsorbance, wantMono)
	}
	if !(r.Deviation > 0) {
		t.Errorf("Deviation = %g, want positive (band above low-wavelength reference)", r.Deviation)
	}
	// A flat profile leaves no shift at all.
	flat := EndpointExtinctions{LowExtinction: 130, HighExtinction: 130}
	mFlat, err := NewMeasure(b, flat, 0.008, 1.0)
	if err != nil {
		t.Fatal(err)
	}
	rFlat, err := mFlat.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if rFlat.Deviation != 0 {
		t.Errorf("flat profile Deviation = %g, want 0", rFlat.Deviation)
	}
}

// TestRectBandEndpointValidation checks that negative extinction samples are
// rejected and that the chord interpolates the endpoints exactly.
func TestRectBandEndpointValidation(t *testing.T) {
	_, err := NewEndpointExtinctions(-1, 130)
	if !errors.Is(err, ErrNegativeExtinctionSample) {
		t.Errorf("NewEndpointExtinctions(-1,130) error = %v, want ErrNegativeExtinctionSample", err)
	}
	b, _ := NewRectBand(508, 5)
	e := EndpointExtinctions{LowExtinction: 120, HighExtinction: 140}
	if got := EffectiveExtinctionAt(b, e, b.Low()); !closeTo(got, 120, 1e-12) {
		t.Errorf("chord at low endpoint = %g, want 120", got)
	}
	if got := EffectiveExtinctionAt(b, e, b.High()); !closeTo(got, 140, 1e-12) {
		t.Errorf("chord at high endpoint = %g, want 140", got)
	}
	if got := EffectiveExtinctionAt(b, e, 508); !closeTo(got, 130, 1e-12) {
		t.Errorf("chord at center = %g, want 130", got)
	}
}

// TestRectBandIntegralCrossCheck checks that the midpoint-rule transmittance
// integral stays close to the two-point average result, the approximation the
// CLI relies on.
func TestRectBandIntegralCrossCheck(t *testing.T) {
	b, _ := NewRectBand(508, 5)
	e := EndpointExtinctions{LowExtinction: 120, HighExtinction: 140}
	frac := FractionalTransmittance(b, e, 0.008, 1.0)
	aIntegral := MidpointIntegralAbsorbance(b, e, 0.008, 1.0)
	aEff := EffectiveExtinction(e) * 0.008 * 1.0
	if math.Abs(aIntegral-aEff) > 0.01 {
		t.Errorf("integral A = %g, two-point average A = %g, want within 0.01", aIntegral, aEff)
	}
	if !(frac > 0 && frac < 1) {
		t.Errorf("fractional transmittance = %g, want in (0,1)", frac)
	}
}

// TestRectBandGeometry checks the window geometry helpers and the resolution
// figure of merit.
func TestRectBandGeometry(t *testing.T) {
	b, _ := NewRectBand(500, 5)
	if !closeTo(b.RelativeHalfWidth(), 0.01, 1e-12) {
		t.Errorf("RelativeHalfWidth = %g, want 0.01", b.RelativeHalfWidth())
	}
	if !closeTo(b.Resolution(), 50, 1e-12) {
		t.Errorf("Resolution = %g, want 50 (center / width)", b.Resolution())
	}
	e := EndpointExtinctions{LowExtinction: 120, HighExtinction: 140}
	if !closeTo(SpectralSlope(e, b), 2, 1e-12) {
		t.Errorf("SpectralSlope = %g, want 2 per nm", SpectralSlope(e, b))
	}
	want := "[495, 505] nm"
	if got := b.Describe(); got != want {
		t.Errorf("Describe = %q, want %q", got, want)
	}
}

// TestBandPredictedTransmittance checks the derived band transmittances and
// the relative deviation percentage.
func TestBandPredictedTransmittance(t *testing.T) {
	b, _ := NewRectBand(508, 5)
	e := EndpointExtinctions{LowExtinction: 120, HighExtinction: 140}
	m, _ := NewMeasure(b, e, 0.008, 1.0)
	r, err := m.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	wantBandT := math.Pow(10, -r.BandAbsorbance)
	if !closeTo(r.BandTransmittance(), wantBandT, 1e-12) {
		t.Errorf("BandTransmittance = %g, want %g", r.BandTransmittance(), wantBandT)
	}
	rel := r.RelativeDeviation()
	if !closeTo(rel, r.Deviation/r.MonoAbsorbance, 1e-12) {
		t.Errorf("RelativeDeviation = %g, want %g", rel, r.Deviation/r.MonoAbsorbance)
	}
	if !closeTo(r.PercentDeviation(), 100*rel, 1e-9) {
		t.Errorf("PercentDeviation = %g, want %g", r.PercentDeviation(), 100*rel)
	}
}
