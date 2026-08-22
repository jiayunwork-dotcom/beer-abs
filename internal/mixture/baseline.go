package mixture

import "strconv"

// BaselineResult carries the blank-corrected reading of a sample.
type BaselineResult struct {
	// SampleAbsorbance is the raw absorbance of the sample.
	SampleAbsorbance float64

	// BlankAbsorbance is the absorbance of the blank (solvent) run.
	BlankAbsorbance float64

	// NetAbsorbance is Sample − Blank.
	NetAbsorbance float64
}

// SubtractBlank computes the net absorbance after subtracting a blank
// reading. Both absorbances must be non-negative.
func SubtractBlank(sample, blank float64) (float64, error) {
	if sample < 0 {
		return 0, &BaselineError{Field: "sample", Value: sample, Reason: "样品吸光度不能为负"}
	}
	if blank < 0 {
		return 0, &BaselineError{Field: "blank", Value: blank, Reason: "空白吸光度不能为负"}
	}
	return sample - blank, nil
}

// CorrectBlank runs the blank subtraction against a second mixture that
// represents the solvent-only run.
func (m Mixture) CorrectBlank(blank Mixture) (BaselineResult, error) {
	sample, err := m.TotalAbsorbance()
	if err != nil {
		return BaselineResult{}, err
	}
	blankA, err := blank.TotalAbsorbance()
	if err != nil {
		return BaselineResult{}, err
	}
	net, err := SubtractBlank(sample, blankA)
	if err != nil {
		return BaselineResult{}, err
	}
	return BaselineResult{
		SampleAbsorbance: sample,
		BlankAbsorbance:  blankA,
		NetAbsorbance:    net,
	}, nil
}

// BaselineError describes an invalid baseline subtraction input.
type BaselineError struct {
	// Field names the offending quantity.
	Field string

	// Value is the invalid value.
	Value float64

	// Reason is a human readable description of the violated constraint.
	Reason string
}

// Error implements the error interface.
func (e *BaselineError) Error() string {
	return "baseline." + e.Field + " = " + strconv.FormatFloat(e.Value, 'g', -1, 64) + ": " + e.Reason
}
