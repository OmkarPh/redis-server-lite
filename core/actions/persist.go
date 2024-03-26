package actions

import (
	"errors"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
	"github.com/OmkarPh/redis-lite/utils"
)

type PersistAction struct{}

func (action *PersistAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'PERSIST' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])

	success := (*kvStore).Persist(key)

	if success {
		return [][]byte{resp.ResolveResponse(1, resp.Response_INTEGERS)}, nil
	} else {
		return [][]byte{resp.ResolveResponse(0, resp.Response_INTEGERS)}, nil
	}
}
