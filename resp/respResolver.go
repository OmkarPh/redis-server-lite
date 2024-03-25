package resp

import (
	"reflect"
	"strconv"
	"strings"
)

type ResponseType int

const (
	Response_SIMPLE_STRING ResponseType = iota
	Response_BULK_STRINGS
	Response_NULL
	Response_NULL_BULK_STRING
	Response_ERRORS
	Response_INTEGERS
	Response_ARRAY
)

func simpleStringResponseResolver(response interface{}) []byte {
	// eg. some_value => ```+some_value\r\n```

	responseType := reflect.TypeOf(response).String()
	responseStr := ""

	switch responseType {
	case "string":
		responseStr = response.(string)
	case "int":
		responseStr = strconv.FormatInt(response.(int64), 10)
	}

	responseBytes := append([]byte{'+'}, []byte(responseStr)...)
	return append(responseBytes, '\r', '\n')
}

func bulkStringsResponseResolver(response interface{}) []byte {
	// eg. some_value => ```$10\r\nsome_value\r\n```

	responseType := reflect.TypeOf(response).String()
	responseStr := ""

	switch responseType {
	case "string":
		responseStr = "$" + strconv.FormatInt(int64(len(response.(string))), 10) + "\r\n" + response.(string)
	case "int":
		str := strconv.FormatInt(response.(int64), 10)
		responseStr = "$" + strconv.FormatInt(int64(len(str)), 10) + "\r\n" + str
	case "[]string":
		responseItems := []string{}
		for _, responseItem := range response.([]string) {
			responseItems = append(responseItems, "$"+strconv.FormatInt(int64(len(responseItem)), 10))
			responseItems = append(responseItems, responseItem)
		}
		responseStr = "*" + strconv.FormatInt(int64(len(response.([]string))), 10) + "\r\n"
		responseStr += strings.Join(responseItems, "\r\n")
	}

	responseBytes := []byte(responseStr)
	return append(responseBytes, '\r', '\n')
}

func nullResponseResolver(_ interface{}) []byte {
	// eg. ```$_\r\n```
	return append([]byte("_"), '\r', '\n')
}

func nullBulkStringResponseResolver(_ interface{}) []byte {
	// eg. ```$-1\r\n```
	return append([]byte("$-1"), '\r', '\n')
}

func errResponseResolver(response interface{}) []byte {
	// eg. -error_message\r\n

	responseType := reflect.TypeOf(response).String()
	responseStr := ""

	switch responseType {
	case "string":
		responseStr = response.(string)
	case "int":
		responseStr = strconv.FormatInt(response.(int64), 10)
	}

	responseBytes := append([]byte("-"), []byte(responseStr)...)
	return append(responseBytes, '\r', '\n')
}

func integerResponseResolver(response interface{}) []byte {
	// eg. 4 => ```:4\r\n```
	responseType := reflect.TypeOf(response).String()
	responseStr := ""

	switch responseType {
	case "string":
		responseStr = response.(string)
	case "int":
		responseStr = strconv.FormatInt(int64(response.(int)), 10)
	case "uint":
		responseStr = strconv.FormatInt(int64(response.(uint)), 10)
	case "int64":
		responseStr = strconv.FormatInt(int64(response.(int64)), 10)
	}

	return append([]byte(":"+responseStr), '\r', '\n')
}

var RespResponseResolvers = map[ResponseType]func(interface{}) []byte{
	Response_SIMPLE_STRING:    simpleStringResponseResolver,
	Response_BULK_STRINGS:     bulkStringsResponseResolver,
	Response_NULL:             nullResponseResolver,
	Response_NULL_BULK_STRING: nullBulkStringResponseResolver,
	Response_ERRORS:           errResponseResolver,
	Response_INTEGERS:         integerResponseResolver,
}

func ResolveResponse(message interface{}, responseType ResponseType) []byte {
	if responseResolver, ok := RespResponseResolvers[responseType]; ok {
		response := responseResolver(message)
		// slog.Debug(fmt.Sprintf("Response:\n%s\n", string(response[:]))
		return response
	}
	return RespResponseResolvers[Response_SIMPLE_STRING](message)
}
