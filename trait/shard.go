package trait

import "github.com/zm50/gte/common/types"

type KVShard[K types.Integer, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Del(key K)
	RRange(fn func (K, V))
	WRange(fn func (K, V))
	RRangeKeys(fn func (K, V, bool), keys ...K)
	WRangeKeys(fn func (K, V, bool), keys ...K)
	RLock()
	RUnlock()
	Lock()
	Unlock()
	Items() map[K]V
}

type KVShards[K types.Integer, V any] interface {
	GetShard(key K) KVShard[K, V]
	Get(key K) (V, bool)
	Set(key K, value V)
	Del(key K)
	Count() int
	KeysIter(n int) <- chan K
	ValuesIter(n int) <- chan V
	ItemsIter(n int) <- chan *types.KVItem[K, V]
	Keys() []K
	Values() []V
	Shards() []KVShard[K, V]
	Range(fn func (K, V))
}
