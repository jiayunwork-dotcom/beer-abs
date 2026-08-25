package mixture

import (
	"errors"
	"strconv"
)

var (
	ErrNoComponents = errors.New("mixture needs at least one component")

	ErrStrayOutOfRange = errors.New("stray fraction must be in [0, 1)")
)

type MixtureError struct {
	Field string

	Value float64

	Reason string

	Cause error
}

func (e *MixtureError) Error() string {
	return "mixture." + e.Field + " = " + strconv.FormatFloat(e.Value, 'g', -1, 64) + ": " + e.Reason
}

func (e *MixtureError) Unwrap() error {
	return e.Cause
}

func NewEmptyMixtureError() *MixtureError {
	return &MixtureError{
		Field:  "components",
		Reason: "混合物至少需要一个组分",
		Cause:  ErrNoComponents,
	}
}

func NewStrayError(value float64) *MixtureError {
	return &MixtureError{
		Field:  "stray_fraction",
		Value:  value,
		Reason: "杂散光分数必须满足 0 ≤ s < 1",
		Cause:  ErrStrayOutOfRange,
	}
}
