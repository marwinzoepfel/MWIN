package message

import "server/client"

// Broadcast sendet eine Nachricht an alle Clients außer dem Sender.
func Broadcast(sender *client.Client, message string) {
	clientList := client.GetClientList() // Zentrale Client-Liste erhalten
	clientList.Broadcast(sender, message)
}
