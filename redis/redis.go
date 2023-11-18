package redis

import (
	"errors"
	"time"

	redigoRedis "github.com/gomodule/redigo/redis"
)

type (
	Redis struct {
		pool             *redigoRedis.Pool
		operationTimeout time.Duration
	}
)

var (
	// TimeoutDefault .
	TimeoutDefault = 500 * time.Millisecond
	// NewPoolMaxIdle .
	NewPoolMaxIdle = 3
	// NewPoolIdleTimeout .
	NewPoolIdleTimeout = 240 * time.Second

	// ErrKeyNotFound .
	ErrKeyNotFound = errors.New("key not found")
	// ErrValueInvalidFormat .
	ErrValueInvalidFormat = errors.New("value format invalid")
)

// NewPool init redigo pool
func NewPool(address string) *redigoRedis.Pool {
	return &redigoRedis.Pool{
		MaxIdle:     NewPoolMaxIdle,
		IdleTimeout: NewPoolIdleTimeout,
		Dial: func() (redigoRedis.Conn, error) {
			return redigoRedis.Dial("tcp", address)
		},
		TestOnBorrow: func(c redigoRedis.Conn, t time.Time) error {
			var err error
			_, err = c.Do("PING")
			return err
		},
	}
}

// NewRedis .
func NewRedis(pool *redigoRedis.Pool, opts ...Option) *Redis {
	args := &options{
		timeout: TimeoutDefault,
	}

	for _, opt := range opts {
		opt(args)
	}

	return &Redis{
		pool:             pool,
		operationTimeout: args.timeout,
	}
}

// Get by key
func (r *Redis) Get(key string) (data any, err error) {
	data, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"GET", key,
	)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrKeyNotFound
	}

	return data, nil
}

// SetEx set with ttl
func (r *Redis) SetEx(key string, data any, ttl time.Duration) (err error) {
	_, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"SETEX", key, ttl.Seconds(), data,
	)

	return err
}

// Set key value
func (r *Redis) Set(key string, data any) (err error) {
	_, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"SET", key, data,
	)

	return err
}

// Del delete key
func (r *Redis) Del(key string) (err error) {
	_, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"DEL", key,
	)

	return err
}

// Keys get keys list
func (r *Redis) Keys(pattern string) (list []string, err error) {
	var data any
	data, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"KEYS", pattern,
	)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, ErrKeyNotFound
	}

	var keys, ok = data.([]any)
	if !ok {
		return nil, ErrKeyNotFound
	}

	for _, key := range keys {
		list = append(list, string(key.([]byte)))
	}

	return list, nil
}

func (r *Redis) ExpireAt(key string, t time.Time) (err error) {
	_, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"EXPIREAT", key, t.Unix(),
	)

	return err
}

// Incr increment key
func (r *Redis) Incr(key string) (val int64, err error) {
	var data any
	data, err = redigoRedis.DoWithTimeout(
		r.pool.Get(),
		r.operationTimeout,
		"INCR", key,
	)
	if err != nil {
		return 0, err
	}

	var ok bool
	val, ok = data.(int64)
	if !ok {
		return 0, ErrValueInvalidFormat
	}

	return val, nil
}
