package pool

import (
	"io"
)

type Pool interface {
	Get() (io.Closer, error)
	Put(c io.Closer)
	Close() error
}

type poolError string

func (e poolError) Error() string {
	return string(e)
}

const Failed = poolError("failed to get connection from pool")
