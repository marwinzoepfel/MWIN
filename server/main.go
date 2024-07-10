package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultPort = ":8080" // Choose your preferred port
const colorRed = "\033[31m" // Angenehmes Rot
const colorReset = "\033[0m"

// Map to store connected clients (key: connection, value: client name)
var clients = make(map[net.Conn]string)
var mutex = &sync.Mutex{} // Mutex for safe concurrent access to clients map

// Function to get client name from the connection
func getClientName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	clientName, _ := reader.ReadString('\n')
	return strings.TrimSpace(clientName)
}

// Function to register a client and send the join message
func registerClient(conn net.Conn) string {
	clientName := getClientName(conn)
	mutex.Lock()
	clients[conn] = clientName
	mutex.Unlock()
	broadcastMessage(conn, fmt.Sprintf("%s ist dem Chat beigetreten!\n", clientName))
	return clientName
}

// Function to handle incoming messages from a client
func handleClientMessages(conn net.Conn, clientName string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			broadcastMessage(conn, fmt.Sprintf(colorRed+"%s hat den Chat verlassen.\n"+colorReset, clientName))
			fmt.Println("Fehler beim Lesen der Nachricht von", clientName, ":", err)
			break
		}
		broadcastMessage(conn, fmt.Sprintf("[%s] %s: %s\n", time.Now().Format("15:04:05"), clientName, strings.TrimRight(message, "\n")))
	}
}

// Function to handle the connection and cleanup
func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientName := registerClient(conn)
	fmt.Printf("Neuer Client verbunden: %s\n", clientName)
	handleClientMessages(conn, clientName)
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	fmt.Printf("Client disconnected: %s\n", clientName)
}

func broadcastMessage(sender net.Conn, message string) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		if client != sender {
			fmt.Fprint(client, message)
		}
	}
}

// Function to get port (input or default)
func getPort() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Port eingeben (Standard: %s): [Wichtig! Ohne :] ", defaultPort)
	portInput, _ := reader.ReadString('\n')
	portInput = strings.TrimSpace(portInput)
	if portInput == "" {
		return defaultPort
	}
	if !strings.HasPrefix(portInput, ":") {
		portInput = ":" + portInput
	}
	return portInput
}

func main() {
	portInput := getPort()
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
