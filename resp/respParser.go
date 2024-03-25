package resp

import (
	"bufio"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type RequestType int

const (
	REQUEST_SIMPLE_STRING RequestType = iota
	REQUEST_ERRORS
	REQUEST_INTEGERS
	REQUEST_BULKSTRINGS
	REQUEST_ARRAYS
	REQUEST_NULL
)

func resolveSimpleString(reader bufio.Reader) []Command {
	simpleStringIdentifier, _ := reader.Peek(1)
	hasIdentifierPrefix := simpleStringIdentifier[0] == byte('+')
	if hasIdentifierPrefix {
		reader.ReadByte()
	}

	tokenString, err := reader.ReadString('\n')
	if err != nil {
		return []Command{
			{
				Action: ACTION_UNKNOWN,
				Args:   []string{"Couldn't parse your request"},
			},
		}
	}

	tokenString = strings.Trim(tokenString, "\r\n")
	return []Command{tokensToCommand(strings.Split(tokenString, " "))}
}

func resolveArrayRequest(reader bufio.Reader) []Command {
	slog.Debug("Serve array request")

	var commands []Command

	for {
		arrayIdentifier, err := reader.ReadByte()
		if err != nil {
			slog.Debug(fmt.Sprintf("Errr reading arrayidentifier: %b\n", arrayIdentifier))
			break
		}

		numOfTokensStr, _ := reader.ReadString('\n')
		numOfTokensStr = strings.Trim(numOfTokensStr, "\r\n")
		numOfTokens, _ := strconv.Atoi(numOfTokensStr)

		if numOfTokens == 0 {
			break
		}

		tokens := []string{}
		for i := 0; i < numOfTokens; i++ {
			reader.ReadLine() // Token length is not needed
			token, _ := reader.ReadString('\n')
			token = strings.Trim(token, "\r\n")
			tokens = append(tokens, token)
			// slog.Debug(fmt.Sprintf("Token[%d]: %s\n", i, token))
		}
		commands = append(commands, tokensToCommand(tokens))
		// slog.Debug(fmt.Sprintf("Command %v\n", commands[len(commands)-1]))
	}
	slog.Debug(fmt.Sprintf("Resolved commands %v\n", commands))
	return commands
}

var RespParsers = map[RequestType]func(reader bufio.Reader) []Command{
	REQUEST_SIMPLE_STRING: resolveSimpleString,
	REQUEST_ARRAYS:        resolveArrayRequest,
}
