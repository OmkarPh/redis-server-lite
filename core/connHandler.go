package core

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/OmkarPh/redis-lite/config"
	"github.com/OmkarPh/redis-lite/core/actions"
	"github.com/OmkarPh/redis-lite/store"
	"github.com/google/uuid"
)

// Read more - https://betterprogramming.pub/redis-internals-client-connects-to-redis-by-tcp-faf6ad1a7df5
func HandleConnection(conn net.Conn, redisConfig *config.RedisConfig, kvStore *store.KvStore) {
	conn_id := uuid.New()
	// slog.Debug(fmt.Sprintf("Serving client %s\n", conn.RemoteAddr().String()))

	defer conn.Close()

	messages := 1
	timeout, timeoutConfigSet := (*redisConfig).GetParam("timeout")
	timeoutPeriod, err := strconv.Atoi(timeout)
	timeoutApplicable := timeoutConfigSet && err == nil
	if timeoutApplicable {
		conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeoutPeriod)))
	}

	for {
		scanner := bufio.NewScanner(conn)
		messageBuffer := make([]byte, 1000)
		scanner.Buffer(messageBuffer, 1000)

		scanned := scanner.Scan()

		if !scanned {
			// slog.Debug("Read message timedout")
			break
		}

		received_msg_str := strings.Trim(string(messageBuffer), "\r\n")
		slog.Debug(fmt.Sprintf("Received (%d): %s", len(received_msg_str), received_msg_str))

		// Reader for progressive command parser
		reader := bufio.NewReader(bytes.NewBuffer(messageBuffer[:]))

		identifierBytes, _ := reader.Peek(1)
		requestType := ResolveRequestType(identifierBytes[0])
		slog.Debug(fmt.Sprint("Serve request", requestType))

		commands := ResolveCommands(*reader, requestType)

		slog.Debug(fmt.Sprintf("Commands: %+v, %d commands\n", commands, len(commands)))

		for _, command := range commands {
			slog.Debug(fmt.Sprintf("Exec Command: %+v, %d args\n", command, len(command.Args)))

			// Execute action based on action key
			if action, ok := actions.RedisActions[command.Action]; ok {
				currCommandResponses, err := action.Execute(kvStore, redisConfig, command.Args...)
				if err != nil {
					slog.Debug("Error executing action:", err)
				}
				for _, response := range currCommandResponses {
					slog.Debug(fmt.Sprint("Sent response", string(response)))
					conn.Write(response)
				}
			} else {
				slog.Debug(fmt.Sprint("Invalid action key: ", command.Action))
			}
		}

		slog.Debug(fmt.Sprintf("Conn %s _ Msg #%d processed\n", conn_id, messages))
		messages++

		if timeoutApplicable {
			conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeoutPeriod)))
		}
	}

	conn.Close()
	slog.Debug(fmt.Sprintf("Closed connection %s", conn_id.String()))
}
