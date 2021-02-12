package utils

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/letian0805/seckill/infrastructure/pool"
)

type RateLimiter interface {
	Push(t pool.Task) bool
	Pop() (pool.Task, bool)
	Close() error
}

type fanInOut struct {
	sync.RWMutex
	queueIn  chan pool.Task
	queueOut chan pool.Task
	timer    time.Timer
	lastTime int64
	rate     int64
	duration time.Duration
	closed   int64
	mode     int
}

const (
	minRate = 1
	minSize = 10

	FanIn  = 1 << 0
	FanOut = 1 << 1
)

func NewRateLimiter(size int, rate int64, mode int) (RateLimiter, error) {
	modeMask := FanIn | FanOut
	if mode > modeMask || modeMask&mode == 0 {
		return nil, errors.New("wrong flag")
	}
	if rate < minRate {
		rate = minRate
	}
	if size < minSize {
		size = minSize
	}
	f := &fanInOut{
		timer:    time.Timer{},
		lastTime: 0,
		rate:     rate,
		duration: time.Second / time.Duration(rate),
		closed:   0,
		mode:     mode,
	}
	if FanIn&mode != 0 {
		f.queueIn = make(chan pool.Task, size)
	}
	if FanOut&mode != 0 {
		f.queueOut = make(chan pool.Task, size)
	}
	if mode == modeMask {
		go f.exchange()
	}
	return f, nil
}

func (f *fanInOut) Push(t pool.Task) bool {
	if atomic.LoadInt64(&f.closed) == 1 {
		return false
	}
	f.RLock()
	defer f.RUnlock()
	if atomic.LoadInt64(&f.closed) == 1 {
		return false
	}

	if FanIn&f.mode != 0 {
		select {
		case f.queueIn <- t:
			return true
		default:
			return false
		}
	} else {
		f.sleep()
		f.queueOut <- t
		return true
	}
}

func (f *fanInOut) Pop() (pool.Task, bool) {
	if FanOut&f.mode != 0 {
		t, ok := <-f.queueOut
		return t, ok
	} else {
		f.sleep()
		t, ok := <-f.queueIn
		return t, ok
	}
}

func (f *fanInOut) sleep() {
	now := time.Now().UnixNano()
	delta := f.duration - time.Duration(now-atomic.LoadInt64(&f.lastTime))
	if delta > time.Millisecond {
		time.Sleep(delta)
	}
	atomic.StoreInt64(&f.lastTime, now)
}

func (f *fanInOut) exchange() {
	for t := range f.queueIn {
		f.sleep()
		f.queueOut <- t
	}
	close(f.queueOut)
}

func (f *fanInOut) Close() error {
	f.Lock()
	defer f.Unlock()
	if atomic.CompareAndSwapInt64(&f.closed, 0, 1) {
		if f.mode&FanIn != 0 {
			close(f.queueIn)
		} else if f.mode == FanOut {
			close(f.queueOut)
		}
	}
	return nil
}
