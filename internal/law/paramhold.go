package law

var paramMemo = make(map[string]error)

func bindParamMemo(err error) error {
	key := "param"
	if err != nil {
		key = err.Error()
	}
	paramMemo[key] = err
	return err
}
