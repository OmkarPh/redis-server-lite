package utils

import (
	"errors"
	"strconv"
	"strings"
)

type SetOptions struct {
	NX                        bool
	XX                        bool
	GET                       bool
	ExpireDuration            bool // EX or PX specified
	ExpiryTimeSeconds         int64
	ExpiryTimeMiliSeconds     int64
	ExpireTimestamp           bool // EXAT or PXAT specified
	ExpiryUnixTimeSeconds     int64
	ExpiryUnixTimeMiliSeconds int64
	KEEPTTL                   bool
}

func ResolveSetOptions(args ...string) (SetOptions, error) {
	options := SetOptions{
		NX:                        false,
		XX:                        false,
		GET:                       false,
		ExpireDuration:            false,
		ExpiryTimeSeconds:         -1,
		ExpiryTimeMiliSeconds:     -1,
		ExpireTimestamp:           false,
		ExpiryUnixTimeSeconds:     -1,
		ExpiryUnixTimeMiliSeconds: -1,
		KEEPTTL:                   false,
	}

	errString := "ERR syntax error"

	for argIdx, arg := range args {
		arg = strings.ToUpper(arg)

		switch arg {
		case "NX":
			options.NX = true
		case "XX":
			options.XX = true
		case "GET":
			options.GET = true
		case "EX":
			if len(args) > argIdx+1 {
				expiryTime, err := strconv.ParseInt(args[argIdx+1], 10, 64)
				if err != nil {
					errString := "ERR value is not an integer or out of range"
					return options, errors.New(errString)
				}
				options.ExpireDuration = true
				options.ExpiryTimeSeconds = expiryTime
			} else {
				return options, errors.New(errString)
			}
		case "PX":
			if len(args) > argIdx+1 {
				expiryTime, err := strconv.ParseInt(args[argIdx+1], 10, 64)
				if err != nil {
					errString := "ERR value is not an integer or out of range"
					return options, errors.New(errString)
				}
				options.ExpireDuration = true
				options.ExpiryTimeMiliSeconds = expiryTime
			} else {
				return options, errors.New(errString)
			}
		case "EXAT":
			if len(args) > argIdx+1 {
				expiryTime, err := strconv.ParseInt(args[argIdx+1], 10, 64)
				if err != nil {
					errString := "ERR value is not an integer or out of range"
					return options, errors.New(errString)
				}
				options.ExpireTimestamp = true
				options.ExpiryUnixTimeSeconds = expiryTime
			} else {
				return options, errors.New(errString)
			}
		case "PXAT":
			if len(args) > argIdx+1 {
				expiryTime, err := strconv.ParseInt(args[argIdx+1], 10, 64)
				if err != nil {
					return options, errors.New("ERR value is not an integer or out of range")
				}
				options.ExpireTimestamp = true
				options.ExpiryUnixTimeMiliSeconds = expiryTime
			} else {
				return options, errors.New(errString)
			}
		case "KEEPTTL":
			options.KEEPTTL = true
		}
	}

	return options, nil
}
