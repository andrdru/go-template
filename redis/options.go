package redis

import (
	"time"
)

type (
	options struct {
		timeout time.Duration
	}

	Option func(*options)
)

func WithTimeout(timeout time.Duration) Option {
	return func(args *options) {
		args.timeout = timeout
	}
}
