package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/stefanuytterhoeven/socketserverclient/shared"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type MessageInterface interface {
	GetType() string
}

type GreetingMessage shared.GreetingMessage

//type GreetingMessage struct {
//	Name string
//	Body string
//}

func (g GreetingMessage) GetType() string {
	return "greeting"
}

type StatusMessage shared.StatusMessage

//type StatusMessage struct {
//	Status string
//	Code   int
//}

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
				InfoLogger.Println("Connection closed by client")
			} else {
				ErrorLogger.Println("Error decoding message type:", err)
			}
			return
		}
		var msgReceived string
		switch msgType {
		case "greeting":
			var msg GreetingMessage
			if err := decoder.Decode(&msg); err != nil {
				ErrorLogger.Println("Error decoding greeting message:", err)
				return
			}
			//fmt.Printf("Greeting Message received: %+v\n", msg)
			msgReceived = fmt.Sprintf("Greeting Message received: %+v", msg)
			fmt.Println(msgReceived)
			InfoLogger.Println(msgReceived)

		case "status":
			var msg StatusMessage
			if err := decoder.Decode(&msg); err != nil {
				ErrorLogger.Println("Error decoding status message:", err)
				return
			}
			//fmt.Printf("Status Message received: %+v\n", msg)
			msgReceived = fmt.Sprintf("Status Message received: %+v", msg)
			fmt.Println(msgReceived)
			InfoLogger.Println(msgReceived)
		default:
			//fmt.Println("Unknown message type received")
			msgReceived = "Unknown message type received"
			fmt.Println(msgReceived)
			WarningLogger.Println(msgReceived)
		}

		writer := bufio.NewWriter(conn)
		//response := "Message received\n"
		response := msgReceived + "\n"
		if _, err := writer.WriteString(response); err != nil {
			ErrorLogger.Println("Error writing response:", err)
			return
		}
		if err := writer.Flush(); err != nil {
			ErrorLogger.Println("Error flushing response:", err)
			return
		}
	}
}

func init() {
	file, err := os.OpenFile("../logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		fmt.Println("logging failed")
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	listener, err := net.Listen(connType, connHost+":"+connPort)

	if err != nil {
		ErrorLogger.Println("Error creating listener:", err)
		os.Exit(1)
	}
	defer listener.Close()
	InfoLogger.Println("Server is over " + connType + " on " + connHost + ":" + connPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			ErrorLogger.Println("Error accepting connection:", err)
			//continue
			return
		}
		InfoLogger.Println("Client connected.")
		InfoLogger.Println("Client " + conn.RemoteAddr().String() + " connected.")
		go handleConnection(conn)
	}
}
