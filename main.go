package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"beer-abs/internal/band"
	"beer-abs/internal/mixture"
)

const usage = `beer-abs: solution absorbance calculator (Beer–Lambert).

Reads a JSON sample description, computes the ideal absorbance
A_tot = Σ ε_i·c_i·L and transmittance T = 10^(−A_tot), then applies the
optional stray-light correction T_obs = (T + s)/(1 + s) and an optional
finite-bandwidth extinction average.

usage:
  beer-abs absorbance <sample.json>
  beer-abs help

sample.json fields:
  components      array of {label, extinction, concentration} or, when a band
                  is given, {label, extinction_low, extinction_high, concentration}
  path_length     optical path L in cm (must be > 0)
  stray_fraction  stray-light fraction s in [0, 1) (default 0)
  band            optional {center, half_width} in nm; when present the two
                  endpoint extinctions are averaged over the rectangular window

Illegal inputs (ε <= 0, c < 0, L <= 0, s outside [0, 1), missing components,
unknown fields) are reported on stderr and exit non-zero.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "absorbance":
		if err := runAbsorbance(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "beer-abs: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "beer-abs: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

// componentSpec is the JSON shape of one absorbing species.
type componentSpec struct {
	Label          string  `json:"label"`
	Extinction     float64 `json:"extinction"`
	ExtinctionLow  float64 `json:"extinction_low"`
	ExtinctionHigh float64 `json:"extinction_high"`
	Concentration  float64 `json:"concentration"`
}

// bandSpec is the JSON shape of the optional rectangular passband.
type bandSpec struct {
	Center    float64 `json:"center"`
	HalfWidth float64 `json:"half_width"`
}

// sampleSpec is the whole JSON input file.
type sampleSpec struct {
	Components    []componentSpec `json:"components"`
	PathLength    float64         `json:"path_length"`
	StrayFraction float64         `json:"stray_fraction"`
	Band          *bandSpec       `json:"band"`
}

func runAbsorbance(args []string) error {
	if len(args) != 1 {
		return errors.New("absorbance needs exactly one sample JSON file")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("read sample: %w", err)
	}
	var sp sampleSpec
	if err := decodeJSON(data, &sp); err != nil {
		return fmt.Errorf("parse sample: %w", err)
	}

	var rect *band.RectBand
	if sp.Band != nil {
		b, err := band.NewRectBand(sp.Band.Center, sp.Band.HalfWidth)
		if err != nil {
			return err
		}
		rect = &b
	}

	m, err := buildMixture(sp, rect)
	if err != nil {
		return err
	}
	analysis, err := m.Analyze()
	if err != nil {
		return err
	}
	printAnalysis(sp, analysis, rect)
	return nil
}

// buildMixture translates the JSON spec into a mixture.Mixture, averaging the
// two endpoint extinctions when a passband is present.
func buildMixture(sp sampleSpec, rect *band.RectBand) (mixture.Mixture, error) {
	comps := make([]mixture.Component, len(sp.Components))
	for i, cs := range sp.Components {
		ext := cs.Extinction
		if rect != nil {
			ext = band.BandAverageExtinction(cs.ExtinctionLow, cs.ExtinctionHigh)
		}
		comps[i] = mixture.Component{
			Label:         cs.Label,
			Extinction:    ext,
			Concentration: cs.Concentration,
		}
	}
	return mixture.New(comps, sp.PathLength, sp.StrayFraction)
}

// decodeJSON rejects unknown fields so a typo in the input cannot silently
// change the measurement.
func decodeJSON(data []byte, out *sampleSpec) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(out)
}

func printAnalysis(sp sampleSpec, a mixture.Analysis, rect *band.RectBand) {
	fmt.Printf("sample: %d component(s), path L = %g cm, stray s = %g\n",
		len(a.Components), a.PathLength, a.StrayFraction)
	for i, c := range a.Components {
		ai := c.Extinction * c.Concentration * a.PathLength
		fmt.Printf("  [%d] %s: eps = %g, c = %g, A_i = %g\n",
			i, c.Label, c.Extinction, c.Concentration, ai)
	}
	if rect != nil {
		for i, cs := range sp.Components {
			epsEff := band.BandAverageExtinction(cs.ExtinctionLow, cs.ExtinctionHigh)
			aMono := cs.ExtinctionLow * cs.Concentration * a.PathLength
			fmt.Printf("  band[%d] %s: window [%g, %g] nm, A_band = %g, A_mono = %g\n",
				i, cs.Label, rect.Low(), rect.High(),
				epsEff*cs.Concentration*a.PathLength, aMono)
		}
	}
	fmt.Printf("A = %g   T = %g\n", a.Absorbance, a.Transmittance)
	fmt.Printf("A_obs = %g   T_obs = %g   s = %g\n",
		a.ObservedAbsorbance, a.ObservedTransmittance, a.StrayFraction)
	dev := a.Deviation
	if a.Ideal {
		dev = 0
	}
	fmt.Printf("deviation = %+.4g  (%s)\n", dev, deviationLabel(a))
}

func deviationLabel(a mixture.Analysis) string {
	if a.Ideal {
		return "ideal instrument"
	}
	if a.Deviation < 0 {
		return "stray light suppresses high A"
	}
	return "stray light negligible at this A"
}
