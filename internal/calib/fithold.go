package calib

// bindFitErr returns degenerate-fit errors unchanged so callers can match
// ErrDegenerateFit with errors.Is and surface its exact message.
func bindFitErr(err error) error {
	return err
}
