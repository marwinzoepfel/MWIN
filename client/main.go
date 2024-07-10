package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"client/handler"
)

const defaultServerAddress = "localhost:8080"
const colorReset = "\033[0m"
const colorCyan = "\033[36m" // Blaue Farbe für die Uhrzeit

func main() {

	handler.Hello()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Serveradresse eingeben (Standard: %s): ", defaultServerAddress)
	serverAddress, _ := reader.ReadString('\n')
	serverAddress = strings.TrimSpace(serverAddress)

	// Wenn die Eingabe leer ist, verwende die Standardadresse
	if serverAddress == "" {
		serverAddress = defaultServerAddress
	}

	// Standardmäßig den Gerätenamen als Client-Namen verwenden
	clientName, err := os.Hostname()
	if err != nil {
		clientName = "Unbekannter Benutzer"
	}

	fmt.Printf("Dein Name (Standard: %s): ", clientName)
	inputName, _ := reader.ReadString('\n')
	inputName = strings.TrimSpace(inputName)

	if inputName != "" {
		clientName = inputName
	}

	// Verbindung zum Server herstellen
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Fehler beim Verbinden zum Server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Du bist erfolgreich beigetreten.")
	// Namen sofort nach dem Verbindungsaufbau senden
	fmt.Fprintln(conn, clientName)

	// Nachrichten vom Server lesen (in einer Goroutine)
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Verbindung zum Server unterbrochen:", err)
				return
			}
			fmt.Print(message) // Nachricht ohne Änderung ausgeben
		}
	}()

	// Nachrichten von der Konsole lesen und an den Server senden
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		currentTime := time.Now().Format("15:04:05")
		fmt.Printf(colorCyan+"[%s] Du: %s\n"+colorReset, currentTime, message) // Gesamte Nachricht in Cyan
		fmt.Fprintln(conn, message)                                            // Nur die Nachricht senden, ohne Uhrzeit und "Du:"
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Fehler beim Lesen von der Konsole:", err)
	}
}
