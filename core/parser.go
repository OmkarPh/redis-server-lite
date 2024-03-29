package core

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/OmkarPh/redis-lite/resp"
)

func ResolveRequestType(identifier byte) resp.RequestType {
	requestTypes := map[byte]resp.RequestType{
		'+': resp.REQUEST_SIMPLE_STRING,
		'-': resp.REQUEST_ERRORS,
		':': resp.REQUEST_INTEGERS,
		'$': resp.REQUEST_BULKSTRINGS,
		'*': resp.REQUEST_ARRAYS,
		'_': resp.REQUEST_NULL,
	}
	slog.Debug(fmt.Sprintf("Resolve %s", string(identifier)))
	if requestType, ok := requestTypes[identifier]; ok {
		return requestType
	}
	return resp.REQUEST_SIMPLE_STRING
}

func ResolveCommands(reader bufio.Reader, requestType resp.RequestType) []resp.Command {
	slog.Debug(fmt.Sprint("Serve request =>", requestType))
	parser, parserImplemented := resp.RespParsers[requestType]
	if !parserImplemented {
		errString := "Request type not implemented"
		return []resp.Command{
			{
				Action: resp.ACTION_UNKNOWN,
				Args:   []string{errString},
			},
		}
	}

	return parser(reader)
}
