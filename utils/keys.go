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
	normalisedKey := strings.ToLower(key)
	if normalisedKey == "key:__rand_int__" || normalisedKey == "__rand_int__" {
		return GenerateRandomKey()
	}
	return key
}
