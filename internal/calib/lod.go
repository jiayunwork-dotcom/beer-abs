package calib

type DetectionLimits struct {
	BlankMean float64
	BlankStd  float64
	LOD       float64
	LOQ       float64
	Slope     float64
}

func LimitsFromBlank(blankAbsorbances []float64, slope float64) (DetectionLimits, error) {
	if slope <= 0 {
		return DetectionLimits{}, ErrBadSlope
	}
	if len(blankAbsorbances) == 0 {
		return DetectionLimits{}, ErrTooFewPoints
	}
	var sum float64
	for i, a := range blankAbsorbances {
		if !finite(a) || a < 0 {
			return DetectionLimits{}, &PointError{Index: i, Field: "blank", Value: a}
		}
		sum += a
	}
	mean := sum / float64(len(blankAbsorbances))
	var ss float64
	for _, a := range blankAbsorbances {
		d := a - mean
		ss += d * d
	}
	std := 0.0
	if len(blankAbsorbances) > 1 {
		std = sqrt(ss / float64(len(blankAbsorbances)-1))
	}
	if std < 0 {
		return DetectionLimits{}, ErrBadBlankStd
	}
	lod := 3 * std / slope
	loq := 10 * std / slope
	return DetectionLimits{
		BlankMean: mean,
		BlankStd:  std,
		LOD:       lod,
		LOQ:       loq,
		Slope:     slope,
	}, nil
}

func (d DetectionLimits) AboveLOD(concentration float64) bool {
	return concentration >= d.LOD
}

func (d DetectionLimits) Quantifiable(concentration float64) bool {
	return concentration >= d.LOQ
}

func RoundTripConcentrations(pts []Point, absorbances []float64) ([]float64, error) {
	fit, err := LinearFit(pts)
	if err != nil {
		return nil, err
	}
	out := make([]float64, len(absorbances))
	for i, a := range absorbances {
		c, err := fit.Invert(a)
		if err != nil && err != ErrOutOfRange {
			return nil, err
		}
		out[i] = c
	}
	return out, nil
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 16; i++ {
		z = 0.5 * (z + x/z)
	}
	return z
}
