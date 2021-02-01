package pool_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/letian0805/seckill/infrastructure/pool"
)

type testTask struct {
	wg *sync.WaitGroup
	ch chan struct{}
	m  bool
	p  int
}

func (t *testTask) Do() {
	if t.m {
		<-t.ch
	}
	t.wg.Done()
}

func newTestTask(wg *sync.WaitGroup) *testTask {
	return &testTask{
		wg: wg,
	}
}

func newMemTask(ch chan struct{}, wg *sync.WaitGroup) *testTask {
	return &testTask{
		wg: wg,
		ch: ch,
		m:  true,
	}
}

func (t *testTask) Priority() int {
	return t.p
}

func newPriorityTask(p int, wg *sync.WaitGroup) *testTask {
	return &testTask{
		wg: wg,
		p:  p,
	}
}

func runTest(b *testing.B, f func(i int, wg *sync.WaitGroup)) *sync.WaitGroup {
	//初始化
	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()
	wg := &sync.WaitGroup{}
	wg.Add(b.N)

	//执行测试
	for i := 0; i < b.N; i++ {
		f(i, wg)
	}

	//输出内存信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	b.ReportMetric(float64(memStats.HeapInuse)/(1024*1024), "heap(MB)")
	b.ReportMetric(float64(memStats.StackInuse)/(1024*1024), "stack(MB)")
	return wg
}

func BenchmarkNoGoroutine(b *testing.B) {
	wg := runTest(b, func(i int, wg *sync.WaitGroup) {
		t := newTestTask(wg)
		t.Do()
	})
	wg.Wait()
	b.StopTimer()
}

const (
	priority = 2
	number   = 2 * priority
)

func BenchmarkWorker(b *testing.B) {
	w := pool.NewWorker(number, b.N)
	wg := runTest(b, func(i int, wg *sync.WaitGroup) {
		w.Push(newTestTask(wg))
	})
	w.Close()
	wg.Wait()
	b.StopTimer()
}

func BenchmarkPriorityWorker(b *testing.B) {
	w := pool.NewPriorityWorker(number, b.N, priority)
	wg := runTest(b, func(i int, wg *sync.WaitGroup) {
		w.Push(newPriorityTask(i%priority, wg))
	})
	w.Close()
	wg.Wait()
	b.StopTimer()
}

func BenchmarkWorkerMem(b *testing.B) {
	w := pool.NewWorker(number, b.N)
	ch := make(chan struct{})
	wg := runTest(b, func(i int, wg *sync.WaitGroup) {
		w.Push(newMemTask(ch, wg))
	})
	close(ch)
	w.Close()
	wg.Wait()
	b.StopTimer()
}

func BenchmarkGoroutineCPU(b *testing.B) {
	wg := runTest(b, func(i int, wg *sync.WaitGroup) {
		go func() {
			t := newTestTask(wg)
			t.Do()
		}()
	})
	wg.Wait()
	b.StopTimer()
}

func BenchmarkGoroutineMem(b *testing.B) {
	ch := make(chan struct{})
	wg := runTest(b, func(i int, wg *sync.WaitGroup) {
		go func() {
			t := newMemTask(ch, wg)
			t.Do()
		}()
	})
	close(ch)
	wg.Wait()
	b.StopTimer()
}
