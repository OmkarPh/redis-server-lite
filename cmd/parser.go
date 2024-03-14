package main

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

func ResolveRequestType(identifier byte) RequestType {
	requestTypes := map[byte]RequestType{
		'+': REQUEST_SIMPLE_STRING,
		'-': REQUEST_ERRORS,
		':': REQUEST_INTEGERS,
		'$': REQUEST_BULKSTRINGS,
		'*': REQUEST_ARRAYS,
		'_': REQUEST_NULL,
	}

	if requestType, ok := requestTypes[identifier]; ok {
		return requestType
	}
	return REQUEST_SIMPLE_STRING
}

func MarshallResponse(message string, responseType ResponseType) []byte {
	if responseResolver, ok := ResponseResolvers[responseType]; ok {
		response := responseResolver(message)
		// slog.Debug(fmt.Sprintf("Response:\n%s\n", string(response[:]))
		return response
	}
	return ResponseResolvers[Response_SIMPLE_STRING](message)
}

type Command struct {
	action ActionKey
	args   []string
}

func ResolveCommands(reader bufio.Reader, requestType RequestType) []Command {
	if requestType != REQUEST_ARRAYS {
		errString := "Request type not implemented"
		return []Command{
			{
				action: UNKNOWN,
				args:   []string{errString},
			},
		}
	}

	var commands []Command

	for {
		arrayIdentifier, err := reader.ReadByte()
		if err != nil {
			slog.Debug(fmt.Sprintf("Errr reading arrayidentifier: %b\n", arrayIdentifier))
			break
		}

		numOfPartsStr, _ := reader.ReadString('\n')
		// slog.Debug(fmt.Sprintf("Parts str: %s\n", numOfPartsStr)))
		numOfPartsStr = strings.Trim(numOfPartsStr, "\r\n")
		numOfParts, _ := strconv.Atoi(numOfPartsStr)
		// slog.Debug(fmt.Sprintf("%d parts to read\n", numOfParts)))

		if numOfParts == 0 {
			break
		}

		currentCommand := Command{
			args: []string{},
		}
		reader.ReadLine()
		actionString, _ := reader.ReadString('\n')
		actionString = strings.ToLower(strings.Trim(actionString, "\r\n"))

		actionArgsStartAtIdx := map[string]int{
			"config": 2,
		}

		argsStartAt := 1
		if idx, ok := actionArgsStartAtIdx[actionString]; ok {
			argsStartAt = idx
		}

		for i := 1; i < numOfParts; i++ {
			reader.ReadLine()
			arg, _ := reader.ReadString('\n')
			arg = strings.Trim(arg, "\r\n")
			slog.Debug(fmt.Sprintf("Arg %d: %s\n", i, arg))

			if i < argsStartAt {
				actionString = actionString + "_" + arg
			} else {
				currentCommand.args = append(currentCommand.args, arg)
			}
		}

		currentCommand.action = ActionKey(strings.ToLower(actionString))
		commands = append(commands, currentCommand)
		slog.Debug(fmt.Sprintf("Added command %v\n", currentCommand))
	}
	return commands
}
