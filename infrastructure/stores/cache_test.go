package stores_test

import (
	"os"
	"strconv"
	"sync"
	"testing"

	. "github.com/letian0805/seckill/infrastructure/stores"
)

func TestIntCache(t *testing.T) {
	c := NewIntCache()
	key := "test"
	c.Set(key, 1)
	if v, ok := c.Get(key); !ok || v != 1 {
		t.Fatal("failed")
	}
	if v := c.Add(key, 5); v != 6 {
		t.Fatal("failed")
	}
	c.Del(key)
	if _, ok := c.Get(key); ok {
		t.Fatal("failed")
	}
}

func TestIntCache_Add(t *testing.T) {
	cache := NewIntCache()
	cases := []struct {
		key    string
		delta  int64
		expect int64
	}{
		{"test1", 0, 0},
		{"test1", 1, 1},
		{"test1", -1, 0},
		{"test1", 0, 0},
		{"test2", 1, 1},
		{"test3", -1, -1},
	}
	for _, c := range cases {
		if cache.Add(c.key, c.delta) != c.expect {
			t.Fatal(c)
		}
	}
}

func TestObjCache(t *testing.T) {
	c := NewObjCache()
	key := "test"
	c.Set(key, int64(1))
	if v, ok := c.Get(key); !ok || v.(int64) != 1 {
		t.Fatal("failed")
		t.Error()
	}
	c.Del(key)
	if _, ok := c.Get(key); ok {
		t.Fatal("failed")
	}
}

func initKeys(b *testing.B) []string {
	var keys = make([]string, 0)
	maxKeyStr := os.Getenv("maxKey")
	maxKey, _ := strconv.Atoi(maxKeyStr)
	if maxKey <= 0 {
		maxKey = b.N
	}
	for i := 0; i < maxKey; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	return keys
}

func initIntCache(b *testing.B, c IntCache, keys []string) {
	l := len(keys)
	for i := 0; i < b.N; i++ {
		c.Set(keys[i%l], int64(i))
	}
}

func initSyncMap(b *testing.B, c sync.Map, keys []string) {
	l := len(keys)
	for i := 0; i < b.N; i++ {
		c.Store(keys[i%l], int64(i))
	}
}

func initObjCache(b *testing.B, c ObjCache, keys []string) {
	l := len(keys)
	for i := 0; i < b.N; i++ {
		c.Set(keys[i%l], int64(i))
	}
}

func BenchmarkIntCache_Add(b *testing.B) {
	keys := initKeys(b)
	c := NewIntCache()
	initIntCache(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Add(keys[i%l], 1)
	}
	b.StopTimer()
}

func benchmarkCacheSet(b *testing.B, setter func(key string, val int64), keys []string) {
	b.ReportAllocs()
	b.StartTimer()
	l := len(keys)
	for i := 0; i < b.N; i++ {
		setter(keys[i%l], int64(i))
	}
	b.StopTimer()
}

func BenchmarkCache_Set(b *testing.B) {
	keys := make([]string, b.N, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = strconv.Itoa(i)
	}
	b.ResetTimer()

	b.Run("intCache", func(b *testing.B) {
		c := NewIntCache()
		setter := func(key string, val int64) {
			c.Set(key, val)
		}
		benchmarkCacheSet(b, setter, keys)
	})
	b.Run("objCache", func(b *testing.B) {
		c := NewObjCache()
		setter := func(key string, val int64) {
			c.Set(key, val)
		}
		benchmarkCacheSet(b, setter, keys)
	})
	b.Run("syncMap", func(b *testing.B) {
		c := sync.Map{}
		setter := func(key string, val int64) {
			c.Store(key, val)
		}
		benchmarkCacheSet(b, setter, keys)
	})
}

func BenchmarkIntCache_Set(b *testing.B) {
	keys := initKeys(b)
	c := NewIntCache()

	b.ReportAllocs()
	b.StartTimer()
	initIntCache(b, c, keys)
	b.StopTimer()
}

func BenchmarkObjCache_Set(b *testing.B) {
	keys := initKeys(b)
	c := NewObjCache()

	b.ReportAllocs()
	b.StartTimer()
	initObjCache(b, c, keys)
	b.StopTimer()
}

func BenchmarkSyncMap_Set(b *testing.B) {
	keys := initKeys(b)
	c := sync.Map{}

	b.ReportAllocs()
	b.StartTimer()
	initSyncMap(b, c, keys)
	b.StopTimer()
}

func BenchmarkIntCache_Get(b *testing.B) {
	keys := initKeys(b)
	c := NewIntCache()
	initIntCache(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Get(keys[i%l])
	}
	b.StopTimer()
}

func BenchmarkObjCache_Get(b *testing.B) {
	keys := initKeys(b)
	c := NewObjCache()
	initObjCache(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Get(keys[i%l])
	}
	b.StopTimer()
}

func BenchmarkSyncMap_Get(b *testing.B) {
	keys := initKeys(b)
	c := sync.Map{}
	initSyncMap(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Load(keys[i%l])
	}
	b.StopTimer()
}

func BenchmarkIntCache_Del(b *testing.B) {
	keys := initKeys(b)
	c := NewIntCache()
	initIntCache(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Del(keys[i%l])
	}
	b.StopTimer()
}

func BenchmarkObjCache_Del(b *testing.B) {
	keys := initKeys(b)
	c := NewObjCache()
	initObjCache(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Del(keys[i%l])
	}
	b.StopTimer()
}

func BenchmarkSyncMap_Del(b *testing.B) {
	keys := initKeys(b)
	c := sync.Map{}
	initSyncMap(b, c, keys)
	l := len(keys)

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.Delete(keys[i%l])
	}
	b.StopTimer()
}
