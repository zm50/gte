package core

import "github.com/zm50/gte/trait"

// ObjPool 对象池，提供对象复用
type ObjPool[T any] struct {
	ch chan T
}

var _ trait.ObjPool[any] = (*ObjPool[any])(nil)

func NewObjPool[T any](size int, objectProvider func() T) *ObjPool[T] {
	pool := &ObjPool[T]{
		ch: make(chan T, size),
	}
	for i := 0; i < size; i++ {
		pool.ch <- objectProvider()
	}

	return pool
}

// Get 从对象池中获取一个对象
func (p *ObjPool[T]) Get() T {
	return <-p.ch
}

// Put 向对象池中放入一个对象
func (p *ObjPool[T]) Put(obj T) {
	if len(p.ch) == cap(p.ch) {
		return
	}

	p.ch <- obj
}
