package utils

import "time"

func Try(fn func() error, retries int, interval time.Duration) (err error) {
	times := 0
	for {
		if retries > 0 && times >= retries {
			break
		}
		if err = fn(); err == nil {
			break
		}
		times++
		time.Sleep(interval)
	}
	return
}
