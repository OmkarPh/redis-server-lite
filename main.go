package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/core"
	"github.com/OmkarPh/redis-lite/store"
)

func main() {
	slog.SetLogLoggerLevel(config.DefaultLoggerLevel)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.PORT))
	if err != nil {
		slog.Error("Error starting server", err)
		return
	}
	defer listener.Close()

	redisConfig := config.NewRedisConfig()
	kvStore := store.NewKvStore(redisConfig)

	fmt.Println("Redis-lite server is up & running on port", config.PORT)
	fmt.Println()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error accepting connection", err)
			return
		}

		slog.Debug(fmt.Sprint("Accept connection from", conn.RemoteAddr().String()))
		// Handle the connection in a new goroutine.
		go core.HandleConnection(conn, redisConfig, kvStore)
	}
}
