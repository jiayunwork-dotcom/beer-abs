package mixture

import (
	"errors"
	"math"
	"testing"

	"beer-abs/internal/band"
	"beer-abs/internal/law"
)

func near(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// TestMixtureAbsorbanceIsAdditive checks the core additivity invariant: the
// mixture absorbance equals the sum of the individual component absorbances
// for a shared path length, without stray light.
func TestMixtureAbsorbanceIsAdditive(t *testing.T) {
	cu := Component{Label: "Cu", Extinction: 125, Concentration: 0.008}
	ni := Component{Label: "Ni", Extinction: 25, Concentration: 0.012}
	const L = 1.0
	aCu, err := cu.Absorbance(L)
	if err != nil {
		t.Fatal(err)
	}
	aNi, err := ni.Absorbance(L)
	if err != nil {
		t.Fatal(err)
	}
	want := aCu + aNi

	m, err := New([]Component{cu, ni}, L, 0)
	if err != nil {
		t.Fatal(err)
	}
	got, err := m.TotalAbsorbance()
	if err != nil {
		t.Fatal(err)
	}
	if !near(got, want, 1e-12) {
		t.Errorf("TotalAbsorbance = %g, want %g (sum of components)", got, want)
	}
	an, err := m.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if !an.IsAdditive() {
		t.Errorf("mixture not reported additive: sum %g vs total %g", want, an.Absorbance)
	}
	parts, err := m.ComponentAbsorbances()
	if err != nil {
		t.Fatal(err)
	}
	sum := 0.0
	for _, p := range parts {
		sum += p
	}
	if !near(sum, got, 1e-12) {
		t.Errorf("ComponentAbsorbances sum = %g, want %g", sum, got)
	}
}

// TestStrayFractionRangeError checks the stray-light domain: s outside
// [0, 1) must be rejected, while every interior value is accepted.
func TestStrayFractionRangeError(t *testing.T) {
	cu := Component{Label: "Cu", Extinction: 125, Concentration: 0.008}
	for _, s := range []float64{-0.1, 1.0, 2.0} {
		_, err := New([]Component{cu}, 1.0, s)
		if !errors.Is(err, ErrStrayOutOfRange) {
			t.Errorf("New with s=%g error = %v, want ErrStrayOutOfRange", s, err)
		}
		if err := ValidateStrayFraction(s); !errors.Is(err, ErrStrayOutOfRange) {
			t.Errorf("ValidateStrayFraction(%g) error = %v, want ErrStrayOutOfRange", s, err)
		}
	}
	for _, s := range []float64{0, 0.001, 0.5, 0.999} {
		if err := ValidateStrayFraction(s); err != nil {
			t.Errorf("ValidateStrayFraction(%g) = %v, want nil", s, err)
		}
	}
}

// TestStrayZeroIsIdeal checks that s = 0 collapses the observed values back
// to the ideal ones exactly.
func TestStrayZeroIsIdeal(t *testing.T) {
	cu := Component{Label: "Cu", Extinction: 125, Concentration: 0.008}
	ni := Component{Label: "Ni", Extinction: 25, Concentration: 0.012}
	m, err := New([]Component{cu, ni}, 1.0, 0)
	if err != nil {
		t.Fatal(err)
	}
	an, err := m.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if !an.Ideal {
		t.Errorf("Ideal = false, want true at s=0")
	}
	if an.ObservedAbsorbance != an.Absorbance {
		t.Errorf("A_obs = %g, want %g (ideal)", an.ObservedAbsorbance, an.Absorbance)
	}
	if an.ObservedTransmittance != an.Transmittance {
		t.Errorf("T_obs = %g, want %g (ideal)", an.ObservedTransmittance, an.Transmittance)
	}
	if an.Deviation != 0 {
		t.Errorf("Deviation = %g, want 0", an.Deviation)
	}
}

// TestObservedTransmittanceWithinUnitInterval checks the stray-light formula
// as a whole: for every valid s and absorbance the observed transmittance
// stays inside (0, 1], so A_obs is real and the deviation direction is fixed.
func TestObservedTransmittanceWithinUnitInterval(t *testing.T) {
	for _, s := range []float64{0.01, 0.05, 0.2, 0.5} {
		for _, c := range []float64{1e-6, 0.001, 0.008, 0.04, 0.2, 1.0} {
			comp := Component{Label: "X", Extinction: 125, Concentration: c}
			m, err := New([]Component{comp}, 1.0, s)
			if err != nil {
				t.Fatal(err)
			}
			tObs, err := m.ObservedTransmittance()
			if err != nil {
				t.Fatal(err)
			}
			if !(tObs > 0 && tObs <= 1) {
				t.Errorf("s=%g c=%g: T_obs = %g, want in (0,1]", s, c, tObs)
			}
			aObs, err := m.ObservedAbsorbance()
			if err != nil {
				t.Fatal(err)
			}
			if math.IsInf(aObs, 0) || math.IsNaN(aObs) {
				t.Errorf("s=%g c=%g: A_obs = %g, want finite", s, c, aObs)
			}
			// T_obs must be above the fully-absorbing floor s/(1+s).
			floor, _ := StrayFloor(s)
			if tObs < floor {
				t.Errorf("s=%g c=%g: T_obs = %g below floor %g", s, c, tObs, floor)
			}
		}
	}
}

// TestHighAbsorbanceNegativeDeviation checks the stray-light signature: at
// high A the observed absorbance drops below the ideal one, while at low A it
// stays close to ideal.
func TestHighAbsorbanceNegativeDeviation(t *testing.T) {
	const s = 0.05
	high := Component{Label: "X", Extinction: 125, Concentration: 0.04}
	mHigh, err := New([]Component{high}, 1.0, s)
	if err != nil {
		t.Fatal(err)
	}
	anHigh, err := mHigh.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if !anHigh.NegativeDeviation() {
		t.Errorf("high A: deviation = %g, want negative (A_obs < A_ideal)", anHigh.Deviation)
	}
	if !(anHigh.ObservedAbsorbance < anHigh.Absorbance) {
		t.Errorf("high A: A_obs = %g, want < A_ideal = %g", anHigh.ObservedAbsorbance, anHigh.Absorbance)
	}

	low := Component{Label: "Y", Extinction: 125, Concentration: 1e-6}
	mLow, err := New([]Component{low}, 1.0, s)
	if err != nil {
		t.Fatal(err)
	}
	anLow, err := mLow.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(anLow.Deviation) > 0.01 {
		t.Errorf("low A: deviation = %g, want close to 0", anLow.Deviation)
	}
}

// TestZeroConcentrationIndependentOfStray checks that a transparent sample
// (c = 0) reports T = 1 and A = 0 no matter what the stray fraction is.
func TestZeroConcentrationIndependentOfStray(t *testing.T) {
	for _, s := range []float64{0, 0.05, 0.9} {
		comp := Component{Label: "Blank", Extinction: 125, Concentration: 0}
		m, err := New([]Component{comp}, 1.0, s)
		if err != nil {
			t.Fatal(err)
		}
		an, err := m.Analyze()
		if err != nil {
			t.Fatal(err)
		}
		if an.Absorbance != 0 {
			t.Errorf("s=%g: A = %g, want 0", s, an.Absorbance)
		}
		if an.Transmittance != 1 {
			t.Errorf("s=%g: T = %g, want 1", s, an.Transmittance)
		}
		if an.ObservedTransmittance != 1 {
			t.Errorf("s=%g: T_obs = %g, want 1", s, an.ObservedTransmittance)
		}
		if an.ObservedAbsorbance != 0 {
			t.Errorf("s=%g: A_obs = %g, want 0", s, an.ObservedAbsorbance)
		}
	}
}

// TestMixtureValidationRejectsBadComponents checks that per-component errors
// propagate out of the mixture constructor with their original sentinels.
func TestMixtureValidationRejectsBadComponents(t *testing.T) {
	bad := Component{Label: "Bad", Extinction: 125, Concentration: -0.001}
	_, err := New([]Component{bad}, 1.0, 0)
	if !errors.Is(err, law.ErrNegativeConcentration) {
		t.Errorf("error = %v, want law.ErrNegativeConcentration", err)
	}
	_, err = New(nil, 1.0, 0)
	if !errors.Is(err, ErrNoComponents) {
		t.Errorf("error = %v, want ErrNoComponents", err)
	}
}

// TestApplyBandAveragesExtinctions checks that applying a band replaces every
// component extinction with the endpoint average.
func TestApplyBandAveragesExtinctions(t *testing.T) {
	cu := Component{Label: "Cu", Extinction: 125, Concentration: 0.008}
	ni := Component{Label: "Ni", Extinction: 25, Concentration: 0.012}
	m, err := New([]Component{cu, ni}, 1.0, 0)
	if err != nil {
		t.Fatal(err)
	}
	b, err := band.NewRectBand(508, 5)
	if err != nil {
		t.Fatal(err)
	}
	m2, err := m.ApplyBandFromSamples(b, []float64{120, 20}, []float64{140, 40})
	if err != nil {
		t.Fatal(err)
	}
	if !near(m2.Components[0].Extinction, 130, 1e-12) {
		t.Errorf("Cu eps = %g, want 130 (average of 120/140)", m2.Components[0].Extinction)
	}
	if !near(m2.Components[1].Extinction, 30, 1e-12) {
		t.Errorf("Ni eps = %g, want 30 (average of 20/40)", m2.Components[1].Extinction)
	}
}

// TestDilutionSplitsAbsorbance checks that a two-fold dilution halves every
// component concentration and therefore halves the total absorbance.
func TestDilutionSplitsAbsorbance(t *testing.T) {
	cu := Component{Label: "Cu", Extinction: 125, Concentration: 0.008}
	m, err := New([]Component{cu}, 1.0, 0)
	if err != nil {
		t.Fatal(err)
	}
	aBefore, err := m.TotalAbsorbance()
	if err != nil {
		t.Fatal(err)
	}
	diluted, err := m.Dilute(2)
	if err != nil {
		t.Fatal(err)
	}
	aAfter, err := diluted.TotalAbsorbance()
	if err != nil {
		t.Fatal(err)
	}
	if !near(aAfter, aBefore/2, 1e-12) {
		t.Errorf("diluted A = %g, want %g (half)", aAfter, aBefore/2)
	}
	if !near(diluted.Components[0].Concentration, 0.004, 1e-12) {
		t.Errorf("diluted c = %g, want 0.004", diluted.Components[0].Concentration)
	}
	if _, err := m.Dilute(0); err == nil {
		t.Errorf("Dilute(0) error = nil, want DilutionError")
	}
}

// TestBaselineSubtraction checks the blank-corrected net absorbance.
func TestBaselineSubtraction(t *testing.T) {
	net, err := SubtractBlank(1.3, 0.2)
	if err != nil {
		t.Fatal(err)
	}
	if !near(net, 1.1, 1e-12) {
		t.Errorf("net = %g, want 1.1", net)
	}
	if _, err := SubtractBlank(-1, 0); err == nil {
		t.Errorf("SubtractBlank(-1,0) error = nil, want BaselineError")
	}
	sample := Component{Label: "S", Extinction: 125, Concentration: 0.008}
	blank := Component{Label: "B", Extinction: 50, Concentration: 0.004}
	mS, _ := New([]Component{sample}, 1.0, 0)
	mB, _ := New([]Component{blank}, 1.0, 0)
	res, err := mS.CorrectBlank(mB)
	if err != nil {
		t.Fatal(err)
	}
	if !near(res.NetAbsorbance, 0.8, 1e-12) {
		t.Errorf("CorrectBlank net = %g, want 0.8 (1.0 - 0.2)", res.NetAbsorbance)
	}
}
