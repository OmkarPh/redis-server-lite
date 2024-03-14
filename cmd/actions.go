package main

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
)

type ActionKey string

const (
	PING       ActionKey = "ping"
	ECHO       ActionKey = "echo"
	GET        ActionKey = "get"
	SET        ActionKey = "set"
	INCR       ActionKey = "incr"
	TYPE       ActionKey = "type"
	UNKNOWN    ActionKey = "unknown"
	CONFIG_GET ActionKey = "config_get"
	CONFIG_SET ActionKey = "config_set"
)

type Action interface {
	Execute(args ...string) ([]byte, error)
}

type PingAction struct{}

func (action *PingAction) Execute(args ...string) ([]byte, error) {
	slog.Debug("Ping action")
	if len(args) > 1 {
		errString := "ERR wrong number of arguments for 'ping' command"
		return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
	}
	if len(args) == 1 {
		return RedisActions[ECHO].Execute(args...)
	}
	return MarshallResponse("PONG", Response_SIMPLE_STRING), nil
}

type EchoAction struct{}

func (action *EchoAction) Execute(args ...string) ([]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'echo' command"
		return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
	}

	slog.Debug("Echo:", args[0])
	return MarshallResponse(args[0], Response_SIMPLE_STRING), nil
}

type GetAction struct{}

func (action *GetAction) Execute(args ...string) ([]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'get' command"
		return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
	}

	key := args[0]
	slog.Debug(fmt.Sprintf("Get action (%s)\n", key))
	value, found := KvEngine.get(key)
	if found {
		return MarshallResponse(value, Response_BULK_STRING), nil
	}
	return MarshallResponse("", Response_NULL), nil
}

type SetAction struct{}

func (action *SetAction) Execute(args ...string) ([]byte, error) {
	if len(args) != 2 {
		errString := "ERR wrong number of arguments for 'set' command"
		return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
	}

	key := args[0]
	value := args[1]

	slog.Debug(fmt.Sprintf("Set action (%s => %s)\n", key, value))
	KvEngine.set(key, value)
	return MarshallResponse("OK", Response_SIMPLE_STRING), nil
}

type TypeAction struct{}

func (action *TypeAction) Execute(args ...string) ([]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'incr' command"
		return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
	}

	// key := args[0]
	// existingValue, found := KvEngine.get(key)

	// Note - There's no integer storage in redis :)
	// @TODO - Implement values of type sets, list, etc

	return MarshallResponse("string", Response_SIMPLE_STRING), nil
}

type IncrAction struct{}

func (action *IncrAction) Execute(args ...string) ([]byte, error) {
	if len(args) != 1 {
		errString := "ERR wrong number of arguments for 'incr' command"
		return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
	}

	key := args[0]
	newValue := 1

	slog.Debug(fmt.Sprintf("Incr action (%s)\n", key))
	existingValue, found := KvEngine.get(key)

	if found {
		existingValueInt, err := strconv.Atoi(existingValue)
		if err != nil {
			errString := "ERR value is not an integer or out of range"
			return MarshallResponse(errString, Response_ERRORS), errors.New(errString)
		}
		newValue = existingValueInt + 1
	}

	newValueString := strconv.FormatInt(int64(newValue), 10)
	KvEngine.set(key, newValueString)
	return MarshallResponse(newValueString, Response_INTEGERS), nil
}

type ConfigAction struct{}

func (action *ConfigAction) Execute(args ...string) ([]byte, error) {
	if len(args) != 1 {
		return MarshallResponse("Not implemented", Response_ERRORS), nil
	}

	configKey := args[0]
	configValue, found := redisConfig.getParam(configKey)

	if !found {
		return MarshallResponse(fmt.Sprintf("Config key %s not available", configKey), Response_ERRORS), nil
	}

	marshalled := []byte("*2\r\n")
	marshalled = append(marshalled, MarshallResponse(configKey, Response_BULK_STRING)...)
	marshalled = append(marshalled, MarshallResponse(configValue, Response_BULK_STRING)...)
	return marshalled, nil
}

type UnknownAction struct{}

func (action *UnknownAction) Execute(args ...string) ([]byte, error) {
	errString := "Request type not implemented !"
	if len(args) > 0 {
		errString = args[0]
	}
	return MarshallResponse(errString, Response_ERRORS), nil
}

var RedisActions = map[ActionKey]Action{
	PING:       &PingAction{},
	ECHO:       &EchoAction{},
	GET:        &GetAction{},
	SET:        &SetAction{},
	INCR:       &IncrAction{},
	TYPE:       &TypeAction{},
	UNKNOWN:    &UnknownAction{},
	CONFIG_GET: &ConfigAction{},
}
