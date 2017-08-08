package maps

import (
	"github.com/shomali11/util/hashes"
)

// NewShardedConcurrentMap creates a new sharded concurrent map
func NewShardedConcurrentMap(options ...ShardOption) *ShardedConcurrentMap {
	shardedConcurrentMap := &ShardedConcurrentMap{
		shards: getNumberOfShards(options...),
	}

	concurrentMaps := make([]*ConcurrentMap, shardedConcurrentMap.shards)
	for i := uint32(0); i < shardedConcurrentMap.shards; i++ {
		concurrentMaps[i] = NewConcurrentMap()
	}

	shardedConcurrentMap.concurrentMaps = concurrentMaps
	return shardedConcurrentMap
}

// ShardedConcurrentMap concurrent map
type ShardedConcurrentMap struct {
	shards         uint32
	concurrentMaps []*ConcurrentMap
}

// Set concurrent set to map
func (c *ShardedConcurrentMap) Set(key string, value interface{}) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Set(key, value)
}

// Get concurrent get from map
func (c *ShardedConcurrentMap) Get(key string) (interface{}, bool) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	return concurrentMap.Get(key)
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMap) Remove(key string) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Remove(key)
}

// Contains concurrent contains in map
func (c *ShardedConcurrentMap) Contains(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// Size concurrent size of map
func (c *ShardedConcurrentMap) Size() int {
	sum := 0
	for _, concurrentMap := range c.concurrentMaps {
		sum += concurrentMap.Size()
	}
	return sum
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMap) Clear() {
	for _, concurrentMap := range c.concurrentMaps {
		concurrentMap.Clear()
	}
}

func (c *ShardedConcurrentMap) getShard(key string) uint32 {
	return hashes.FNV32(key) % uint32(c.shards)
}
