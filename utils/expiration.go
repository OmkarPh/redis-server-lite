package utils

import (
	"time"
)

type ValueWithExpiration struct {
	Expiry time.Time
}

func IsExpired(expiration time.Time) bool {
	return !expiration.IsZero() && expiration.Before(time.Now())
}
