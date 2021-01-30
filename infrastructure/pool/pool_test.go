package pool

import (
	"io"
	"sync"
	"testing"
)

type testCloser struct {
}

func (c *testCloser) Close() error {
	return nil
}

func newCloser() (io.Closer, error) {
	return &testCloser{}, nil
}

func testPool(b *testing.B, p Pool) {
	var data = make([]io.Closer, b.N, b.N)
	for i := 0; i < b.N; i++ {
		data[i], _ = p.Get()
	}
	ch := make(chan struct{})
	wg1 := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	wg1.Add(b.N)
	wg2.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(c io.Closer) {
			wg1.Done()
			<-ch
			p.Put(c)
			p.Get()
			wg2.Done()
		}(data[i])
	}
	wg1.Wait()
	b.ReportAllocs()
	b.StartTimer()
	close(ch)
	wg2.Wait()
	b.StopTimer()
}

func BenchmarkChanPool(b *testing.B) {
	p := NewChanPool("test", b.N, newCloser)
	testPool(b, p)
}

func BenchmarkRingBufferPool(b *testing.B) {
	p := NewRingBufferPool("test", b.N, newCloser)
	testPool(b, p)
}
