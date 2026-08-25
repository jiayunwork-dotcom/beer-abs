package mixture

import "strconv"

type BaselineResult struct {
	SampleAbsorbance float64

	BlankAbsorbance float64

	NetAbsorbance float64
}

func SubtractBlank(sample, blank float64) (float64, error) {
	if sample < 0 {
		return 0, &BaselineError{Field: "sample", Value: sample, Reason: "样品吸光度不能为负"}
	}
	if blank < 0 {
		return 0, &BaselineError{Field: "blank", Value: blank, Reason: "空白吸光度不能为负"}
	}
	return sample - blank, nil
}

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

type BaselineError struct {
	Field string

	Value float64

	Reason string
}

func (e *BaselineError) Error() string {
	return "baseline." + e.Field + " = " + strconv.FormatFloat(e.Value, 'g', -1, 64) + ": " + e.Reason
}
