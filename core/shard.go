package core

import (
	"sync"

	"github.com/zm50/gte/common/types"
	"github.com/zm50/gte/trait"
)

// KVShard 键值对分片
type KVShard[K types.Integer, V any] struct {
	items map[K]V
	sync.RWMutex
}

var _ trait.KVShard[int, int] = (*KVShard[int, int])(nil)

// NewKVShard 创建一个新的键值对分片
func NewKVShard[K types.Integer, V any]() trait.KVShard[K, V] {
	return &KVShard[K, V]{
		items: make(map[K]V),
		RWMutex: sync.RWMutex{},
	}
}

// Get 基于键获取分片中的值
func (s *KVShard[K, V]) Get(key K) (V, bool) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.items[key]
	return value, ok
}

// Set 基于键设置分片中的值
func (s *KVShard[K, V]) Set(key K, value V) {
	s.Lock()
	defer s.Unlock()

	s.items[key] = value
}

// Del 基于键删除分片中的键值对
func (s *KVShard[K, V]) Del(key K) {
	s.Lock()
	defer s.Unlock()

	delete(s.items, key)
}

// RRange 加读锁并遍历分片中所有的键值对
func (s *KVShard[K, V]) RRange(fn func (K, V)) {
	s.RLock()
	defer s.RUnlock()

	for key, value := range s.items {
		fn(key, value)
	}
}

// WRange 加写锁并遍历分片中所有的键值对
func (s *KVShard[K, V]) WRange(fn func (K, V)) {
	s.Lock()
	defer s.Unlock()

	for key, value := range s.items {
		fn(key, value)
	}
}

// RRangeKeys 加读锁并遍历分片中指定的键值对
func (s *KVShard[K, V]) RRangeKeys(fn func (K, V, bool), keys ...K) {
	s.RLock()
	defer s.RUnlock()

	for _, key := range keys {
		if value, ok := s.items[key]; ok {
			fn(key, value, ok)
		}
	}
}

// WRangeKeys 加写锁并遍历分片中指定的键值对
func (s *KVShard[K, V]) WRangeKeys(fn func (K, V, bool), keys ...K) {
	s.Lock()
	defer s.Unlock()

	for _, key := range keys {
		if value, ok := s.items[key]; ok {
			fn(key, value, ok)
		}
	}
}

// Items 获取分片中所有的键值对
func (s *KVShard[K, V]) Items() map[K]V {
	return s.items
}

// KVShards 键值对分片集合
type KVShards[K types.Integer, V any] struct {
	shards []trait.KVShard[K, V]
}

var _ trait.KVShards[int, any] = (*KVShards[int, any])(nil)

// NewKVShards 创建一个新的键值对分片集合
func NewKVShards[K types.Integer, V any](numShards int) *KVShards[K, V] {
	shards := make([]trait.KVShard[K, V], numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewKVShard[K, V]()
	}

	return &KVShards[K, V]{
		shards: shards,
	}
}

// GetShard 获取分片
func (s *KVShards[K, V]) GetShard(key K) trait.KVShard[K, V] {
	return s.shards[key % K(len(s.shards))]
}

// Get 基于键获取分片集合中的值
func (s *KVShards[K, V]) Get(key K) (V, bool) {
	return s.GetShard(key).Get(key)
}

// Set 基于键设置分片集合中的值
func (s *KVShards[K, V]) Set(key K, value V) {
	s.GetShard(key).Set(key, value)
}

// Del 基于键删除分片集合中的键值对
func (s *KVShards[K, V]) Del(key K) {
	s.GetShard(key).Del(key)
}

// Count 计算分片集合中键值对的数量
func (s *KVShards[K, V]) Count() int {
	count := 0
	for i := 0; i < len(s.shards); i++ {
		shard := s.shards[i]
		shard.RLock()
		count += len(shard.Items())
		shard.RUnlock()
	}
	return count
}

// Keys 获取分片集合中所有键的迭代器
func (s *KVShards[K, V]) KeysIter(n int) <- chan K {
	keysCh := make(chan K, n)
	go func() {
		n := len(s.shards)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			shard := s.shards[i]
			go func(shard trait.KVShard[K, V]) {
				shard.RLock()
				for key := range shard.Items() {
					keysCh <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(keysCh)
	}()

	return keysCh
}

// Values 获取分片集合中所有值的迭代器
func (s *KVShards[K, V]) ValuesIter(n int) <- chan V {
	valuesCh := make(chan V, n)
	go func() {
		n := len(s.shards)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			shard := s.shards[i]
			go func(shard trait.KVShard[K, V]) {
				shard.RLock()
				for _, value := range shard.Items() {
					valuesCh <- value
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(valuesCh)
	}()

	return valuesCh
}

// ItemsIter 获取分片集合中所有键值对的迭代器
func (s *KVShards[K, V]) ItemsIter(n int) <- chan *types.KVItem[K, V] {
	itemsCh := make(chan *types.KVItem[K, V], n)
	go func() {
		n := len(s.shards)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			shard := s.shards[i]
			go func(shard trait.KVShard[K, V]) {
				shard.RLock()
				for key, value := range shard.Items() {
					itemsCh <- &types.KVItem[K, V]{Key: key, Value: value}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(itemsCh)
	}()

	return itemsCh
}

// Keys 获取分片集合中所有键
func (s *KVShards[K, V]) Keys() []K {
	count := s.Count()
	keys := make([]K, 0, count)
	for key := range s.KeysIter(count) {
		keys = append(keys, key)
	}

	return keys
}

// Values 获取分片集合中所有值
func (s *KVShards[K, V]) Values() []V {
	count := s.Count()
	values := make([]V, 0, count)
	for value := range s.ValuesIter(count) {
		values = append(values, value)
	}

	return values
}

// Items 获取分片集合中所有键值对
func (s *KVShards[K, V]) Items() []*types.KVItem[K, V] {
	count := s.Count()
	items := make([]*types.KVItem[K, V], 0, count)
	for item := range s.ItemsIter(count) {
		items = append(items, item)
	}

	return items
}

// Shards 获取所有的分片信息
func (s *KVShards[K, V]) Shards() []trait.KVShard[K, V] {
	return s.shards
}

// Range 遍历分片集合中所有的键值对
func (s *KVShards[K, V]) Range(fn func (K, V)) {
	for i := 0; i < len(s.shards); i++ {
		shard := s.shards[i]
		shard.RLock()
		for key, value := range shard.Items() {
			fn(key, value)
		}
		shard.RUnlock()
	}
}
