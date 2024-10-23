package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"strings"
)

var (
	clients   = make(map[net.Conn]string) // List of clients and their names
	clientMux sync.Mutex                  // Mutex to protect clients map
)

func main() {
	var err error

	// Start listening for incoming TCP connections
	ln, err := net.Listen("tcp", "10.81.17.131:8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Chat server started on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}


// Handles client communication
func handleClient(conn net.Conn) {
	defer conn.Close()

	clientMux.Lock()
	clients[conn] = conn.RemoteAddr().String() // Add the client to the list
	clientMux.Unlock()

	fmt.Println("New client connected:", conn.RemoteAddr())

	// Broadcast that the client has joined
	broadcastMessage(fmt.Sprintf("%s joined the chat", conn.RemoteAddr().String()), conn)

	// Read incoming messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.TrimSpace(msg) == "" {
			continue
		}

		// Broadcast the message to all clients
		broadcastMessage(fmt.Sprintf("%s: %s", clients[conn], msg), conn)
	}

	// Handle client disconnection
	clientMux.Lock()
	delete(clients, conn)
	clientMux.Unlock()

	fmt.Println("Client disconnected:", conn.RemoteAddr())
	broadcastMessage(fmt.Sprintf("%s left the chat", conn.RemoteAddr().String()), conn)
}


// Broadcasts a message to all clients except the sender
func broadcastMessage(message string, sender net.Conn) {
	fmt.Println("Broadcasting:", message)
	clientMux.Lock()
	for client := range clients {
		if client != sender {
			fmt.Fprintln(client, message)
		}
	}
	clientMux.Unlock()
}

