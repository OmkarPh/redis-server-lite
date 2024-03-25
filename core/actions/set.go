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

type SetAction struct{}

func (action *SetAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 2 {
		errString := "ERR wrong number of arguments for 'SET' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])
	value := args[1]
	slog.Debug(fmt.Sprintf("Set action (%s => %s)\n", key, value))

	(*kvStore).Set(key, value)
	return [][]byte{resp.ResolveResponse("OK", resp.Response_SIMPLE_STRING)}, nil
}
