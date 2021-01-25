package stock

import (
	"sync"
	"testing"

	"github.com/letian0805/seckill/infrastructure/utils"
)

func TestMemStock(t *testing.T) {
	var (
		st  Stock
		err error
		val int64
	)
	if st, err = NewMemStock("101", "1001"); err != nil {
		t.Fatal(err)
	}
	if err = st.Set(10, 1); err != nil {
		t.Fatal(err)
	}
	if val, err = st.Get(); err != nil {
		t.Fatal(err)
	} else if val != 10 {
		t.Fatal("not equal 10")
	}
	if val, err = st.Sub(); err != nil {
		t.Fatal(err)
	} else if val != 9 {
		t.Fatal("not equal 9")
	}
	if err = st.Del(); err != nil {
		t.Fatal(err)
	}
	if val, err = st.Get(); err != nil {
		t.Fatal(err)
	} else if val != 0 {
		t.Fatal("not equal 0")
	}
}

func BenchmarkMemStock(b *testing.B) {
	b.ReportAllocs()
	b.StartTimer()
	c := utils.NewMemCache(10)
	for i := 0; i < b.N; i++ {
		c.Set("101", int64(i))
		c.Get("101")
		c.Del("101")
	}
	b.StopTimer()
}

func BenchmarkSyncMap(b *testing.B) {
	b.ReportAllocs()
	b.StartTimer()
	c := sync.Map{}
	for i := 0; i < b.N; i++ {
		c.Store("101", int64(i))
		c.Load("101")
		c.Delete("101")
	}
	b.StopTimer()
}
