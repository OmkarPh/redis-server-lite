package engine

import (
	"errors"
	"strconv"
	"sync"
)

type SimpleEngine struct {
	mutex    sync.Mutex
	kv_store map[string]string
	// kv_store map[string]interface{}
}

func NewSimpleEngine() *SimpleEngine {
	return &SimpleEngine{
		kv_store: map[string]string{
			"dev": "Omkar Phansopkar",
		},
	}
}

func (engine *SimpleEngine) Get(key string) (string, bool) {
	if value, ok := engine.kv_store[key]; ok {
		return value, true
	}
	return "", false
}

func (engine *SimpleEngine) Set(key string, value string) {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()
	engine.kv_store[key] = value
}

func (engine *SimpleEngine) Incr(key string) (string, error) {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()

	newValue := 1

	existingValue, found := engine.kv_store[key]

	if found {
		existingValueInt, err := strconv.Atoi(existingValue)
		if err != nil {
			errString := "ERR value is not an integer or out of range"
			return "", errors.New(errString)
		}
		newValue = existingValueInt + 1
	}

	return strconv.FormatInt(int64(newValue), 10), nil
}

func (engine *SimpleEngine) Decr(key string) (string, error) {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()

	newValue := -1

	existingValue, found := engine.kv_store[key]

	if found {
		existingValueInt, err := strconv.Atoi(existingValue)
		if err != nil {
			errString := "ERR value is not an integer or out of range"
			return "", errors.New(errString)
		}
		newValue = existingValueInt - 1
	}

	return strconv.FormatInt(int64(newValue), 10), nil
}
