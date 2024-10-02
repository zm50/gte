package gnet

import (
	"sync"

	"github.com/go75/gte/trait"
)

// ConnShard 连接分片，存储一些客户端连接
type ConnShard struct {
	conns map[int32]trait.Connection
	shardLock sync.RWMutex
}

// ConnShards 连接分片集合
type ConnShards struct {
	shards []*ConnShard
}

func (s *ConnShards) GetShard(connID int32) *ConnShard {
	shardIdx := connID % int32(len(s.shards))

	return s.shards[shardIdx]
}

func (s *ConnShards) GetConn(connID int32, conn trait.Connection) {

}
