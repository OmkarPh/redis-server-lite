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

type SetAction struct{}

func (action *SetAction) Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error) {
	if len(args) < 2 {
		errString := "ERR wrong number of arguments for 'SET' command"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	key := utils.ResolvePossibleKeyDirectives(args[0])
	value := args[1]
	slog.Debug(fmt.Sprintf("Set action (%s => %s)\n", key, value))

	optionsResolved, err := utils.ResolveSetOptions(args[2:]...)

	if err != nil {
		return [][]byte{resp.ResolveResponse(err.Error(), resp.Response_ERRORS)}, err
	}

	options := store.SetOptions(optionsResolved)

	if (optionsResolved.NX && optionsResolved.XX) || (optionsResolved.ExpireDuration && optionsResolved.ExpireTimestamp) || (optionsResolved.KEEPTTL && (optionsResolved.ExpireDuration || optionsResolved.ExpireTimestamp)) {
		errString := "ERR syntax error"
		return [][]byte{resp.ResolveResponse(errString, resp.Response_ERRORS)}, errors.New(errString)
	}

	set, prevValue, err := (*kvStore).Set(key, value, options)

	if err != nil {
		return [][]byte{resp.ResolveResponse(err.Error(), resp.Response_ERRORS)}, err
	}

	if set {
		if optionsResolved.GET {
			if prevValue == "" {
				return [][]byte{resp.ResolveResponse(nil, resp.Response_NULL)}, nil
			}
			return [][]byte{resp.ResolveResponse(prevValue, resp.Response_BULK_STRINGS)}, nil
		}
		return [][]byte{resp.ResolveResponse("OK", resp.Response_SIMPLE_STRING)}, nil
	}
	return [][]byte{resp.ResolveResponse(nil, resp.Response_NULL)}, nil
}
