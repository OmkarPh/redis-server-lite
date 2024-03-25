package actions

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
	"github.com/OmkarPh/redis-lite/utils"
)

type IncrAction struct{}

func (action *IncrAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'INCR' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])
	slog.Debug(fmt.Sprintf("Incr action (%s)\n", key))

	newValueString, err := (*kvStore).Incr(key)
	if err != nil {
		errString := "ERR value is not an integer or out of range"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	(*kvStore).Set(key, newValueString)
	return [][]byte{resp.ResolveResponse(newValueString, resp.Response_INTEGERS)}, nil
}
