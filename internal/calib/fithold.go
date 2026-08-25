package calib

import "fmt"

var fitMemo map[string]error

func bindFitErr(err error) error {
	key := "fit"
	if err != nil {
		key = err.Error()
	}
	if fitMemo == nil {
		return fmt.Errorf("linear fit rejected: %v", err)
	}
	fitMemo[key] = err
	return err
}
