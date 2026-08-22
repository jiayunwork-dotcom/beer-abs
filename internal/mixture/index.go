package mixture

import "fmt"

// IndexError reports an out-of-range component index.
type IndexError struct {
	// Index is the requested position.
	Index int

	// Length is the number of components in the mixture.
	Length int
}

// Error implements the error interface.
func (e *IndexError) Error() string {
	return fmt.Sprintf("component index %d out of range (mixture has %d components)", e.Index, e.Length)
}

// OutOfRange reports whether the index fell outside the valid range.
func (e *IndexError) OutOfRange() bool {
	return e.Index < 0 || e.Index >= e.Length
}
