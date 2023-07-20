package cache

import "sync"

type LocalMap[T any] struct {
	mu    sync.RWMutex
	store map[string]*T
}

func NewLocalMap[T any]() *LocalMap[T] {
	return &LocalMap[T]{
		store: make(map[string]*T),
	}
}

func (l *LocalMap[T]) Set(key string, value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.store[key] = &value
}

func (l *LocalMap[T]) GetData() map[string]*T {
	return l.store
}

func (l *LocalMap[T]) Get(key string) (*T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if t, ok := l.store[key]; ok {
		return t, ok
	}
	return nil, false
}

//type LocalMap[T any] struct {
//	store sync.Map
//}
//
//func (l *LocalMap[T]) Set(key string, value T) {
//	l.store.Store(key, value)
//}
//
//func (l *LocalMap[T]) Get(key string) (*T, bool) {
//	data, ok := l.store.Load(key)
//
//	if !ok {
//		return nil, false
//	}
//
//	if res, ok := data.(T); ok {
//		return &res, ok
//	}
//
//	return nil, false
//}
