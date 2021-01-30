package stores

import (
	"sync"
	"sync/atomic"
)

type IntCache interface {
	Get(key string) (int64, bool)
	Set(key string, val int64)
	Add(key string, delta int64) int64
	Del(key string)
	Keys() []string
}

type ObjCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	Del(key string)
	Keys() []string
}

type intCache struct {
	sync.RWMutex
	data map[string]*int64
}

func NewIntCache() IntCache {
	return &intCache{
		data: make(map[string]*int64),
	}
}

func (c *intCache) getPtr(key string) *int64 {
	c.RLock()
	vp, _ := c.data[key]
	c.RUnlock()
	return vp
}

func (c *intCache) Set(key string, val int64) {
	vp := c.getPtr(key)
	if vp != nil {
		atomic.StoreInt64(vp, val)
	} else {
		vp = new(int64)
		*vp = val
		c.Lock()
		c.data[key] = vp
		c.Unlock()
	}
}

func (c *intCache) Get(key string) (int64, bool) {
	vp := c.getPtr(key)
	if vp != nil {
		return atomic.LoadInt64(vp), true
	}
	return 0, false
}

func (c *intCache) Add(key string, delta int64) int64 {
	vp := c.getPtr(key)
	if vp != nil {
		return atomic.AddInt64(vp, delta)
	} else {
		var val int64
		var ok bool
		c.Lock()
		if vp, ok = c.data[key]; ok {
			val = atomic.AddInt64(vp, delta)
		} else {
			val = delta
			vp = &val
			c.data[key] = vp
		}
		c.Unlock()
		return val
	}
}

func (c *intCache) Del(key string) {
	vp := c.getPtr(key)
	if vp != nil {
		c.Lock()
		delete(c.data, key)
		c.Unlock()
	}
}

func (c *intCache) Keys() []string {
	keys := make([]string, 0)
	c.RLock()
	for k, _ := range c.data {
		keys = append(keys, k)
	}
	c.RUnlock()
	return keys
}

type objCache struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewObjCache() ObjCache {
	return &objCache{
		data: make(map[string]interface{}),
	}
}

func (oc *objCache) Set(key string, data interface{}) {
	oc.Lock()
	oc.data[key] = data
	oc.Unlock()
}

func (oc *objCache) Get(key string) (interface{}, bool) {
	oc.RLock()
	v, ok := oc.data[key]
	oc.RUnlock()
	return v, ok
}

func (oc *objCache) Del(key string) {
	if _, ok := oc.Get(key); ok {
		oc.Lock()
		delete(oc.data, key)
		oc.Unlock()
	}
}

func (oc *objCache) Keys() []string {
	keys := make([]string, 0)
	oc.RLock()
	for k, _ := range oc.data {
		keys = append(keys, k)
	}
	oc.RUnlock()
	return keys
}
