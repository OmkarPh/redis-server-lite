package store

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/OmkarPh/redis-lite/utils"
	"github.com/cespare/xxhash/v2"
)

type ShardedKvStore struct {
	shardFactor uint64
	mutex       []sync.RWMutex
	kv_stores   []map[string]StoredValue
}

func NewShardedKvStore(shardFactor int) *ShardedKvStore {
	newKvStore := &ShardedKvStore{
		shardFactor: uint64(shardFactor),
		mutex:       make([]sync.RWMutex, shardFactor),
		kv_stores:   make([]map[string]StoredValue, shardFactor),
	}
	for i := 0; i < shardFactor; i++ {
		newKvStore.kv_stores[i] = make(map[string]StoredValue)
	}
	return newKvStore
}

// var shardIdxCache map[string]uint64 = make(map[string]uint64)
func (kvStore *ShardedKvStore) resolveShardIdx(key string) uint64 {
	// if shardIdx, ok := shardIdxCache[key]; ok {
	// 	return shardIdx
	// }
	hasher := xxhash.New()
	hasher.Write([]byte(key))
	shardIdx := hasher.Sum64() % kvStore.shardFactor
	// shardIdxCache[key] = shardIdx
	slog.Debug(fmt.Sprintf("Shard idx for %s => %d", key, shardIdx))
	return shardIdx
}

func (kvStore *ShardedKvStore) Get(key string) (string, bool) {
	shardIdx := kvStore.resolveShardIdx(key)
	if data, ok := kvStore.kv_stores[shardIdx][key]; ok {
		if utils.IsExpired(data.Expiry) {
			kvStore.mutex[shardIdx].Lock()
			defer kvStore.mutex[shardIdx].Unlock()
			delete(kvStore.kv_stores[shardIdx], key)
			return "", false
		}
		return data.Value, true
	}
	return "", false
}

func (kvStore *ShardedKvStore) Set(key string, value string) {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()
	kvStore.kv_stores[shardIdx][key] = StoredValue{
		Value:  value,
		Expiry: time.Time{},
	}
}

func (kvStore *ShardedKvStore) Del(key string) bool {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()
	_, existed := kvStore.kv_stores[shardIdx][key]
	delete(kvStore.kv_stores[shardIdx], key)
	return existed
}

func (kvStore *ShardedKvStore) Expire(key string, seconds int64) (bool, error) {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()

	newExpiry := time.Now().Add(time.Second * time.Duration(seconds))

	existingData, exists := kvStore.kv_stores[shardIdx][key]
	if !exists {
		return false, errors.New("ERR key doesn't exist")
	}

	kvStore.kv_stores[shardIdx][key] = StoredValue{
		Value:  existingData.Value,
		Expiry: newExpiry,
	}

	return true, nil
}

func (kvStore *ShardedKvStore) Ttl(key string) int {
	shardIdx := kvStore.resolveShardIdx(key)

	/*
		-2 => Key doesn't exist
		-1 => No expiry set
		0,1,2,3 ... => No. of seconds remaining
	*/

	if data, ok := kvStore.kv_stores[shardIdx][key]; ok {
		if utils.IsExpired(data.Expiry) {
			kvStore.mutex[shardIdx].Lock()
			defer kvStore.mutex[shardIdx].Unlock()
			delete(kvStore.kv_stores[shardIdx], key)
			return -2
		}
		if data.Expiry.IsZero() {
			return -1
		}
		return int(data.Expiry.Sub(time.Now()).Abs().Seconds())
	}
	return -2
}

func (kvStore *ShardedKvStore) updateWithOffset(key string, offset int, defaultValue string) (string, error) {
	shardIdx := kvStore.resolveShardIdx(key)

	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()

	newValue := defaultValue
	newExpiry := time.Time{}

	data, found := kvStore.kv_stores[shardIdx][key]

	if found {
		if utils.IsExpired(data.Expiry) {
			newValue = defaultValue
			newExpiry = time.Time{}
		} else {
			existingValueInt, err := strconv.Atoi(data.Value)
			if err != nil {
				errString := "ERR value is not an integer or out of range"
				return "", errors.New(errString)
			}
			newValue = strconv.FormatInt(int64(existingValueInt+offset), 10)
			newExpiry = data.Expiry
		}
	}

	kvStore.kv_stores[shardIdx][key] = StoredValue{
		Value:  newValue,
		Expiry: newExpiry,
	}
	return newValue, nil
}

func (kvStore *ShardedKvStore) Incr(key string) (string, error) {
	return kvStore.updateWithOffset(key, 1, "1")
}

func (kvStore *ShardedKvStore) Decr(key string) (string, error) {
	return kvStore.updateWithOffset(key, -1, "-1")
}
