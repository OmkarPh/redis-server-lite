package actions

import (
	"errors"
	"strconv"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
	"github.com/OmkarPh/redis-lite/utils"
)

type ExpireAction struct{}

func (action *ExpireAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) < 2 || len(args) > 3 {
		errString := "ERR wrong number of arguments for 'EXPIRE' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])
	seconds, err := strconv.ParseInt(args[1], 10, 64)

	// extraOption := ""
	// if len(args) > 1{
	// 	extraOption := args[1]
	// }

	if err != nil {
		errString := "ERR invalid seconds arguments for 'EXPIRE' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	success, _ := (*kvStore).Expire(key, seconds)

	if success {
		return [][]byte{resp.ResolveResponse(1, resp.Response_INTEGERS)}, nil
	} else {
		return [][]byte{resp.ResolveResponse(0, resp.Response_INTEGERS)}, nil
	}
}
