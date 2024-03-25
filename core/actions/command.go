package actions

import (
	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type CommandAction struct{}

func (action *CommandAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	return [][]byte{resp.ResolveResponse("Docs not available yet.", resp.Response_SIMPLE_STRING)}, nil
}
