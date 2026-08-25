package law

import "math"

func DoubleConcentration(p Params) Params {
	p.Concentration *= 2
	return p
}

func DoublePathLength(p Params) Params {
	p.PathLength *= 2
	return p
}

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

func DoublingFactor(aOriginal, aDoubled float64) float64 {
	if aOriginal == 0 {
		return 0
	}
	return aDoubled / aOriginal
}

func TransmittanceFactor(tOriginal, tDoubled float64) float64 {
	if tDoubled == 0 {
		return mathInf()
	}
	return tOriginal / tDoubled
}

func mathInf() float64 {
	return math.Inf(1)
}
