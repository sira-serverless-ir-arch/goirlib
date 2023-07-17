package cache

import (
	"container/list"
	"sync"
)

type Page[T any] struct {
	key   string
	value T
}

type LRUCache[T any] struct {
	capacity int
	cache    map[string]*list.Element
	pages    *list.List
	sync.RWMutex
}

func NewLRUCache[T any](capacity int) *LRUCache[T] {
	return &LRUCache[T]{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		pages:    list.New(),
	}
}

func (l *LRUCache[T]) Get(key string) (*T, bool) {
	l.RLock()
	defer l.RUnlock()
	if elem, ok := l.cache[key]; ok {
		l.pages.MoveToFront(elem)
		return &elem.Value.(*Page[T]).value, true
	}
	return nil, false
}

func (l *LRUCache[T]) Put(key string, value T) {
	l.Lock()
	defer l.Unlock()
	if elem, ok := l.cache[key]; ok {
		l.pages.MoveToFront(elem)
		elem.Value.(*Page[T]).value = value
		return
	}

	if l.pages.Len() >= l.capacity {
		delete(l.cache, l.pages.Back().Value.(*Page[T]).key)
		l.pages.Remove(l.pages.Back())
	}

	l.cache[key] = l.pages.PushFront(&Page[T]{key: key, value: value})
}
