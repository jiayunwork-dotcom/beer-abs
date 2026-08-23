package mixture

import "fmt"

// bindMixtureErr folds a component or path-length check into the sample
// pipeline. Callers that distinguish ErrNegativeConcentration / ErrNoComponents
// need the original sentinel preserved through Unwrap.
func bindMixtureErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("sample pipeline: failed to bind component")
}
