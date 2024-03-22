package actions

import (
	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/engine"
	"github.com/OmkarPh/redis-lite/resp"
)

type CommandAction struct{}

func (action *CommandAction) Execute(kvEngine *engine.KvEngine, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	return [][]byte{resp.ResolveResponse("Docs not available yet.", resp.Response_SIMPLE_STRING)}, nil
}
