// Package band models a finite rectangular spectral passband. A monochromator
// does not deliver a single wavelength; it passes a window of wavelengths. For
// a rectangular window of half width W around a center λ0, a smooth extinction
// profile ε(λ) is approximated here by its linear chord across the two window
// endpoints, and the effective extinction is the chord's mean over the window.
//
//	ε_eff = (ε(λ0−W) + ε(λ0+W)) / 2
//
// The band-averaged absorbance A_band = ε_eff·c·L therefore differs from the
// single-wavelength A(λ0) = ε(λ0)·c·L whenever the extinction varies across
// the window. This is the "thin path, finite bandwidth" deviation the CLI
// reports alongside the ideal value.
package band

import "errors"

// Sentinel errors for the rectangular-band parameter checks.
var (
	// ErrCenterNotPositive reports λ0 ≤ 0.
	ErrCenterNotPositive = errors.New("band center wavelength must be positive")

	// ErrHalfWidthNotPositive reports W ≤ 0.
	ErrHalfWidthNotPositive = errors.New("band half width must be positive")

	// ErrEndpointOrder reports a low endpoint at or above the high endpoint.
	ErrEndpointOrder = errors.New("low endpoint must be below the high endpoint")
)

// BandError describes a single invalid band parameter.
type BandError struct {
	// Field names the offending parameter.
	Field string

	// Value is the value that failed validation.
	Value float64

	// Reason is a human readable description of the violated constraint.
	Reason string

	// Cause is the sentinel error this error unwraps to.
	Cause error
}

// Error implements the error interface.
func (e *BandError) Error() string {
	return "band." + e.Field + ": " + e.Reason
}

// Unwrap exposes the sentinel error for errors.Is.
func (e *BandError) Unwrap() error {
	return e.Cause
}

// NewCenterError builds a BandError for a non-positive center wavelength.
func NewCenterError(value float64) *BandError {
	return &BandError{
		Field:  "center",
		Value:  value,
		Reason: "通带中心波长必须大于 0",
		Cause:  ErrCenterNotPositive,
	}
}

// NewHalfWidthError builds a BandError for a non-positive half width.
func NewHalfWidthError(value float64) *BandError {
	return &BandError{
		Field:  "half_width",
		Value:  value,
		Reason: "通带半宽必须大于 0",
		Cause:  ErrHalfWidthNotPositive,
	}
}

// NewEndpointOrderError builds a BandError for an inverted endpoint pair.
func NewEndpointOrderError(low, high float64) *BandError {
	return &BandError{
		Field:  "endpoints",
		Value:  low,
		Reason: "低端波长必须小于高端波长",
		Cause:  ErrEndpointOrder,
	}
}
