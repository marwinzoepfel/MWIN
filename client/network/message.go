package network

import (
	"bufio"
	"client/notificatioin"
	"strings"

	"fmt"
	"net"
	"time"
)

// ReceiveMessages empfängt Nachrichten vom Server und ruft die Callback-Funktion auf.
func ReceiveMessages(conn net.Conn, onMessageReceived func(string)) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection to server interrupted:", err)
			return
		}
		if len(message) == 0 {
			continue
		}
		//currentTime := time.Now().Format("15:04:05")
		message = message[:len(message)-1]

		formattedMessage := strings.ReplaceAll(message, "\033[38;5;153m%s\033", "")
		message = fmt.Sprintf("\033[38;5;153m%s\033[0m", message)

		notificatioin.SendMessageToOS(formattedMessage)
		onMessageReceived(message)
	}
}

// SendMessage sendet eine Nachricht an den Server und aktualisiert die GUI.
func SendMessage(conn net.Conn, message string, onMessageSent func(string)) {
	currentTime := time.Now().Format("15:04:05")
	formattedMessage := fmt.Sprintf("[%s] You: %s", currentTime, message)

	onMessageSent(formattedMessage)

	_, err := fmt.Fprintln(conn, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		// Hier könntest du eine Fehlerbehandlung hinzufügen, z.B. die Verbindung schließen
	}
}
