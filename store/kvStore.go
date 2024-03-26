package store

import (
	"fmt"
	"strconv"
	"time"

	"github.com/OmkarPh/redis-lite/config"
)

type StoreType string

const (
	SIMPLE_STORE  StoreType = "simple"
	SHARDED_STORE StoreType = "sharded"
)

// References
// https://redis.io/commands/expire/#how-expires-are-handled-in-the-replication-link-and-aof-file

type StoredValue struct {
	Value string
	// Value interface{}
	Expiry time.Time
}
type ExpireOptions struct {
	NX bool
	XX bool
	GT bool
	LT bool
}
type KvStore interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Del(key string) bool
	Incr(key string) (string, error)
	Decr(key string) (string, error)
	Expire(key string, seconds int64, options ExpireOptions) (bool, error)
	Persist(key string) bool
	Ttl(key string) int
}

var StoreGenerators = map[StoreType]func(rc *config.RedisConfig) *KvStore{
	SIMPLE_STORE: func(rc *config.RedisConfig) *KvStore {
		var newKvStore KvStore = NewSimpleStore()
		return &newKvStore
	},
	SHARDED_STORE: func(rc *config.RedisConfig) *KvStore {
		shardFactorStr, found := rc.GetParam("shardfactor")
		shardFactor := 10 // Default
		if found {
			var err error
			shardFactor, err = strconv.Atoi(shardFactorStr)
			if err != nil {
				shardFactor = 10 // Default
			}
		}
		var newKvStore KvStore = NewShardedKvStore(shardFactor)
		return &newKvStore
	},
}

func NewKvStore(rc *config.RedisConfig) *KvStore {
	kv_store_config, found := rc.GetParam("kv_store")
	kvStoreType := StoreType(kv_store_config)
	if !found {
		kvStoreType = SIMPLE_STORE
	}
	fmt.Println("Using Store type:", kvStoreType)
	return StoreGenerators[kvStoreType](rc)
}
