package handler

import (
	"fmt"
	"net"
	"sync"
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
