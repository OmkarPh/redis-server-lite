package actions

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/engine"
	"github.com/OmkarPh/redis-lite/resp"
)

type EchoAction struct{}

func (action *EchoAction) Execute(kvEngine *engine.KvEngine, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'echo' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	slog.Debug(fmt.Sprint("Echo:", args[0]))
	return [][]byte{resp.ResolveResponse(args[0], resp.Response_SIMPLE_STRING)}, nil
}
