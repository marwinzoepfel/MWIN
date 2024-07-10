package interfaces

import (
	"fmt"
	"net"
	"sync"
)

// Broadcaster defines the interface for broadcasting messages.
type Broadcaster interface {
	Broadcast(sender Client, message string)
}

// Client represents a connected client.
type Client struct {
	Conn        net.Conn
	Name        string
	clientList  ClientList // Embed ClientList for direct access to its methods
	broadcaster Broadcaster
}

// Add adds the client to the client list.
func (c *Client) Add() {
	c.clientList.Add(c)
}

// Remove removes the client from the client list.
func (c *Client) Remove() {
	c.clientList.Remove(c)
}

// ClientList manages a list of connected clients.
type ClientList struct {
	clients map[net.Conn]*Client
	mutex   sync.Mutex
}

// NewClientList creates a new ClientList instance.
func NewClientList() *ClientList {
	return &ClientList{
		clients: make(map[net.Conn]*Client),
	}
}

// Add adds a client to the list.
func (cl *ClientList) Add(client *Client) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.clients[client.Conn] = client
}

// Remove removes a client from the list.
func (cl *ClientList) Remove(client *Client) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	delete(cl.clients, client.Conn)
}

// Broadcast sends a message to all clients except the sender.
func (cl *ClientList) Broadcast(sender *Client, message string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	for _, client := range cl.clients {
		if client != sender {
			fmt.Fprint(client.Conn, message)
		}
	}
}
