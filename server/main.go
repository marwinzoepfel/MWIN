package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"server/utils"
	"strings"
	"sync"
	"time"
)

const defaultPort = ":8080" // Choose your preferred port

// Map to store connected clients (key: connection, value: client name)
var clients = make(map[net.Conn]string)
var mutex = &sync.Mutex{} // Mutex for safe concurrent access to clients map

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Client-Namen empfangen
	reader := bufio.NewReader(conn)
	clientName, _ := reader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)

	mutex.Lock()
	clients[conn] = clientName
	mutex.Unlock()

	broadcastMessage(conn, fmt.Sprintf("%s ist dem Chat beigetreten!\n", clientName))

	fmt.Printf("Neuer Client verbunden: %s\n", clientName)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			broadcastMessage(conn, fmt.Sprintf(utils.ColorRed+"%s hat den Chat verlassen.\n"+utils.ColorReset, clientName))

			fmt.Println("Fehler beim Lesen der Nachricht von", clientName, ":", err)
			break
		}

		currentTime := time.Now().Format("15:04:05")

		// Broadcast the message to all clients (with correct sender name)
		mutex.Lock()
		for client := range clients {
			if client != conn {
				// Nachricht mit Client-Namen senden, aber nur den Zeilenumbruch entfernen
				fmt.Fprintf(client, "[%s] %s: %s\n", currentTime, clientName, strings.TrimRight(message, "\n"))
			}
		}
		mutex.Unlock()
	}
	// Remove the client from the map when they disconnect
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

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Port eingeben (Standard: %s): [Wichtig! Ohne :] ", defaultPort)
	portInput, _ := reader.ReadString('\n')
	portInput = strings.TrimSpace(portInput)

	// Wenn die Eingabe leer ist, verwende den Standardport
	if portInput == "" {
		portInput = defaultPort
	} else {
		// Füge einen Doppelpunkt hinzu, wenn der Benutzer keinen eingegeben hat
		if !strings.HasPrefix(portInput, ":") {
			portInput = ":" + portInput
		}
	}

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
		// Now call the handleConnection function to handle the new connection
		go handleConnection(conn)
	}
}
