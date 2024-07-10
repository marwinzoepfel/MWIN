package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"server/client"
	"strings"
)

const defaultPort = ":8080"

func main() {
	port := getPortFromUser()
	listener := startServer(port)
	defer listener.Close()

	fmt.Println("Chat-Server gestartet auf Port", port)

	acceptConnections(listener)
}

func getPortFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Port eingeben (Standard: %s): ", defaultPort)

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

func startServer(port string) net.Listener {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Fehler beim Lauschen:", err)
		os.Exit(1) // Beende das Programm bei einem Fehler
	}
	return listener
}

func acceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := client.NewClient(conn)
		go client.Handle()
	}
}
