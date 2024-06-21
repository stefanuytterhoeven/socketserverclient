package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

type GreetingMessage struct {
	Name string
	Body string
}

type StatusMessage struct {
	Status string
	Code   int
}

func sendMessage(conn net.Conn, msgType string, msg interface{}) error {
	writer := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(writer)

	if err := encoder.Encode(msgType); err != nil {
		return fmt.Errorf("error encoding message type: %w", err)
	}
	if err := encoder.Encode(msg); err != nil {
		return fmt.Errorf("error encoding message: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}
	return nil
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	var counter int = 0
	defer conn.Close()

	// Stuur een greeting message
	if err := sendMessage(conn, "greeting", GreetingMessage{Name: "Alice", Body: "Hello, server!"}); err != nil {
		fmt.Println("Error sending greeting message:", err)
		return
	}

	// Stuur een status message
	if err := sendMessage(conn, "status", StatusMessage{Status: "OK", Code: 200}); err != nil {
		fmt.Println("Error sending status message:", err)
		return
	}

	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by server")
			} else {
				fmt.Println("Error reading response:", err)
			}
			return
		}
		fmt.Printf("Server response: %s", response)
		counter += 1
		if counter >= 2 {
			return
		}
	}
}
