package handler

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Function to get client name from the connection
func getClientName(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	clientName, _ := reader.ReadString('\n')
	return strings.TrimSpace(clientName)
}

// Function to register a client and send the join message
func RegisterClient(conn net.Conn, mutex *sync.Mutex, clients map[net.Conn]string) string {
	clientName := getClientName(conn)
	mutex.Lock()
	clients[conn] = clientName
	mutex.Unlock()
	message := fmt.Sprintf("%s ist dem Chat beigetreten!\n", clientName)
	BroadcastMessage(conn, message, mutex, clients)
	return clientName
}
