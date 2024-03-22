package actions

import (
	"errors"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/engine"
	"github.com/OmkarPh/redis-lite/resp"
)

type TypeAction struct{}

func (action *TypeAction) Execute(kvEngine *engine.KvEngine, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'incr' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	// key := args[0]
	// existingValue, found := KvEngine.get(key)

	// Note - There's no integer storage in redis :)
	// @TODO - Implement values of type sets, list, etc

	return [][]byte{resp.ResolveResponse("string", resp.Response_SIMPLE_STRING)}, nil
}
