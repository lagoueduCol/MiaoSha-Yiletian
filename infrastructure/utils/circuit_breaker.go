package utils

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Counter int64

func (c *Counter) Add() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *Counter) Load() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func (c *Counter) Reset() {
	atomic.StoreInt64((*int64)(c), 0)
}

type CircuitBreaker struct {
	totalCounter Counter
	failsCounter Counter

	duration       int64
	latencyLimit   int64
	totalLimit     int64
	failsRateLimit int64

	recoverFailsRate int64
	lastTime         int64
	allow            int64
}

type CBOption func(cb *CircuitBreaker)

const (
	minDuration  = 100
	minTotal     = 1000
	minFailsRate = 2
)

func WithDuration(duration int64) CBOption {
	return func(cb *CircuitBreaker) {
		cb.duration = duration
	}
}

func WithLatencyLimit(latencyLimit int64) CBOption {
	return func(cb *CircuitBreaker) {
		cb.latencyLimit = latencyLimit
	}
}

func WithFailsLimit(failsRateLimit int64) CBOption {
	return func(cb *CircuitBreaker) {
		cb.failsRateLimit = failsRateLimit
	}
}

func WithTotalLimit(totalLimit int64) CBOption {
	return func(cb *CircuitBreaker) {
		cb.totalLimit = totalLimit
	}
}

func NewCircuitBreaker(opts ...CBOption) *CircuitBreaker {
	cb := &CircuitBreaker{
		totalCounter:   0,
		failsCounter:   0,
		duration:       0,
		lastTime:       0,
		failsRateLimit: 0,
		latencyLimit:   0,
		totalLimit:     0,
		allow:          1,
	}
	for _, opt := range opts {
		opt(cb)
	}
	if cb.duration < minDuration {
		cb.duration = minDuration
	}
	if cb.totalLimit < minTotal {
		cb.totalLimit = minTotal
	}
	if cb.failsRateLimit < minFailsRate {
		cb.failsRateLimit = minFailsRate
	}
	cb.recoverFailsRate = cb.failsRateLimit / 2
	return cb
}

func (cb *CircuitBreaker) Allow(f func() bool) bool {
	fails := cb.failsCounter.Load()
	total := cb.totalCounter.Load()
	start := time.Now().UnixNano() / int64(time.Millisecond)
	if start > cb.lastTime+cb.duration {
		atomic.StoreInt64(&cb.lastTime, start)
		cb.failsCounter.Reset()
		cb.totalCounter.Reset()
		atomic.StoreInt64(&cb.allow, 1)
	}
	cb.totalCounter.Add()
	allow := !(total > 0 && fails*100/cb.failsRateLimit >= total || total >= cb.totalLimit)
	if atomic.LoadInt64(&cb.allow) == 0 {
		if fails*100/cb.recoverFailsRate > total {
			allow = false
		} else if allow {
			atomic.StoreInt64(&cb.allow, 1)
		}
	} else if !allow {
		atomic.StoreInt64(&cb.allow, 0)
	}
	if !allow {
		logrus.Error("not allowed")
		return false
	}
	ok := f()
	end := time.Now().UnixNano() / int64(time.Millisecond)
	if (cb.latencyLimit > 0 && end-start >= cb.latencyLimit) || !ok {
		cb.failsCounter.Add()
	}
	return true
}
