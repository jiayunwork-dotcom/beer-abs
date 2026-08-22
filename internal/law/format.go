package law

import (
	"fmt"
	"strings"
)

// FormatConfig controls how Evaluate results are rendered by the CLI.
type FormatConfig struct {
	// Decimals is the number of decimal places for floating point values.
	Decimals int

	// ShowTransmittance toggles the T = 10^(−A) line.
	ShowTransmittance bool
}

// DefaultFormat returns the standard rendering used by the CLI.
func DefaultFormat() FormatConfig {
	return FormatConfig{Decimals: 4, ShowTransmittance: true}
}

// String renders a Result on a single line:
//
//	A = 1.0000   T = 0.1000
func (r Result) String() string {
	return r.StringWith(DefaultFormat())
}

// StringWith renders a Result with a custom FormatConfig.
func (r Result) StringWith(cfg FormatConfig) string {
	f := "%." + itoaDecimals(cfg.Decimals) + "f"
	parts := []string{
		"A = " + fmt.Sprintf(f, r.Absorbance),
	}
	if cfg.ShowTransmittance {
		parts = append(parts, "T = "+fmt.Sprintf(f, r.Transmittance))
	}
	return strings.Join(parts, "   ")
}

// Header returns the field labels aligned with String's output, useful for
// column-style CLI output.
func Header(cfg FormatConfig) string {
	parts := []string{"A"}
	if cfg.ShowTransmittance {
		parts = append(parts, "T")
	}
	return strings.Join(parts, "   ")
}

// itoaDecimals converts the decimal count into a format width suffix.
func itoaDecimals(n int) string {
	switch n {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	default:
		return "4"
	}
}
