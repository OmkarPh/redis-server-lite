package actions

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/engine"
	"github.com/OmkarPh/redis-lite/resp"
)

type IncrAction struct{}

func (action *IncrAction) Execute(kvEngine *engine.KvEngine, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'incr' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := args[0]
	slog.Debug(fmt.Sprintf("Incr action (%s)\n", key))

	newValueString, err := (*kvEngine).Incr(key)
	if err != nil {
		errString := "ERR value is not an integer or out of range"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	(*kvEngine).Set(key, newValueString)
	return [][]byte{resp.ResolveResponse(newValueString, resp.Response_INTEGERS)}, nil
}
