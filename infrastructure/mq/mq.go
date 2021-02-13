package mq

import (
	"io"

	"github.com/letian0805/seckill/infrastructure/pool"
)

type Queue interface {
	Producer
	Consumer
	io.Closer
}

type Producer interface {
	Produce(task pool.Task) error
}

type Consumer interface {
	Consume() (pool.Task, error)
}

type Factory interface {
	New(name string) (Queue, error)
	NewProducer(name string) (Producer, error)
	NewConsumer(name string) (Consumer, error)
}

type FactoryFunc func(name string) (Queue, error)

func (f FactoryFunc) New(name string) (Queue, error) {
	return f(name)
}

func (f FactoryFunc) NewProducer(name string) (Producer, error) {
	return f.New(name)
}

func (f FactoryFunc) NewConsumer(name string) (Consumer, error) {
	return f.New(name)
}

var queueFactories = make(map[string]Factory)

func Register(tp string, f Factory) {
	if _, ok := queueFactories[tp]; ok {
		panic("duplicate queue factory " + tp)
	}
	queueFactories[tp] = f
}

func NewFactory(tp string) Factory {
	return queueFactories[tp]
}
