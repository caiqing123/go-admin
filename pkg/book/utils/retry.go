package utils

import (
	"time"
)

func Retry(attempts int, sleep time.Duration, fn func() error) error {

	var (
		err    error
		nsleep = sleep
	)

	if err = fn(); err != nil {
		if s, ok := err.(Stop); ok {
			return s.error
		}

		for attempts--; attempts > 0; attempts-- {
			time.Sleep(nsleep)
			nsleep += sleep
			if err = fn(); err == nil {
				return nil
			}
		}
	}
	return err
}

type Stop struct {
	error
}

func NoRetryError(err error) Stop {
	return Stop{err}
}
