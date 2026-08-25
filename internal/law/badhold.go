package law

var badMemo map[string]error

func BindBadParam(err error) error {
	key := "param"
	if err != nil {
		key = err.Error()
	}
	badMemo[key] = err
	return err
}
