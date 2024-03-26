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
	if len(args) < 2 {
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

	expireOptions := store.ExpireOptions{
		NX: false,
		XX: false,
		GT: false,
		LT: false,
	}

	if len(args) > 2 {
		for _, option := range args[2:] {
			switch option {
			case "NX":
				expireOptions.NX = true
			case "XX":
				expireOptions.XX = true
			case "GT":
				expireOptions.GT = true
			case "LT":
				expireOptions.LT = true
			}
		}
	}

	success, _ := (*kvStore).Expire(key, seconds, expireOptions)

	if success {
		return [][]byte{resp.ResolveResponse(1, resp.Response_INTEGERS)}, nil
	} else {
		return [][]byte{resp.ResolveResponse(0, resp.Response_INTEGERS)}, nil
	}
}
