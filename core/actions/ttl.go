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

type TtlAction struct{}

func (action *TtlAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'TTL' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])
	slog.Debug(fmt.Sprintf("TTL action (%s)\n", key))

	seconds := (*kvStore).Ttl(key)

	return [][]byte{resp.ResolveResponse(seconds, resp.Response_INTEGERS)}, nil
}
