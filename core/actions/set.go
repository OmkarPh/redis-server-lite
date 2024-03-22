package actions

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/engine"
	"github.com/OmkarPh/redis-lite/resp"
)

type SetAction struct{}

func (action *SetAction) Execute(kvEngine *engine.KvEngine, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 2 {
		errString := "ERR wrong number of arguments for 'set' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := args[0]
	if key == "key:__rand_int__" {
		key = strconv.FormatInt(int64(rand.Uint64()), 10)
	}

	value := args[1]

	slog.Debug(fmt.Sprintf("Set action (%s => %s)\n", key, value))
	(*kvEngine).Set(key, value)
	return [][]byte{resp.ResolveResponse("OK", resp.Response_SIMPLE_STRING)}, nil
}
