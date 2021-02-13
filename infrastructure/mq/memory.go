package mq

import (
	"errors"
	"fmt"

	"github.com/letian0805/seckill/infrastructure/pool"
	"github.com/letian0805/seckill/infrastructure/utils"
	"github.com/spf13/viper"
)

type memoryQueue struct {
	q utils.RateLimiter
}

func init() {
	Register("memory", FactoryFunc(memoryQueueFactory))
}

func memoryQueueFactory(name string) (Queue, error) {
	rate := viper.GetInt64(fmt.Sprintf("queue.%s.rate", name))
	size := viper.GetInt(fmt.Sprintf("queue.%s.size", name))
	q, _ := utils.NewRateLimiter(size, rate, utils.FanIn)
	mq := &memoryQueue{
		q: q,
	}

	return mq, nil
}

func (mq *memoryQueue) Produce(task pool.Task) error {
	if ok := mq.q.Push(task); !ok {
		return errors.New("queue producer error")
	}
	return nil
}

func (mq *memoryQueue) Consume() (pool.Task, error) {
	t, ok := mq.q.Pop()
	if !ok {
		return nil, errors.New("queue consumer error")
	}
	return t, nil
}

func (mq *memoryQueue) Close() error {
	return mq.q.Close()
}
