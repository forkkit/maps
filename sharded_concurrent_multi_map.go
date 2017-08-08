package maps

import (
	"github.com/shomali11/util/hashes"
)

// NewShardedConcurrentMultiMap creates a new sharded concurrent map
func NewShardedConcurrentMultiMap(options ...ShardOption) *ShardedConcurrentMultiMap {
	shardedConcurrentMultiMap := &ShardedConcurrentMultiMap{
		shards: getNumberOfShards(options...),
	}

	concurrentMaps := make([]*ConcurrentMultiMap, shardedConcurrentMultiMap.shards)
	for i := uint32(0); i < shardedConcurrentMultiMap.shards; i++ {
		concurrentMaps[i] = NewConcurrentMultiMap()
	}

	shardedConcurrentMultiMap.concurrentMaps = concurrentMaps
	return shardedConcurrentMultiMap
}

// ShardedConcurrentMultiMap concurrent map
type ShardedConcurrentMultiMap struct {
	shards         uint32
	concurrentMaps []*ConcurrentMultiMap
}

// Set concurrent set to map
func (c *ShardedConcurrentMultiMap) Set(key string, values []interface{}) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Set(key, values)
}

// Append concurrent append to map
func (c *ShardedConcurrentMultiMap) Append(key string, value interface{}) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Append(key, value)
}

// Get concurrent get from map
func (c *ShardedConcurrentMultiMap) Get(key string) ([]interface{}, bool) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	return concurrentMap.Get(key)
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMultiMap) Remove(key string) {
	shard := c.getShard(key)
	concurrentMap := c.concurrentMaps[shard]
	concurrentMap.Remove(key)
}

// Contains concurrent contains in map
func (c *ShardedConcurrentMultiMap) Contains(key string) bool {
	_, ok := c.Get(key)
	return ok
}

// Size concurrent size of map
func (c *ShardedConcurrentMultiMap) Size() int {
	sum := 0
	for _, concurrentMap := range c.concurrentMaps {
		sum += concurrentMap.Size()
	}
	return sum
}

// Remove concurrent remove from map
func (c *ShardedConcurrentMultiMap) Clear() {
	for _, concurrentMap := range c.concurrentMaps {
		concurrentMap.Clear()
	}
}

func (c *ShardedConcurrentMultiMap) getShard(key string) uint32 {
	return hashes.FNV32(key) % uint32(c.shards)
}
