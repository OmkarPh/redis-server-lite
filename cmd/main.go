package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Read more - https://betterprogramming.pub/redis-internals-client-connects-to-redis-by-tcp-faf6ad1a7df5
func HandleConnection(conn net.Conn) {
	conn_id := uuid.New()
	// slog.Debug(fmt.Sprintf("Serving client %s\n", conn.RemoteAddr().String()))

	defer conn.Close()

	messages := 1
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

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

		var response []byte
		commands := ResolveCommands(*reader, requestType)

		slog.Debug(fmt.Sprintf("Commands: %+v, %d commands\n", commands, len(commands)))

		for _, command := range commands {
			slog.Debug(fmt.Sprintf("Exec Command: %+v, %d args\n", command, len(command.args)))

			// Execute action based on action key
			if action, ok := RedisActions[command.action]; ok {
				currCommandResponse, err := action.Execute(command.args...)
				if err != nil {
					slog.Debug("Error executing action:", err)
				}
				response = append(response, currCommandResponse...)
			} else {
				slog.Debug(fmt.Sprint("Invalid action key: ", command.action))
			}
		}

		slog.Debug(fmt.Sprint("Sent response", string(response)))
		conn.Write(response)

		slog.Debug(fmt.Sprintf("Conn %s _ Msg #%d processed\n", conn_id, messages))
		messages++
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	}

	conn.Close()
	slog.Debug(fmt.Sprintf("Closed connection %s", conn_id.String()))
}

func main() {
	slog.SetLogLoggerLevel(DefaultLoggerLevel)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		slog.Debug("Error starting server", err)
		return
	}
	defer listener.Close()

	fmt.Println("Redis-lite server is up & running on port", PORT)
	fmt.Println()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Debug("Error accepting connection", err)
			return
		}

		slog.Debug(fmt.Sprint("Accept connection from", conn.RemoteAddr().String()))
		// Handle the connection in a new goroutine.
		go HandleConnection(conn)
	}
}
