package test

import (
	"fmt"
	"github.com/sira-serverless-ir-arch/goirlib/cache"
	"strconv"
	"testing"
	"time"
)

//func BenchmarkAsyncMapPut(b *testing.B) {
//	m := cache.NewAsyncMap[int]()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		m.Put(strconv.Itoa(i), i)
//	}
//}

//func BenchmarkAsyncMapGet(b *testing.B) {
//	m := cache.NewAsyncMap[int]()
//	for i := 0; i < b.N; i++ {
//		m.Put(strconv.Itoa(i), i)
//	}
//
//	// Wait for all values to be processed by the writer goroutine
//	time.Sleep(4 * time.Second)
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		m.Get(strconv.Itoa(i))
//	}
//}

//	func BenchmarkSyncMapStore(b *testing.B) {
//		var m sync.Map
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ {
//			m.Store(strconv.Itoa(i), i)
//		}
//	}
//
//	func BenchmarkSyncMapLoad(b *testing.B) {
//		var m sync.Map
//		for i := 0; i < b.N; i++ {
//			m.Store(strconv.Itoa(i), i)
//		}
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ {
//			m.Load(strconv.Itoa(i))
//		}
//	}
func BenchmarkAsyncMap(b *testing.B) {
	asyncMap := cache.NewAsyncMap[int]()

	numGoroutines := 1000
	b.SetParallelism(numGoroutines)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < 1000; i++ {
				asyncMap.Put(strconv.Itoa(i), i)
			}
			for i := 0; i < 1000; i++ {
				asyncMap.Get(strconv.Itoa(i))
			}
		}
	})
}

func TestXyz(t *testing.T) {
	c := cache.NewCacheRCU[int]()
	c.Put("a", 2)
	c.Put("b", 3)
	c.Put("c", 4)
	time.Sleep(5 * time.Second)
	v, _ := c.Get("b")
	fmt.Printf("%v\n", *v)
}

func BenchmarkRCUCacheConcurrent(b *testing.B) {
	c := cache.NewCacheRCU[int]()

	numGoroutines := 1000
	b.SetParallelism(numGoroutines)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < 1000; i++ {
				c.Put(strconv.Itoa(i), i)
			}
			for i := 0; i < 1000; i++ {
				c.Get(strconv.Itoa(i))
			}
		}
	})
}

func BenchmarkAsyncMapGet(b *testing.B) {
	m := cache.NewAsyncMap[int]()
	for i := 0; i < 1000; i++ {
		m.Put(fmt.Sprint(i), i)
	}

	// Wait for 5 seconds to let writer goroutine of AsyncMap do its work
	time.Sleep(5 * time.Second)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = m.Get(fmt.Sprint(b.N % 1000))
		}
	})
}

func BenchmarkCacheGet(b *testing.B) {
	c := cache.NewCacheRCU[int]()
	for i := 0; i < 1000; i++ {
		c.Put(fmt.Sprint(i), i)
	}

	// Wait for 5 seconds to let writer goroutine of CacheRCU do its work
	time.Sleep(5 * time.Second)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = c.Get(fmt.Sprint(b.N % 1000))
		}
	})
}

func BenchmarkAsyncMapPut(b *testing.B) {
	m := cache.NewAsyncMap[int]()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Put(fmt.Sprint(b.N), b.N)
		}
	})
}

func BenchmarkCachePut(b *testing.B) {
	c := cache.NewCacheRCU[int]()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Put(fmt.Sprint(b.N), b.N)
		}
	})
}
