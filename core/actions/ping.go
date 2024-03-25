package actions

import (
	"errors"
	"log/slog"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type PingAction struct{}

func (action *PingAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	slog.Debug("Ping action")
	if len(args) > 1 {
		errString := "ERR wrong number of arguments for 'PING' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}
	if len(args) == 1 {
		return RedisActions[resp.ACTION_ECHO].Execute(kvStore, redisConfig, args...)
	}
	return [][]byte{resp.ResolveResponse("PONG", resp.Response_SIMPLE_STRING)}, nil
}
