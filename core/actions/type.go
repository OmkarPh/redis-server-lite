package actions

import (
	"errors"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type TypeAction struct{}

func (action *TypeAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'TYPE' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	// key := args[0]
	// existingValue, found := kvStore.get(key)

	// Note - There's no integer storage in redis :)
	// @TODO - Implement values of type sets, list, etc

	return [][]byte{resp.ResolveResponse("string", resp.Response_SIMPLE_STRING)}, nil
}
