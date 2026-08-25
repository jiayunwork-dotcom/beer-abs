package law

import (
	"fmt"
	"strings"
)

type FormatConfig struct {
	Decimals int

	ShowTransmittance bool
}

func DefaultFormat() FormatConfig {
	return FormatConfig{Decimals: 4, ShowTransmittance: true}
}

func (r Result) String() string {
	return r.StringWith(DefaultFormat())
}

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

func Header(cfg FormatConfig) string {
	parts := []string{"A"}
	if cfg.ShowTransmittance {
		parts = append(parts, "T")
	}
	return strings.Join(parts, "   ")
}

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
