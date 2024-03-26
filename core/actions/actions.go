package actions

import (
	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/resp"
	"github.com/OmkarPh/redis-lite/store"
)

type Action interface {
	Execute(kvStore *store.KvStore, redisConfig *config.RedisConfig, args ...string) ([][]byte, error)
}

var RedisActions = map[resp.ActionKey]Action{
	resp.ACTION_PING:    &PingAction{},
	resp.ACTION_ECHO:    &EchoAction{},
	resp.ACTION_GET:     &GetAction{},
	resp.ACTION_SET:     &SetAction{},
	resp.ACTION_DEL:     &DelAction{},
	resp.ACTION_INCR:    &IncrAction{},
	resp.ACTION_DECR:    &DecrAction{},
	resp.ACTION_TYPE:    &TypeAction{},
	resp.ACTION_UNKNOWN: &UnknownAction{},
	resp.ACTION_CONFIG:  &ConfigAction{},
	resp.ACTION_COMMAND: &CommandAction{},
	resp.ACTION_EXISTS:  &ExistsAction{},
	resp.ACTION_TTL:     &TtlAction{},
	resp.ACTION_EXPIRE:  &ExpireAction{},
	resp.ACTION_PERSIST: &PersistAction{},
}
