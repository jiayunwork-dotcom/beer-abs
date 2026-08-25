package calib

import "fmt"

var (
	ErrTooFewPoints   = fmt.Errorf("calib: need at least 2 calibration points")
	ErrDegenerateFit  = fmt.Errorf("calib: concentration span is zero")
	ErrBadBlankStd    = fmt.Errorf("calib: blank standard deviation must be non-negative")
	ErrBadSlope       = fmt.Errorf("calib: calibration slope must be positive")
	ErrOutOfRange     = fmt.Errorf("calib: predicted concentration outside calibration span")
	ErrNegativeAbsorb = fmt.Errorf("calib: absorbance must be finite and non-negative")
)

type PointError struct {
	Index int
	Field string
	Value float64
}

func (e *PointError) Error() string {
	return fmt.Sprintf("calib: point %d field %s invalid value %g", e.Index, e.Field, e.Value)
}
