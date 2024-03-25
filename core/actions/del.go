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

type DelAction struct{}

func (action *DelAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) < 1 {
		errString := "ERR wrong number of arguments for 'DEL' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	keys := utils.MapOver(args, utils.ResolvePossibleKeyDirectives)

	existedKeys := 0
	for _, key := range keys {
		existed := (*kvStore).Del(key)
		slog.Debug(fmt.Sprintf("Exists action (%s) => %t\n", key, existed))
		if existed {
			existedKeys++
		}
	}

	return [][]byte{resp.ResolveResponse(existedKeys, resp.Response_INTEGERS)}, nil
}
