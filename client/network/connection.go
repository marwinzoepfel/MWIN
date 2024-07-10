package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const defaultServerAddress = "localhost:8080"

// GetServerAddress prompts the user for a server address, using the default if none is provided.
func GetServerAddress(reader *bufio.Reader) string {
	fmt.Printf("Enter server address (default: %s): ", defaultServerAddress)
	serverAddress, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading server address:", err)
		return defaultServerAddress // Fall back to default
	}

	serverAddress = strings.TrimSpace(serverAddress)
	if serverAddress == "" {
		return defaultServerAddress
	}

	return serverAddress
}

// GetClientName prompts the user for a client name, using the hostname as default.
func GetClientName(reader *bufio.Reader) string {
	defaultName, err := os.Hostname()
	if err != nil {
		defaultName = "Unknown User"
	}

	fmt.Printf("Your name (default: %s): ", defaultName)
	inputName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return defaultName // Fall back to default
	}

	inputName = strings.TrimSpace(inputName)
	if inputName != "" {
		return inputName
	}

	return defaultName
}

// ConnectToServer establishes a TCP connection to the specified server address.
func ConnectToServer(serverAddress string) (net.Conn, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s: %v", serverAddress, err)
	}
	return conn, nil
}

// SendInitialName sends the client's name to the server immediately after connecting.
func SendInitialName(conn net.Conn, clientName string) error {
	_, err := fmt.Fprintln(conn, clientName)
	if err != nil {
		return fmt.Errorf("error sending name: %v", err)
	}
	return nil
}
