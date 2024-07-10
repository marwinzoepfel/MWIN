package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	colorReset = "\033[0m"
	colorCyan  = "\033[36m"
)

// ReceiveMessages continuously reads and prints messages from the server connection.
func ReceiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection to server interrupted:", err)
			return // Exit the goroutine if there's an error
		}

		// Check if the message is empty (e.g., due to a disconnection)
		if len(message) == 0 {
			continue
		}

		// Handle potential errors when printing the message (e.g., if the terminal is closed)
		if _, err := fmt.Print(message); err != nil {
			fmt.Println("Error printing message:", err)
			return // Exit the goroutine if there's an error
		}
	}
}

// SendMessages continuously reads messages from the user's input and sends them to the server.
func SendMessages(conn net.Conn, reader *bufio.Reader) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		currentTime := time.Now().Format("15:04:05")

		// Print the formatted message with a cyan color for the timestamp and "You:"
		_, err := fmt.Printf(colorCyan+"[%s] You: %s\n"+colorReset, currentTime, message)
		if err != nil {
			fmt.Println("Error printing message:", err)
			return // Exit the function if there's an error
		}

		// Send only the message content, without the timestamp and "You:"
		_, err = fmt.Fprintln(conn, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return // Exit the function if there's an error
		}
	}

	// Handle potential errors from the scanner
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from console:", err)
	}
}
