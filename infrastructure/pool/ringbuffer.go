package pool

import (
	"sync/atomic"
	"unsafe"
)

type RingBuffer struct {
	count int32
	size  int32
	head  int32
	tail  int32
	buf   []unsafe.Pointer
}

func NewRingBuffer(size int32) *RingBuffer {
	return &RingBuffer{
		size: size,
		head: 0,
		tail: 0,
		buf:  make([]unsafe.Pointer, size),
	}
}

// Get方法从buf中取出对象
func (r *RingBuffer) Get() interface{} {
	// 在高并发开始的时候，队列容易空，直接判断空性能最优
	if atomic.LoadInt32(&r.count) <= 0 {
		return nil
	}

	// 当扣减数量后没有超，就从队列里取出对象
	if atomic.AddInt32(&r.count, -1) >= 0 {
		idx := (atomic.AddInt32(&r.head, 1) - 1) % r.size
		if obj := atomic.LoadPointer(&r.buf[idx]); obj != unsafe.Pointer(nil) {
			o := *(*interface{})(obj)
			atomic.StorePointer(&r.buf[idx], nil)
			return o
		}
	} else {
		// 当减数量超了，再加回去
		atomic.AddInt32(&r.count, 1)
	}
	return nil
}

// Put方法将对象放回到buf中。如果buf满了，返回false
func (r *RingBuffer) Put(obj interface{}) bool {
	// 在高并发结束的时候，队列容易满，直接判满性能最优
	if atomic.LoadInt32(&r.count) >= r.size {
		return false
	}
	// 当增加数量后没有超，就将对象放到队列里
	if atomic.AddInt32(&r.count, 1) <= r.size {
		idx := (atomic.AddInt32(&r.tail, 1) - 1) % r.size
		atomic.StorePointer(&r.buf[idx], unsafe.Pointer(&obj))
		return true
	}
	// 当加的数量超了，再减回去
	atomic.AddInt32(&r.count, -1)
	return false
}
