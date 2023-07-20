package cache

import (
	"fmt"
	"hash/crc32"
)

type Shard[T any] struct {
	Fragments int
	Shard     map[int]*LocalMap[T]
}

func NewShardMap[T any](fragments int) *Shard[T] {
	shardMap := &Shard[T]{
		Fragments: fragments,
		Shard:     make(map[int]*LocalMap[T]),
	}

	for i := 0; i < fragments; i++ {
		shardMap.Shard[i] = NewLocalMap[T]()
	}

	return shardMap
}

//func (m *Shard[T]) Put(key string, value *RCU[T]) {
//	shardId := m.getShardId(key)
//	m.Shard[shardId] = value
//}

func (m *Shard[T]) Get(key string) (*LocalMap[T], bool) {
	shardId := m.getShardId(key)
	if iMap, ok := m.Shard[shardId]; ok {
		return iMap, ok
	}

	panic(fmt.Sprintf("Shard not found for key %v", key))

}

func (m *Shard[T]) getShardId(key string) int {
	hash := crc32.ChecksumIEEE([]byte(key))
	return int(hash) % m.Fragments
}
