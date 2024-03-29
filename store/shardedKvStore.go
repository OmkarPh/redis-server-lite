package store

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/OmkarPh/redis-lite/utils"
	"github.com/spaolacci/murmur3"
)

type ShardedKvStore struct {
	shardFactor uint32
	mutex       []sync.RWMutex
	kv_stores   []map[string]StoredValue
}

func NewShardedKvStore(shardFactor uint32) *ShardedKvStore {
	newKvStore := &ShardedKvStore{
		shardFactor: shardFactor,
		mutex:       make([]sync.RWMutex, shardFactor),
		kv_stores:   make([]map[string]StoredValue, shardFactor),
	}
	for i := uint32(0); i < shardFactor; i++ {
		newKvStore.kv_stores[i] = make(map[string]StoredValue)
	}
	fmt.Printf("%d Shards created\n", shardFactor)
	return newKvStore
}

// var shardIdxCache map[string]uint64 = make(map[string]uint64)
func (kvStore *ShardedKvStore) resolveShardIdx(key string) uint32 {
	// if shardIdx, ok := shardIdxCache[key]; ok {
	// 	return shardIdx
	// }

	// hasher := xxhash.New()
	// hasher.Write([]byte(key))
	// shardIdx := uint32(hasher.Sum64() % uint64(kvStore.shardFactor))

	h := murmur3.New32()
	h.Write([]byte(key))

	var shardIdx uint32
	if kvStore.shardFactor&(kvStore.shardFactor-1) == 0 {
		shardIdx = h.Sum32() & (kvStore.shardFactor - 1) // Faster than modulus (when shardfactor is power of 2)
	} else {
		shardIdx = h.Sum32() % kvStore.shardFactor
	}

	slog.Debug(fmt.Sprintf("Shard idx for %s => %d", key, shardIdx))
	// shardIdxCache[key] = shardIdx
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

func (kvStore *ShardedKvStore) Set(key string, value string, options SetOptions) (bool, string, error) {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()

	existingData, exists := kvStore.kv_stores[shardIdx][key]

	// Convert options to ExpirationTimeOptions
	expirationOptions := utils.ExpirationTimeOptions{
		NX:                        options.NX,
		XX:                        options.XX,
		ExpireDuration:            options.ExpireDuration,
		ExpiryTimeSeconds:         options.ExpiryTimeSeconds,
		ExpiryTimeMiliSeconds:     options.ExpiryTimeMiliSeconds,
		ExpireTimestamp:           options.ExpireTimestamp,
		ExpiryUnixTimeSeconds:     options.ExpiryUnixTimeSeconds,
		ExpiryUnixTimeMiliSeconds: options.ExpiryUnixTimeMiliSeconds,
		KEEPTTL:                   options.KEEPTTL,
	}

	if (options.NX && exists) || (options.XX && !exists) {
		return false, "", nil
	}

	expiryTime, canSet, err := utils.ResolveExpirationTime(expirationOptions, exists, existingData.Expiry)

	if !canSet {
		return false, "", err
	}

	kvStore.kv_stores[shardIdx][key] = StoredValue{
		Value:  value,
		Expiry: expiryTime,
	}
	return true, existingData.Value, nil
}

func (kvStore *ShardedKvStore) Del(key string) bool {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()
	_, existed := kvStore.kv_stores[shardIdx][key]
	delete(kvStore.kv_stores[shardIdx], key)
	return existed
}

func (kvStore *ShardedKvStore) DeleteIfExpired(keysCount int) int {
	removedKeys := 0
	for shardIdx := uint32(0); shardIdx < kvStore.shardFactor; shardIdx++ {
		kvStore.mutex[shardIdx].Lock()
		for key, data := range kvStore.kv_stores[shardIdx] {
			if utils.IsExpired(data.Expiry) {
				// slog.Debug(fmt.Sprintf("Deleting expired key %s", key))
				delete(kvStore.kv_stores[shardIdx], key)
				removedKeys++
			}
		}

		kvStore.mutex[shardIdx].Unlock()
	}
	return removedKeys
}

func (kvStore *ShardedKvStore) Expire(key string, seconds int64, options ExpireOptions) (bool, error) {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()

	newExpiry := time.Now().Add(time.Second * time.Duration(seconds))

	existingData, exists := kvStore.kv_stores[shardIdx][key]
	if !exists {
		return false, errors.New("ERR key doesn't exist")
	}

	if options.NX && !existingData.Expiry.IsZero() {
		return false, nil
	}
	if options.XX && existingData.Expiry.IsZero() {
		return false, nil
	}
	if options.LT && !existingData.Expiry.IsZero() && newExpiry.After(existingData.Expiry) {
		return false, nil
	}
	if options.GT && !existingData.Expiry.IsZero() && newExpiry.Before(existingData.Expiry) {
		return false, nil
	}

	kvStore.kv_stores[shardIdx][key] = StoredValue{
		Value:  existingData.Value,
		Expiry: newExpiry,
	}

	return true, nil
}

func (kvStore *ShardedKvStore) Persist(key string) bool {
	shardIdx := kvStore.resolveShardIdx(key)
	kvStore.mutex[shardIdx].Lock()
	defer kvStore.mutex[shardIdx].Unlock()

	existingData, exists := kvStore.kv_stores[shardIdx][key]
	if !exists || existingData.Expiry.IsZero() {
		return false
	}

	kvStore.kv_stores[shardIdx][key] = StoredValue{
		Value:  existingData.Value,
		Expiry: time.Time{},
	}

	return true
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
		return int(time.Until(data.Expiry).Seconds())
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
