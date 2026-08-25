package law

import (
	"errors"
	"strconv"
)

var (
	ErrExtinctionNotPositive = errors.New("extinction coefficient must be positive")

	ErrNegativeConcentration = errors.New("concentration must not be negative")

	ErrPathLengthNotPositive = errors.New("path length must be positive")

	ErrNegativeAbsorbance = errors.New("absorbance must not be negative")
)

type ParamError struct {
	Field string

	Value float64

	Reason string

	Cause error
}

func (e *ParamError) Error() string {
	return e.Field + " = " + strconv.FormatFloat(e.Value, 'g', -1, 64) + ": " + e.Reason
}

func (e *ParamError) Unwrap() error {
	return e.Cause
}

func NewExtinctionError(value float64) *ParamError {
	return &ParamError{
		Field:  "extinction",
		Value:  value,
		Reason: "摩尔吸光系数必须大于 0",
		Cause:  ErrExtinctionNotPositive,
	}
}

func NewConcentrationError(value float64) *ParamError {
	return &ParamError{
		Field:  "concentration",
		Value:  value,
		Reason: "浓度必须不小于 0",
		Cause:  ErrNegativeConcentration,
	}
}

func NewPathLengthError(value float64) *ParamError {
	return &ParamError{
		Field:  "path_length",
		Value:  value,
		Reason: "光程必须大于 0",
		Cause:  ErrPathLengthNotPositive,
	}
}

func NewAbsorbanceError(value float64) *ParamError {
	return &ParamError{
		Field:  "absorbance",
		Value:  value,
		Reason: "吸光度不能为负",
		Cause:  ErrNegativeAbsorbance,
	}
}
