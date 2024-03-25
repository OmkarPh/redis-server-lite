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

type GetAction struct{}

func (action *GetAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'GET' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])
	slog.Debug(fmt.Sprintf("Get action (%s)\n", key))

	value, found := (*kvStore).Get(key)
	if found {
		return [][]byte{resp.ResolveResponse(value, resp.Response_BULK_STRINGS)}, nil
	}
	return [][]byte{resp.ResolveResponse("", resp.Response_NULL)}, nil
}
