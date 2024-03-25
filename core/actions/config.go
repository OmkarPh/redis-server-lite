package actions

import (
	"fmt"
	"strings"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type ConfigAction struct{}

func (action *ConfigAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) != 2 {
		return [][]byte{resp.ResolveResponse("Not implemented", resp.Response_ERRORS)}, nil
	}

	subAction := strings.ToLower(args[0])
	configKey := args[1]
	response := []byte{}

	switch subAction {
	case "get":
		configValue, found := (*config.RedisConfig).GetParam(redisConfig, configKey)
		if !found {
			return [][]byte{resp.ResolveResponse(fmt.Sprintf("Config key %s not available", configKey), resp.Response_ERRORS)}, nil
		}
		response = resp.ResolveResponse([]string{configKey, configValue}, resp.Response_BULK_STRINGS)
	case "set":
		// @TODO
		response = resp.ResolveResponse("OK", resp.Response_SIMPLE_STRING)
	}

	return [][]byte{response}, nil
}
