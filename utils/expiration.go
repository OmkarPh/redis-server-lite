package utils

import (
	"errors"
	"time"
)

type ValueWithExpiration struct {
	Expiry time.Time
}

func IsExpired(expiration time.Time) bool {
	return !expiration.IsZero() && expiration.Before(time.Now())
}

type ExpirationTimeOptions struct {
	NX                        bool
	XX                        bool
	ExpireDuration            bool // EX or PX specified
	ExpiryTimeSeconds         int64
	ExpiryTimeMiliSeconds     int64
	ExpireTimestamp           bool // EXAT or PXAT specified
	ExpiryUnixTimeSeconds     int64
	ExpiryUnixTimeMiliSeconds int64
	KEEPTTL                   bool
}

func ResolveExpirationTime(options ExpirationTimeOptions, exists bool, existingExpiry time.Time) (time.Time, bool, error) {
	expiryTime := time.Time{}

	if options.KEEPTTL && exists {
		expiryTime = existingExpiry
	}

	if options.ExpireDuration {
		if options.ExpiryTimeSeconds != -1 {
			expiryTime = time.Now().Add(time.Second * time.Duration(options.ExpiryTimeSeconds))
		} else if options.ExpiryTimeMiliSeconds != -1 {
			expiryTime = time.Now().Add(time.Millisecond * time.Duration(options.ExpiryTimeMiliSeconds))
		} else {
			errString := "ERR invalid expire duration in SET"
			return expiryTime, false, errors.New(errString)
		}
	}

	if options.ExpireTimestamp {
		if options.ExpiryUnixTimeSeconds != -1 {
			expiryTime = time.Unix(options.ExpiryUnixTimeSeconds, 0)
		} else if options.ExpiryUnixTimeMiliSeconds != -1 {
			nanoseconds := options.ExpiryUnixTimeMiliSeconds * 1000000 * int64(time.Nanosecond)
			expiryTime = time.Unix(0, nanoseconds)
		} else {
			errString := "ERR invalid expire timestamp in SET"
			return expiryTime, false, errors.New(errString)
		}
	}

	return expiryTime, true, nil
}
