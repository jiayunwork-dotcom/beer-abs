package band

import (
	"errors"
	"strconv"
)

// ErrNegativeExtinctionSample reports a sampled extinction below zero.
var ErrNegativeExtinctionSample = errors.New("extinction sample must not be negative")

// formatFloat renders a float compactly for error messages.
func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// ExtinctionSampleError describes a single invalid extinction sample.
type ExtinctionSampleError struct {
	// Field names the sampled quantity.
	Field string

	// Value is the invalid sample.
	Value float64

	// Reason is a human readable description.
	Reason string
}

// Error implements the error interface.
func (e *ExtinctionSampleError) Error() string {
	return e.Field + " = " + formatFloat(e.Value) + ": " + e.Reason
}

// Unwrap exposes the sentinel error for errors.Is.
func (e *ExtinctionSampleError) Unwrap() error {
	return ErrNegativeExtinctionSample
}

// NewExtinctionSampleError builds an ExtinctionSampleError.
func NewExtinctionSampleError(field string, value float64) *ExtinctionSampleError {
	return &ExtinctionSampleError{
		Field:  field,
		Value:  value,
		Reason: "带宽端点处的摩尔吸光系数不能为负",
	}
}
