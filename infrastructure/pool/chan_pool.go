package pool

import (
	"io"
)

type chanPool struct {
	name    string
	size    int
	ch      chan io.Closer
	newFunc func() (io.Closer, error)
}

func NewChanPool(name string, size int, newFunc func() (io.Closer, error)) Pool {
	return &chanPool{
		name:    name,
		size:    size,
		ch:      make(chan io.Closer, size),
		newFunc: newFunc,
	}
}

func (p *chanPool) Get() (io.Closer, error) {
	select {
	case c := <-p.ch:
		return c, nil
	default:
		if p.newFunc != nil {
			return p.newFunc()
		}
	}
	return nil, Failed
}

func (p *chanPool) Put(c io.Closer) {
	if c == nil {
		return
	}
	select {
	case p.ch <- c:
		break
	default:
		_ = c.Close()
	}
}

func (p *chanPool) Close() error {
	close(p.ch)
	p.ch = nil
	return nil
}
