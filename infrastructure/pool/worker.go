package pool

import (
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

type Worker interface {
	Push(t Task) bool
	Close() error
}

type Task interface {
	Do()
}

type TaskFunc func()

func (tf TaskFunc) Do() {
	tf()
}

type worker struct {
	number   int
	size     int
	closed   int32
	taskPool chan Task
	wg       sync.WaitGroup
}

const (
	minBufferSize = 10
	minNumber     = 2
)

func NewWorker(number int, size int) Worker {
	if number < minNumber {
		number = minNumber
	}
	if size < minBufferSize {
		size = minNumber
	}
	w := &worker{
		number:   number,
		size:     size,
		taskPool: make(chan Task, size),
	}
	w.wg.Add(number)
	for i := 0; i < number; i++ {
		go w.run()
	}
	return w
}

func (w *worker) run() {
	defer w.wg.Done()
	for task := range w.taskPool {
		w.process(task)
	}
}

func (w *worker) process(t Task) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	t.Do()
}

func (w *worker) Push(t Task) bool {
	if w.isClosed() {
		return false
	}

	w.taskPool <- t

	return true
}

func (w *worker) Close() error {
	if !w.isClosed() && atomic.CompareAndSwapInt32(&w.closed, 0, 1) {
		close(w.taskPool)
		w.wg.Wait()
	}
	return nil
}

func (w *worker) isClosed() bool {
	return atomic.LoadInt32(&w.closed) == 1
}

type PriorityTask interface {
	Priority() int
	Do()
}

type priorityWorker struct {
	priorities int
	number     int
	size       int
	closed     int32
	workers    []Worker
}

func NewPriorityWorker(number, size, priorities int) Worker {
	if priorities < minNumber {
		priorities = minNumber
	}

	number = (number - 1 + priorities) / priorities
	if number < minNumber {
		number = minNumber
	}

	size = (size - 1 + priorities) / priorities
	if size < minBufferSize {
		size = minBufferSize
	}

	w := &priorityWorker{
		priorities: priorities,
		number:     number,
		size:       size,
		closed:     0,
		workers:    make([]Worker, priorities),
	}
	for i := 0; i < priorities; i++ {
		w.workers[i] = NewWorker(number, size)
	}
	return w
}

func (pw *priorityWorker) Push(t Task) bool {
	if pw.isClosed() {
		return false
	}
	if pt, ok := t.(PriorityTask); !ok {
		return pw.workers[pw.priorities-1].Push(t)
	} else {
		p := pt.Priority()
		if p < 0 {
			p = 0
		} else if p >= pw.priorities {
			p = pw.priorities - 1
		}
		return pw.workers[p].Push(t)
	}
}

func (pw *priorityWorker) Close() error {
	if !pw.isClosed() && atomic.CompareAndSwapInt32(&pw.closed, 0, 1) {
		for _, w := range pw.workers {
			w.Close()
		}
		return nil
	}
	return nil
}

func (pw *priorityWorker) isClosed() bool {
	return atomic.LoadInt32(&pw.closed) == 1
}
