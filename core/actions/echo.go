package actions

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type EchoAction struct{}

func (action *EchoAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'ECHO' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	slog.Debug(fmt.Sprint("Echo:", args[0]))
	return [][]byte{resp.ResolveResponse(args[0], resp.Response_SIMPLE_STRING)}, nil
}
