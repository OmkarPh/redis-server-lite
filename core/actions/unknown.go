package actions

import (
	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type UnknownAction struct{}

func (action *UnknownAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	errString := "Request type not implemented !"
	if len(args) > 0 {
		errString = args[0]
	}
	return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, nil
}
