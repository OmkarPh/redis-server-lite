package main

import (
	"strconv"
)

type ResponseType int

const (
	Response_SIMPLE_STRING ResponseType = iota
	Response_BULK_STRING
	Response_NULL
	Response_ERRORS
	Response_INTEGERS
	// Response_ARRAYS
)

func simpleStringResponseResolver(response string) []byte {
	responseBytes := []byte{}
	responseBytes = append(responseBytes, '+')
	responseBytes = append(responseBytes, []byte(response)...)
	return append(responseBytes, '\r', '\n')
}

func bulkStringResponseResolver(response string) []byte {
	responseBytes := []byte{}
	responseBytes = append(responseBytes, []byte("$"+strconv.FormatInt(int64(len(response)), 10))...)
	responseBytes = append(responseBytes, '\r', '\n')
	responseBytes = append(responseBytes, []byte(response)...)
	return append(responseBytes, '\r', '\n')
}

func nullResponseResolver(response string) []byte {
	return append([]byte("$-1"), '\r', '\n')
}

func errResponseResolver(response string) []byte {
	responseBytes := []byte{}
	responseBytes = append(responseBytes, '-')
	responseBytes = append(responseBytes, []byte(response)...)
	return append(responseBytes, '\r', '\n')
}

func integerResponseResolver(response string) []byte {
	return append([]byte(":"+response), '\r', '\n')
}

var ResponseResolvers = map[ResponseType]func(string) []byte{
	Response_SIMPLE_STRING: simpleStringResponseResolver,
	Response_BULK_STRING:   bulkStringResponseResolver,
	Response_NULL:          nullResponseResolver,
	Response_ERRORS:        errResponseResolver,
	Response_INTEGERS:      integerResponseResolver,
}
