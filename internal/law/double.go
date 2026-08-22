package law

import "math"

// DoubleConcentration returns a copy of p with the concentration doubled.
// Because A = ε·c·L, doubling c doubles A, and because T = 10^(−A), the
// transmittance of the doubled sample is the square of the original T.
func DoubleConcentration(p Params) Params {
	p.Concentration *= 2
	return p
}

// DoublePathLength returns a copy of p with the path length doubled.
// Doubling L doubles A in exactly the same way as doubling the concentration.
func DoublePathLength(p Params) Params {
	p.PathLength *= 2
	return p
}

// DoubleAbsorbanceConcentration doubles a concentration and returns the
// resulting absorbance plus the squared transmittance in one call.
func DoubleAbsorbanceConcentration(p Params) (doubleA float64, squaredT float64, err error) {
	p2 := DoubleConcentration(p)
	a1, err := p.Absorbance()
	if err != nil {
		return 0, 0, err
	}
	a2, err := p2.Absorbance()
	if err != nil {
		return 0, 0, err
	}
	return a2, TransmittanceSquared(a1), nil
}

// DoubleAbsorbancePath doubles the path length and returns the new absorbance
// and the squared transmittance.
func DoubleAbsorbancePath(p Params) (doubleA float64, squaredT float64, err error) {
	p2 := DoublePathLength(p)
	a1, err := p.Absorbance()
	if err != nil {
		return 0, 0, err
	}
	a2, err := p2.Absorbance()
	if err != nil {
		return 0, 0, err
	}
	return a2, TransmittanceSquared(a1), nil
}

// DoublingFactor is the ratio A_doubled / A_original; for a linear law it is
// exactly 2 within floating point tolerance.
func DoublingFactor(aOriginal, aDoubled float64) float64 {
	if aOriginal == 0 {
		return 0
	}
	return aDoubled / aOriginal
}

// TransmittanceFactor is the ratio T_original / T_doubled. When the
// absorbance doubles, this factor equals 1/T_original, i.e. the transmittance
// itself is squared.
func TransmittanceFactor(tOriginal, tDoubled float64) float64 {
	if tDoubled == 0 {
		return mathInf()
	}
	return tOriginal / tDoubled
}

// mathInf returns positive infinity, the natural value of a division whose
// denominator is exactly zero.
func mathInf() float64 {
	return math.Inf(1)
}
