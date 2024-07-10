package handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const defaultPort = ":8080" // Choose your preferred port

// Function to get port (input or default)
func GetPort() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Port eingeben (Standard: %s): [Wichtig! Ohne :] ", defaultPort)
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
