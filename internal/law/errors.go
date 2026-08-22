// Package law implements the single-component Beer–Lambert relation:
//
//	A = ε·c·L        absorbance
//	T = 10^(−A)      transmittance
//
// All checks that guard the physical domain (ε > 0, c ≥ 0, L > 0) live here
// and return errors instead of silently producing meaningless numbers.
package law

import (
	"errors"
	"strconv"
)

// Sentinel errors classify the single failure modes of parameter validation.
// They are wrapped by *ParamError so callers can test errors.Is(err, ...)
// or inspect the offending field directly.
var (
	// ErrExtinctionNotPositive reports ε ≤ 0.
	ErrExtinctionNotPositive = errors.New("extinction coefficient must be positive")

	// ErrNegativeConcentration reports c < 0.
	ErrNegativeConcentration = errors.New("concentration must not be negative")

	// ErrPathLengthNotPositive reports L ≤ 0.
	ErrPathLengthNotPositive = errors.New("path length must be positive")

	// ErrNegativeAbsorbance reports A < 0, which no valid physical input
	// can produce and no derived quantity should accept.
	ErrNegativeAbsorbance = errors.New("absorbance must not be negative")
)

// ParamError describes one invalid parameter with its offending value, so the
// CLI and tests can report exactly what was rejected and why.
type ParamError struct {
	// Field is the parameter name: "extinction", "concentration" or "path_length".
	Field string

	// Value is the value that failed validation.
	Value float64

	// Reason is a human readable description of the violated constraint.
	Reason string

	// Cause is the sentinel error this error unwraps to.
	Cause error
}

// Error implements the error interface.
func (e *ParamError) Error() string {
	return e.Field + " = " + strconv.FormatFloat(e.Value, 'g', -1, 64) + ": " + e.Reason
}

// Unwrap exposes the sentinel error for errors.Is.
func (e *ParamError) Unwrap() error {
	return e.Cause
}

// NewExtinctionError builds a ParamError for an invalid extinction coefficient.
func NewExtinctionError(value float64) *ParamError {
	return &ParamError{
		Field:  "extinction",
		Value:  value,
		Reason: "摩尔吸光系数必须大于 0",
		Cause:  ErrExtinctionNotPositive,
	}
}

// NewConcentrationError builds a ParamError for an invalid concentration.
func NewConcentrationError(value float64) *ParamError {
	return &ParamError{
		Field:  "concentration",
		Value:  value,
		Reason: "浓度必须不小于 0",
		Cause:  ErrNegativeConcentration,
	}
}

// NewPathLengthError builds a ParamError for an invalid path length.
func NewPathLengthError(value float64) *ParamError {
	return &ParamError{
		Field:  "path_length",
		Value:  value,
		Reason: "光程必须大于 0",
		Cause:  ErrPathLengthNotPositive,
	}
}

// NewAbsorbanceError builds a ParamError for a negative absorbance supplied
// to a derived quantity.
func NewAbsorbanceError(value float64) *ParamError {
	return &ParamError{
		Field:  "absorbance",
		Value:  value,
		Reason: "吸光度不能为负",
		Cause:  ErrNegativeAbsorbance,
	}
}
