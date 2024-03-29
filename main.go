package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/core"
	"github.com/OmkarPh/redis-lite/store"
)

func main() {

	fmt.Println("Redis-lite server v0.2")
	fmt.Println("Port:", config.PORT, ", PID:", os.Getpid())
	fmt.Println("Visit - https://github.com/OmkarPh/redis-server-lite")
	fmt.Println()

	slog.SetLogLoggerLevel(config.DefaultLoggerLevel)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.PORT))
	if err != nil {
		slog.Error("Error starting server", err)
		return
	}
	defer listener.Close()

	redisConfig := config.NewRedisConfig()
	kvStore := store.NewKvStore(redisConfig)

	go core.CleanExpiredKeys(kvStore)

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
