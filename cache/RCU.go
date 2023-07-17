package cache

import (
	"sync/atomic"
	"time"
)

type RCU[T any] struct {
	data atomic.Value
	Ch   chan TransferData[T]
}

func NewCacheRCU[T any]() *RCU[T] {
	c := &RCU[T]{
		Ch: make(chan TransferData[T], 1000),
	}
	c.data.Store(make(map[string]T))
	go c.writer()
	return c
}

//type TransferData[T any] struct {
//	Key   string
//	Value T
//}

func (c *RCU[T]) Get(key string) (*T, bool) {
	data := c.data.Load().(map[string]T)
	value, ok := data[key]
	return &value, ok
}

func (c *RCU[T]) Put(key string, value T) {
	c.Ch <- TransferData[T]{Key: key, Value: value}
}

func (c *RCU[T]) writer() {
	buffer := make(map[string]T)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case td := <-c.Ch:
			buffer[td.Key] = td.Value
		case <-ticker.C:
			dataCopy := make(map[string]T)
			data := c.data.Load().(map[string]T)

			for key, value := range data {
				dataCopy[key] = value
			}
			for key, value := range buffer {
				dataCopy[key] = value
			}

			c.data.Store(dataCopy)
			buffer = make(map[string]T)
		}
	}
}
