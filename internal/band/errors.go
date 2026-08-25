package band

import "errors"

var (
	ErrCenterNotPositive = errors.New("band center wavelength must be positive")

	ErrHalfWidthNotPositive = errors.New("band half width must be positive")

	ErrEndpointOrder = errors.New("low endpoint must be below the high endpoint")
)

type BandError struct {
	Field string

	Value float64

	Reason string

	Cause error
}

func (e *BandError) Error() string {
	return "band." + e.Field + ": " + e.Reason
}

func (e *BandError) Unwrap() error {
	return e.Cause
}

func NewCenterError(value float64) *BandError {
	return &BandError{
		Field:  "center",
		Value:  value,
		Reason: "通带中心波长必须大于 0",
		Cause:  ErrCenterNotPositive,
	}
}

func NewHalfWidthError(value float64) *BandError {
	return &BandError{
		Field:  "half_width",
		Value:  value,
		Reason: "通带半宽必须大于 0",
		Cause:  ErrHalfWidthNotPositive,
	}
}

func NewEndpointOrderError(low, high float64) *BandError {
	return &BandError{
		Field:  "endpoints",
		Value:  low,
		Reason: "低端波长必须小于高端波长",
		Cause:  ErrEndpointOrder,
	}
}
