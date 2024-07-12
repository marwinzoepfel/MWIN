package main

import (
	"fmt"
	"net"
	"server/handler"
	"sync"
)

// Map to store connected clients (key: connection, value: client name)
var clients = make(map[net.Conn]string)
var mutex = &sync.Mutex{} // Mutex for safe concurrent access to clients map

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
		go handler.HandleConnection(conn, mutex, clients) // Handle each connection in a separate goroutine
	}
}
