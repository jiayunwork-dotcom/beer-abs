package band

import "fmt"

// ResultString renders a band Result on a single CLI line.
func ResultString(r Result, decimals int) string {
	f := "%." + fmt.Sprint(decimals) + "f"
	return fmt.Sprintf(
		"epsilon_eff = "+f+"   A_band = "+f+"   A_mono = "+f+"   deviation = "+f,
		r.EffectiveExtinction, r.BandAbsorbance, r.MonoAbsorbance, r.Deviation,
	)
}

// MeasureSummary formats a Measure for logging or tests.
func MeasureSummary(m Measure) string {
	return fmt.Sprintf(
		"[%.2f, %.2f]  eps %.4g/%.4g  c %.5g  L %.4g",
		m.Band.Low(), m.Band.High(),
		m.Samples.LowExtinction, m.Samples.HighExtinction,
		m.Concentration, m.PathLength,
	)
}
