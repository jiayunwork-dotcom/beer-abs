package mixture

import "fmt"

type IndexError struct {
	Index int

	Length int
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("component index %d out of range (mixture has %d components)", e.Index, e.Length)
}

func (e *IndexError) OutOfRange() bool {
	return e.Index < 0 || e.Index >= e.Length
}
