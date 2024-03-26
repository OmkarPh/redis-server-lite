package resp

import "strings"

type ActionKey string

const (
	ACTION_PING    ActionKey = "ping"
	ACTION_ECHO    ActionKey = "echo"
	ACTION_GET     ActionKey = "get"
	ACTION_SET     ActionKey = "set"
	ACTION_DEL     ActionKey = "del"
	ACTION_INCR    ActionKey = "incr"
	ACTION_DECR    ActionKey = "decr"
	ACTION_TYPE    ActionKey = "type"
	ACTION_UNKNOWN ActionKey = "unknown"
	ACTION_CONFIG  ActionKey = "config"
	ACTION_COMMAND ActionKey = "command"
	ACTION_EXISTS  ActionKey = "exists"
	ACTION_EXPIRE  ActionKey = "expire"
	ACTION_PERSIST ActionKey = "persist"
	ACTION_TTL     ActionKey = "ttl"
)

type Command struct {
	Action ActionKey
	Args   []string
}

func tokensToCommand(tokens []string) Command {
	return Command{
		Action: ActionKey(strings.ToLower(tokens[0])),
		Args:   tokens[1:],
	}
}
