package types

// KVItem 键值对
type KVItem[K Integer, V any] struct {
	Key K
	Value V
}
