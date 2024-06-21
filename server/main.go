package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

type MessageInterface interface {
	GetType() string
}

type GreetingMessage struct {
	Name string
	Body string
}

func (g GreetingMessage) GetType() string {
	return "greeting"
}

type StatusMessage struct {
	Status string
	Code   int
}

func (s StatusMessage) GetType() string {
	return "status"
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	decoder := gob.NewDecoder(reader)

	for {
		var msgType string
		if err := decoder.Decode(&msgType); err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error decoding message type:", err)
			}
			return
		}

		switch msgType {
		case "greeting":
			var msg GreetingMessage
			if err := decoder.Decode(&msg); err != nil {
				fmt.Println("Error decoding greeting message:", err)
				return
			}
			fmt.Printf("GreetingMessage received: %+v\n", msg)

		case "status":
			var msg StatusMessage
			if err := decoder.Decode(&msg); err != nil {
				fmt.Println("Error decoding status message:", err)
				return
			}
			fmt.Printf("StatusMessage received: %+v\n", msg)

		default:
			fmt.Println("Unknown message type received")
		}

		writer := bufio.NewWriter(conn)
		response := "Message received\n"
		if _, err := writer.WriteString(response); err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing response:", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen(connType, connHost+":"+connPort)

	if err != nil {
		fmt.Println("Error creating listener:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is over " + connType + " on " + connHost + ":" + connPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			//continue
			return
		}
		fmt.Println("Client connected.")
		fmt.Println("Client " + conn.RemoteAddr().String() + " connected.")
		go handleConnection(conn)
	}
}
