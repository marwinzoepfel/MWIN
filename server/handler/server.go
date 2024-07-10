package handler

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

const defaultPort = ":8080" // Choose your preferred port

// Function to get port (input or default)
func GetPort() string {
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

// Function to handle the connection and cleanup
func HandleConnection(conn net.Conn, mutex *sync.Mutex, clients map[net.Conn]string) {
	defer conn.Close()
	clientName := RegisterClient(conn, mutex, clients)
	fmt.Printf("Neuer Client verbunden: %s\n", clientName)
	handleClientMessages(conn, clientName, mutex, clients)
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
	fmt.Printf("Client disconnected: %s\n", clientName)
}
