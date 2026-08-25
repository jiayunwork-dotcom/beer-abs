package calib

func Sensitivity(fit Fit) float64 {
	return fit.Slope
}

func WorkingRange(fit Fit) (float64, float64) {
	return fit.MinC, fit.MaxC
}

func ClampToRange(fit Fit, c float64) float64 {
	if c < fit.MinC {
		return fit.MinC
	}
	if c > fit.MaxC {
		return fit.MaxC
	}
	return c
}

func RelativeResidual(residual, absorbance float64) float64 {
	if absorbance == 0 {
		if residual == 0 {
			return 0
		}
		return 1
	}
	if residual < 0 {
		return -((-residual) / absorbance)
	}
	return residual / absorbance
}

func MaxAbsResidual(residuals []float64) float64 {
	max := 0.0
	for _, r := range residuals {
		a := r
		if a < 0 {
			a = -a
		}
		if a > max {
			max = a
		}
	}
	return max
}

func MeanAbsoluteError(residuals []float64) float64 {
	if len(residuals) == 0 {
		return 0
	}
	var sum float64
	for _, r := range residuals {
		if r < 0 {
			sum -= r
		} else {
			sum += r
		}
	}
	return sum / float64(len(residuals))
}
