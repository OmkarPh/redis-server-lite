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

type GetAction struct{}

func (action *GetAction) Execute(kvEngine *engine.KvEngine, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'get' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := args[0]
	if key == "key:__rand_int__" {
		key = strconv.FormatInt(int64(rand.Uint64()), 10)
	}

	slog.Debug(fmt.Sprintf("Get action (%s)\n", key))
	value, found := (*kvEngine).Get(key)
	if found {
		return [][]byte{resp.ResolveResponse(value, resp.Response_BULK_STRINGS)}, nil
	}
	return [][]byte{resp.ResolveResponse("", resp.Response_NULL)}, nil
}
