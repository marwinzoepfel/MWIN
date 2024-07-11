package main

import (
	"bufio"
	"fmt"
	"os"

	"client/gui"     // Passe den Importpfad entsprechend an
	"client/network" // Passe den Importpfad entsprechend an
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	serverAddress := network.GetServerAddress(reader)
	clientName := network.GetClientName(reader)

	conn, err := network.ConnectToServer(serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully joined.")
	if err := network.SendInitialName(conn, clientName); err != nil {
		fmt.Println("Error sending name:", err)
		return
	}

	// Start GUI
	gui.RunChat(conn, reader)
}
