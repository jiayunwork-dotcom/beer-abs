// Package mixture assembles several absorbing components into one sample and
// applies the two instrument-level corrections the CLI understands:
//
//   - additivity: with non-interacting species the total absorbance is the
//     sum of the individual absorbances, A_tot = Σ ε_i·c_i·L, all sharing the
//     same path length L;
//   - stray light: a fraction s of the detector signal comes from light that
//     never passed through the sample, so the observed transmittance is
//
//	T_obs = (T + s) / (1 + s)      A_obs = −log10(T_obs)
//
// When s = 0 the observed values collapse back to the ideal ones, and for
// large absorbances stray light pushes A_obs below A_ideal.
package mixture

import (
	"errors"
	"strconv"
)

// Sentinel errors for mixture-level validation.
var (
	// ErrNoComponents reports an empty component list.
	ErrNoComponents = errors.New("mixture needs at least one component")

	// ErrStrayOutOfRange reports s outside [0, 1).
	ErrStrayOutOfRange = errors.New("stray fraction must be in [0, 1)")
)

// MixtureError describes an invalid mixture-level parameter.
type MixtureError struct {
	// Field names the offending parameter.
	Field string

	// Value is the invalid value.
	Value float64

	// Reason is a human readable description of the violated constraint.
	Reason string

	// Cause is the sentinel error this error unwraps to.
	Cause error
}

// Error implements the error interface.
func (e *MixtureError) Error() string {
	return "mixture." + e.Field + " = " + strconv.FormatFloat(e.Value, 'g', -1, 64) + ": " + e.Reason
}

// Unwrap exposes the sentinel error for errors.Is.
func (e *MixtureError) Unwrap() error {
	return e.Cause
}

// NewEmptyMixtureError builds an error for an empty component list.
func NewEmptyMixtureError() *MixtureError {
	return &MixtureError{
		Field:  "components",
		Reason: "混合物至少需要一个组分",
		Cause:  ErrNoComponents,
	}
}

// NewStrayError builds an error for an out-of-range stray fraction.
func NewStrayError(value float64) *MixtureError {
	return &MixtureError{
		Field:  "stray_fraction",
		Value:  value,
		Reason: "杂散光分数必须满足 0 ≤ s < 1",
		Cause:  ErrStrayOutOfRange,
	}
}
