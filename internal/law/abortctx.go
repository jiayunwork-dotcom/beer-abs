package law

import "context"

func abortAbsContext() error {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
