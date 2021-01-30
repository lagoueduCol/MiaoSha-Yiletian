package pool

import (
	"io"
	"sync/atomic"
)

type ringBufferPool struct {
	closed  int32
	name    string
	rb      *RingBuffer
	newFunc func() (io.Closer, error)
}

func NewRingBufferPool(name string, size int, newFunc func() (io.Closer, error)) Pool {
	return &ringBufferPool{
		name:    name,
		rb:      NewRingBuffer(int32(size)),
		newFunc: newFunc,
	}
}

func (p *ringBufferPool) Get() (io.Closer, error) {
	var err error
	var c io.Closer
	if atomic.LoadInt32(&p.closed) != 0 {
		return nil, Failed
	}
	obj := p.rb.Get()
	if c, _ = obj.(io.Closer); c != io.Closer(nil) {
		return c, err
	} else if p.newFunc != nil {
		return p.newFunc()
	}
	return nil, Failed
}

func (p *ringBufferPool) Put(c io.Closer) {
	if c == io.Closer(nil) {
		return
	}
	if atomic.LoadInt32(&p.closed) != 0 || !p.rb.Put(c) {
		_ = c.Close()
	}
}

func (p *ringBufferPool) Close() error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return nil
	}
	for obj := p.rb.Get(); obj != nil; obj = p.rb.Get() {
		if c, ok := obj.(io.Closer); ok {
			_ = c.Close()
		}
	}
	return nil
}
