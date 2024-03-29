package core

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/OmkarPh/redis-lite/store"
)

func CleanExpiredKeys(kvStore *store.KvStore) {
	checkInterval := 5 * time.Second
	keysCount := 25
	timeoutKeyscountThreshold := 100

	for {
		expiredKeys := (*kvStore).DeleteIfExpired(keysCount)
		if expiredKeys > 0 {
			slog.Debug(fmt.Sprintf("Removed %d expired keys", expiredKeys))
		}
		if expiredKeys < timeoutKeyscountThreshold {
			time.Sleep(checkInterval)
		}
	}
}
