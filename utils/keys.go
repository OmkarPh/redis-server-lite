package utils

import (
	"math/rand"
	"strconv"
	"strings"
)

func ResolvePossibleKeyDirectives(key string) string {
	key = strings.ToLower(key)
	if key == "key:__rand_int__" || key == "__rand_int__" {
		return strconv.FormatInt(int64(rand.Uint64()), 10)
	}
	return key
}
