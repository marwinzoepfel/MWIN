package handler

import (
	"bufio"
	"fmt"
	"net"
	"server/utils"
	"strings"
	"sync"
	"time"
)

func BroadcastMessage(sender net.Conn, message string, mutex *sync.Mutex, clients map[net.Conn]string) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		if client != sender {
			fmt.Fprint(client, message)
		}
	}
}

// Function to handle incoming messages from a client
func handleClientMessages(conn net.Conn, clientName string, mutex *sync.Mutex, clients map[net.Conn]string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			broadcast_message := fmt.Sprintf(utils.ColorRed+"%s hat den Chat verlassen.\n"+utils.ColorReset, clientName)
			BroadcastMessage(conn, broadcast_message, mutex, clients)
			fmt.Println("Fehler beim Lesen der Nachricht von", clientName, ":", err)
			break
		}
		broadcast_message := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format("15:04:05"), clientName, strings.TrimRight(message, "\n"))
		BroadcastMessage(conn, broadcast_message, mutex, clients)
	}
}
