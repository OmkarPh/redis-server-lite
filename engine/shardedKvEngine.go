package engine

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log/slog"
	"strconv"
	"sync"
)

type ShardedEngine struct {
	shardFactor int
	mutex       []sync.Mutex
	kv_stores   []map[string]string
}

func NewShardedEngine(shardFactor int) *ShardedEngine {
	newEngine := &ShardedEngine{
		shardFactor: shardFactor,
		mutex:       make([]sync.Mutex, shardFactor),
		kv_stores:   make([]map[string]string, shardFactor),
	}
	for i := 0; i < shardFactor; i++ {
		newEngine.kv_stores[i] = make(map[string]string)
	}
	return newEngine
}

// @TODO - faster shard algo
func (engine *ShardedEngine) resolveShardIdx(key string) uint32 {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	shardIdx := hasher.Sum32() % uint32(engine.shardFactor)
	slog.Debug(fmt.Sprintf("Shard idx for %s => %d", key, shardIdx))
	return shardIdx
}

func (engine *ShardedEngine) Get(key string) (string, bool) {
	shardIdx := engine.resolveShardIdx(key)
	if value, ok := engine.kv_stores[shardIdx][key]; ok {
		return value, true
	}
	return "", false
}

func (engine *ShardedEngine) Set(key string, value string) {
	shardIdx := engine.resolveShardIdx(key)
	engine.mutex[shardIdx].Lock()
	defer engine.mutex[shardIdx].Unlock()
	engine.kv_stores[shardIdx][key] = value
}

func (engine *ShardedEngine) Incr(key string) (string, error) {
	shardIdx := engine.resolveShardIdx(key)
	engine.mutex[shardIdx].Lock()
	defer engine.mutex[shardIdx].Unlock()

	newValue := 1

	existingValue, found := engine.kv_stores[shardIdx][key]

	if found {
		existingValueInt, err := strconv.Atoi(existingValue)
		if err != nil {
			errString := "ERR value is not an integer or out of range"
			return "", errors.New(errString)
		}
		newValue = existingValueInt + 1
	}

	return strconv.FormatInt(int64(newValue), 10), nil
}

func (engine *ShardedEngine) Decr(key string) (string, error) {
	shardIdx := engine.resolveShardIdx(key)
	engine.mutex[shardIdx].Lock()
	defer engine.mutex[shardIdx].Unlock()

	newValue := -1

	existingValue, found := engine.kv_stores[shardIdx][key]

	if found {
		existingValueInt, err := strconv.Atoi(existingValue)
		if err != nil {
			errString := "ERR value is not an integer or out of range"
			return "", errors.New(errString)
		}
		newValue = existingValueInt - 1
	}

	return strconv.FormatInt(int64(newValue), 10), nil
}
