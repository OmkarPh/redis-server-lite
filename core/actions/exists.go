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

type ExistsAction struct{}

func (action *ExistsAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) < 1 {
		errString := "ERR wrong number of arguments for 'EXISTS' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	keys := utils.MapOver(args, utils.ResolvePossibleKeyDirectives)

	existingKeys := 0
	for _, key := range keys {
		_, exists := (*kvStore).Get(key)
		slog.Debug(fmt.Sprintf("Exists action (%s) => %t\n", key, exists))
		if exists {
			existingKeys++
		}
	}

	return [][]byte{resp.ResolveResponse(existingKeys, resp.Response_INTEGERS)}, nil
}
