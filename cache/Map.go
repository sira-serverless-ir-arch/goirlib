package cache

import (
	"sync"
	"time"
)

//type CacheMap struct {
//	store sync.Map
//}
//
//func (c *CacheMap) Set(key string, value interface{}) {
//	c.store.Store(key, value)
//}
//
//func (c *CacheMap) Get(key string) (interface{}, bool) {
//	return c.store.Load(key)
//}

type SyncMap[T any] struct {
	Data map[string]T
	Ch   chan TransferData[T]
	sync.RWMutex
}

func NewSyncMap[T any]() *SyncMap[T] {
	i := &SyncMap[T]{
		Data: make(map[string]T),
		Ch:   make(chan TransferData[T], 1000),
	}

	//go i.writer()
	return i
}

type TransferData[T any] struct {
	Key   string
	Value T
}

func (i *SyncMap[T]) writer() {

	m := sync.Mutex{}
	buffer := make(map[string]T)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case transferData := <-i.Ch:
			m.Lock()
			buffer[transferData.Key] = transferData.Value
			m.Unlock()
		case <-ticker.C:

			m.Lock()
			bufferCopy := make(map[string]T, len(i.Data)+len(buffer))
			for key, value := range buffer {
				bufferCopy[key] = value
			}
			buffer = make(map[string]T)
			m.Unlock()

			i.RLock()
			for key, value := range i.Data {
				bufferCopy[key] = value
			}
			i.RUnlock()

			i.Lock()
			i.Data = bufferCopy
			i.Unlock()
		}
	}
}

func (i *SyncMap[T]) Get(key string) (*T, bool) {
	i.RLock()
	defer i.RUnlock()
	if v, ok := i.Data[key]; ok {
		return &v, ok
	}
	return nil, false
}

func (i *SyncMap[T]) Put(key string, value T) {
	i.Ch <- TransferData[T]{
		Key:   key,
		Value: value,
	}
}
