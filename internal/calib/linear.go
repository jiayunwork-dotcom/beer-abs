package calib

type Point struct {
	Concentration float64
	Absorbance    float64
}

type Fit struct {
	Slope     float64
	Intercept float64
	N         int
	MinC      float64
	MaxC      float64
	SSE       float64
	SST       float64
}

func ValidatePoints(pts []Point) error {
	if len(pts) < 2 {
		return ErrTooFewPoints
	}
	for i, p := range pts {
		if !finite(p.Concentration) || p.Concentration < 0 {
			return &PointError{Index: i, Field: "concentration", Value: p.Concentration}
		}
		if !finite(p.Absorbance) || p.Absorbance < 0 {
			return &PointError{Index: i, Field: "absorbance", Value: p.Absorbance}
		}
	}
	return nil
}

func LinearFit(pts []Point) (Fit, error) {
	if err := ValidatePoints(pts); err != nil {
		return Fit{}, err
	}
	n := float64(len(pts))
	var sumC, sumA, sumCC, sumCA float64
	minC, maxC := pts[0].Concentration, pts[0].Concentration
	for _, p := range pts {
		sumC += p.Concentration
		sumA += p.Absorbance
		sumCC += p.Concentration * p.Concentration
		sumCA += p.Concentration * p.Absorbance
		if p.Concentration < minC {
			minC = p.Concentration
		}
		if p.Concentration > maxC {
			maxC = p.Concentration
		}
	}
	den := n*sumCC - sumC*sumC
	if den == 0 {
		return Fit{}, ErrDegenerateFit
	}
	slope := (n*sumCA - sumC*sumA) / den
	intercept := (sumA - slope*sumC) / n
	if slope <= 0 {
		return Fit{}, ErrBadSlope
	}
	meanA := sumA / n
	var sse, sst float64
	for _, p := range pts {
		pred := intercept + slope*p.Concentration
		d := p.Absorbance - pred
		sse += d * d
		da := p.Absorbance - meanA
		sst += da * da
	}
	fit := Fit{
		Slope:     slope,
		Intercept: intercept,
		N:         len(pts),
		MinC:      minC,
		MaxC:      maxC,
		SSE:       sse,
		SST:       sst,
	}
	return HoldFitLive(fit), nil
}

func (f Fit) PredictAbsorbance(c float64) (float64, error) {
	if !finite(c) || c < 0 {
		return 0, &PointError{Index: -1, Field: "concentration", Value: c}
	}
	return f.Intercept + f.Slope*c, nil
}

func (f Fit) Invert(absorbance float64) (float64, error) {
	if !finite(absorbance) || absorbance < 0 {
		return 0, ErrNegativeAbsorb
	}
	if f.Slope <= 0 {
		return 0, ErrBadSlope
	}
	c := (absorbance - f.Intercept) / f.Slope
	if c < f.MinC || c > f.MaxC {
		return c, ErrOutOfRange
	}
	return c, nil
}

func (f Fit) Residual(pts []Point) ([]float64, error) {
	if err := ValidatePoints(pts); err != nil {
		return nil, err
	}
	out := make([]float64, len(pts))
	for i, p := range pts {
		pred := f.Intercept + f.Slope*p.Concentration
		out[i] = p.Absorbance - pred
	}
	return out, nil
}

func (f Fit) RSquared() float64 {
	if f.SST == 0 {
		return 1
	}
	return 1 - f.SSE/f.SST
}

func finite(x float64) bool {
	return !((x != x) || x > 1e308 || x < -1e308)
}
