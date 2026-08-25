package spectrum

import "fmt"

var (
	ErrEmptyScan     = fmt.Errorf("spectrum: empty wavelength scan")
	ErrUnequalLength = fmt.Errorf("spectrum: wavelength and value length mismatch")
	ErrBadStep       = fmt.Errorf("spectrum: wavelength step must be positive")
	ErrNoPeak        = fmt.Errorf("spectrum: no interior peak")
)

type SampleError struct {
	Index int
	Field string
	Value float64
}

func (e *SampleError) Error() string {
	return fmt.Sprintf("spectrum: sample %d field %s invalid value %g", e.Index, e.Field, e.Value)
}
