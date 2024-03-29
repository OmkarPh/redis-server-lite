package utils

import (
	"math/rand"
	"strconv"
	"strings"
)

func GenerateRandomKey() string {
	return strconv.FormatInt(int64(rand.Uint64()), 10)
}

func ResolvePossibleKeyDirectives(key string) string {
	key = strings.ToLower(key)
	if key == "key:__rand_int__" || key == "__rand_int__" {
		return GenerateRandomKey()
	}
	return key
}
