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

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type GreetingMessage shared.GreetingMessage

type StatusMessage shared.StatusMessage

func sendMessage(conn net.Conn, msgType string, msg interface{}) error {
	writer := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(writer)
	InfoLogger.Println("Sending "+msgType+" ", msg)
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

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		ErrorLogger.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	var counter int = 0
	defer conn.Close()

	// Stuur een greeting message

	if err := sendMessage(conn, "greeting", GreetingMessage{Name: "Alice", Body: "Hello, server!"}); err != nil {
		ErrorLogger.Println("Error sending greeting message:", err)
		return
	}

	// Stuur een status message
	if err := sendMessage(conn, "status", StatusMessage{Status: "OK", Code: 200}); err != nil {
		ErrorLogger.Println("Error sending status message:", err)
		return
	}

	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				InfoLogger.Println("Connection closed by server")
			} else {
				WarningLogger.Println("Error reading response:", err)
			}
			return
		}
		InfoLogger.Printf("Server response: %s", response)
		counter += 1
		if counter >= 2 {
			return
		}
	}
}
