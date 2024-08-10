package network

import (
    "bufio"
    "client/notification"
    "strings"

    "fmt"
    "net"
    "time"
)

// ReceiveMessages receives messages from the server and calls the callback function.
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

        notification.SendMessageToOS(formattedMessage) 
        onMessageReceived(message)
    }
}

// SendMessage sends a message to the server and updates the GUI.
func SendMessage(conn net.Conn, message string, onMessageSent func(string)) {
    currentTime := time.Now().Format("15:04:05")
    formattedMessage := fmt.Sprintf("[%s] You: %s", currentTime, message)

    onMessageSent(formattedMessage)

    _, err := fmt.Fprintln(conn, message)
    if err != nil {
        fmt.Println("Error sending message:", err)
        // You could add error handling here, e.g., closing the connection
    }
}
