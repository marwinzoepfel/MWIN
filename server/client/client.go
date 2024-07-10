package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"server/interfaces"
	"server/utils"
)

// NewClient creates a new Client and initializes its fields.
func NewClient(conn net.Conn, clientList interfaces.ClientList) *interfaces.Client {
	return &interfaces.Client{
		Conn:       conn,
		clientList: clientList,
	}
}

// Handle is responsible for handling the communication with a single client.
func (c *interfaces.Client) Handle() {
	defer c.Conn.Close()
	reader := bufio.NewReader(c.Conn)

	// Get and set the client's name
	if name, err := reader.ReadString('\n'); err == nil {
		c.Name = strings.TrimSpace(name)
	} else {
		// Handle the error gracefully if the client disconnects before sending the name
		fmt.Println("Client disconnected before providing a name:", err)
		return // Exit the handle function since we don't have a valid client name
	}

	// Add the client to the client list
	c.Add()

	// Broadcast a welcome message
	c.clientList.Broadcast(c, fmt.Sprintf(utils.ColorGreen+"%s ist dem Chat beigetreten!\n"+utils.ColorReset, c.Name))
	fmt.Printf("Neuer Client verbunden: %s\n", c.Name)

	// Start the main message loop
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.clientList.Broadcast(c, fmt.Sprintf(utils.ColorRed+"%s hat den Chat verlassen.\n"+utils.ColorReset, c.Name))
			fmt.Println("Fehler beim Lesen der Nachricht von", c.Name, ":", err)
			break
		}

		currentTime := time.Now().Format("15:04:05")
		c.clientList.Broadcast(c, fmt.Sprintf("[%s] %s: %s", currentTime, c.Name, strings.TrimRight(message, "\n")))
	}

	// Remove the client from the client list
	c.Remove()
	fmt.Printf("Client disconnected: %s\n", c.Name)
}
