package main

import (
	"bufio"
	"fmt"
	"net"
	"server/handler"
	"strings"
	"sync"
	"time"
)

const colorRed = "\033[31m" // Angenehmes Rot
const colorReset = "\033[0m"

// Map to store connected clients (key: connection, value: client name)
var clients = make(map[net.Conn]string)
var mutex = &sync.Mutex{} // Mutex for safe concurrent access to clients map

// Function to handle the connection and cleanup
func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientName := handler.RegisterClient(conn, mutex, clients)
	fmt.Printf("Neuer Client verbunden: %s\n", clientName)
	handleClientMessages(conn, clientName)
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	fmt.Printf("Client disconnected: %s\n", clientName)
}

// Function to handle incoming messages from a client
func handleClientMessages(conn net.Conn, clientName string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			broadcast_message := fmt.Sprintf(colorRed+"%s hat den Chat verlassen.\n"+colorReset, clientName)
			handler.BroadcastMessage(conn, broadcast_message, mutex, clients)
			fmt.Println("Fehler beim Lesen der Nachricht von", clientName, ":", err)
			break
		}
		broadcast_message := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format("15:04:05"), clientName, strings.TrimRight(message, "\n"))
		handler.BroadcastMessage(conn, broadcast_message, mutex, clients)
	}
}

func main() {
	portInput := handler.GetPort()
	listener, err := net.Listen("tcp", portInput)
	if err != nil {
		fmt.Println("Fehler beim Lauschen:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat-Server gestartet auf Port", portInput)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn) // Handle each connection in a separate goroutine
	}
}
