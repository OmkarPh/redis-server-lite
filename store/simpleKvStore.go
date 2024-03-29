package store

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/OmkarPh/redis-lite/utils"
)

type SimpleKvStore struct {
	mutex    sync.RWMutex // https://upstash.com/blog/upgradable-rwlock-for-go
	kv_store map[string]StoredValue
}

func NewSimpleStore() *SimpleKvStore {
	return &SimpleKvStore{
		kv_store: map[string]StoredValue{
			"dev": {
				Value:  "Omkar Phansopkar",
				Expiry: time.Time{},
			},
		},
	}
}

func (kvStore *SimpleKvStore) Get(key string) (string, bool) {
	if data, ok := kvStore.kv_store[key]; ok {
		if utils.IsExpired(data.Expiry) {
			kvStore.mutex.Lock()
			defer kvStore.mutex.Unlock()
			delete(kvStore.kv_store, key)
			return "", false
		}
		return data.Value, true
	}
	return "", false
}

func (kvStore *SimpleKvStore) Set(key string, value string) {
	kvStore.mutex.Lock()
	defer kvStore.mutex.Unlock()
	kvStore.kv_store[key] = StoredValue{
		Value:  value,
		Expiry: time.Time{},
	}
}

func (kvStore *SimpleKvStore) Del(key string) bool {
	kvStore.mutex.Lock()
	defer kvStore.mutex.Unlock()
	_, existed := kvStore.kv_store[key]
	delete(kvStore.kv_store, key)
	return existed
}

func (kvStore *SimpleKvStore) DeleteIfExpired(keysCount int) int {
	removedKeys := 0
	kvStore.mutex.Lock()
	defer kvStore.mutex.Unlock()
	for key, data := range kvStore.kv_store {
		if utils.IsExpired(data.Expiry) {
			// slog.Debug(fmt.Sprintf("Deleting expired key %s", key))
			delete(kvStore.kv_store, key)
			removedKeys++
		}
	}
	return removedKeys
}

func (kvStore *SimpleKvStore) Expire(key string, seconds int64, options ExpireOptions) (bool, error) {
	kvStore.mutex.Lock()
	defer kvStore.mutex.Unlock()

	newExpiry := time.Now().Add(time.Second * time.Duration(seconds))

	existingData, exists := kvStore.kv_store[key]
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

	kvStore.kv_store[key] = StoredValue{
		Value:  existingData.Value,
		Expiry: newExpiry,
	}

	return true, nil
}

func (kvStore *SimpleKvStore) Persist(key string) bool {
	kvStore.mutex.Lock()
	defer kvStore.mutex.Unlock()

	existingData, exists := kvStore.kv_store[key]
	if !exists {
		return false
	}

	kvStore.kv_store[key] = StoredValue{
		Value:  existingData.Value,
		Expiry: time.Time{},
	}

	return true
}

func (kvStore *SimpleKvStore) Ttl(key string) int {
	/*
		-2 => Key doesn't exist
		-1 => No expiry set
		0,1,2,3 ... => No. of seconds remaining
	*/

	if data, ok := kvStore.kv_store[key]; ok {
		if utils.IsExpired(data.Expiry) {
			kvStore.mutex.Lock()
			defer kvStore.mutex.Unlock()
			delete(kvStore.kv_store, key)
			return -2
		}
		if data.Expiry.IsZero() {
			return -1
		}
		return int(time.Until(data.Expiry).Seconds())
	}
	return -2
}

func (kvStore *SimpleKvStore) updateWithOffset(key string, offset int, defaultValue string) (string, error) {
	kvStore.mutex.Lock()
	defer kvStore.mutex.Unlock()

	newValue := defaultValue
	newExpiry := time.Time{}

	data, found := kvStore.kv_store[key]

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

	kvStore.kv_store[key] = StoredValue{
		Value:  newValue,
		Expiry: newExpiry,
	}
	return newValue, nil
}

func (kvStore *SimpleKvStore) Incr(key string) (string, error) {
	return kvStore.updateWithOffset(key, 1, "1")
}

func (kvStore *SimpleKvStore) Decr(key string) (string, error) {
	return kvStore.updateWithOffset(key, -1, "-1")
}
