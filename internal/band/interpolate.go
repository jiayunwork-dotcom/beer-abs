package band

import (
	"errors"
	"strconv"
)

var ErrNegativeExtinctionSample = errors.New("extinction sample must not be negative")

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

type ExtinctionSampleError struct {
	Field string

	Value float64

	Reason string
}

func (e *ExtinctionSampleError) Error() string {
	return e.Field + " = " + formatFloat(e.Value) + ": " + e.Reason
}

func (e *ExtinctionSampleError) Unwrap() error {
	return ErrNegativeExtinctionSample
}

func NewExtinctionSampleError(field string, value float64) *ExtinctionSampleError {
	return &ExtinctionSampleError{
		Field:  field,
		Value:  value,
		Reason: "带宽端点处的摩尔吸光系数不能为负",
	}
}
