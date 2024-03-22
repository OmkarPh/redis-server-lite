package engine

import (
	"fmt"
	"strconv"

	"github.com/OmkarPh/redis-lite/config"
)

type EngineType string

const (
	SIMPLE_ENGINE  EngineType = "simple"
	SHARDED_ENGINE EngineType = "sharded"
)

type KvEngine interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Incr(key string) (string, error)
	Decr(key string) (string, error)
}

var EngineGenerators = map[EngineType]func(rc *config.RedisConfig) *KvEngine{
	SIMPLE_ENGINE: func(rc *config.RedisConfig) *KvEngine {
		var newEngine KvEngine = NewSimpleEngine()
		return &newEngine
	},
	SHARDED_ENGINE: func(rc *config.RedisConfig) *KvEngine {
		shardFactorStr, found := rc.GetParam("shardfactor")
		shardFactor := 10 // Default
		if found {
			var err error
			shardFactor, err = strconv.Atoi(shardFactorStr)
			if err != nil {
				shardFactor = 10 // Default
			}
		}
		var newEngine KvEngine = NewShardedEngine(shardFactor)
		return &newEngine
	},
}

func NewKvEngine(rc *config.RedisConfig) *KvEngine {
	engine, found := rc.GetParam("kv_engine")
	engineType := EngineType(engine)
	if !found {
		engineType = SIMPLE_ENGINE
	}
	fmt.Println("Using Engine type:", engineType)
	return EngineGenerators[engineType](rc)
}
