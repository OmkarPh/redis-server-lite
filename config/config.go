package config

import (
	"log/slog"
)

const PORT = 6379

// const DefaultLoggerLevel = slog.LevelDebug

const DefaultLoggerLevel = slog.LevelInfo

const ShardFactor = 10
