package stores

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type cacheValue struct {
	val        interface{}
	expiration int64
}

func (v *cacheValue) expired() bool {
	return atomic.LoadInt64(&v.expiration) <= time.Now().UnixNano()
}

func (v *cacheValue) expire(expiration int64) {
	if expiration > 0 {
		atomic.StoreInt64(&v.expiration, time.Now().UnixNano()+expiration*int64(time.Second))
	} else if expiration <= 0 {
		atomic.StoreInt64(&v.expiration, expiration)
		if expiration == 0 {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&v.val)), nil)
		}
	}
}

type memCache struct {
	sync.RWMutex
	data    map[string]*cacheValue
	ticker  *time.Ticker
	stopped int32
}

type Cache interface {
	Set(key string, data interface{})
	Get(key string) interface{}
	Expire(key string, expiration int64) error
	Del(key string)
	Close()
}

const minGCInterval = 10

func NewMemCache(gcIntervalMS int64) Cache {
	if gcIntervalMS < minGCInterval {
		gcIntervalMS = minGCInterval
	}
	c := &memCache{
		data:   make(map[string]*cacheValue),
		ticker: time.NewTicker(time.Duration(gcIntervalMS) * time.Millisecond),
	}
	go c.startGC()
	return c
}

func (mc *memCache) Set(key string, data interface{}) {
	mc.Lock()
	v, ok := mc.data[key]
	if !ok {
		v = &cacheValue{
			val:        data,
			expiration: -1,
		}
		mc.data[key] = v
	} else {
		v.val = data
	}
	mc.Unlock()
}

func (mc *memCache) Get(key string) interface{} {
	var val interface{}
	mc.RLock()
	if v, ok := mc.data[key]; ok {
		val = v.val
	}
	mc.RUnlock()
	return val
}

func (mc *memCache) Del(key string) {
	// 仅仅标记为失效，让 gc 来回收
	mc.RLock()
	if v, ok := mc.data[key]; ok {
		v.expire(0)
	}
	mc.RUnlock()
}

func (mc *memCache) Expire(key string, expiration int64) error {
	var err error
	mc.RLock()
	v, ok := mc.data[key]
	if ok {
		v.expire(expiration)
	} else {
		err = errors.New("not exists")
	}
	mc.RUnlock()
	return err
}

func (mc *memCache) Close() {
	mc.ticker.Stop()
}

func (mc *memCache) startGC() {
	for _ = range mc.ticker.C {
		mc.doGC()
	}
}

func (mc *memCache) doGC() {
	keys := make([]string, 0)
	mc.RLock()
	for k, v := range mc.data {
		if v.expired() {
			keys = append(keys, k)
		}
	}
	mc.RUnlock()
	for _, k := range keys {
		mc.Lock()
		delete(mc.data, k)
		mc.Unlock()
	}
}
