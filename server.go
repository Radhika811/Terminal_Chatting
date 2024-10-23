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


	fmt.Println("New client connected:", conn.RemoteAddr())

	// Broadcast that the client has joined
	broadcastMessage(fmt.Sprintf("%s joined the chat", conn.RemoteAddr().String()), conn)

	// Read incoming messages from the client
	

	// Handle client disconnection

	fmt.Println("Client disconnected:", conn.RemoteAddr())
	broadcastMessage(fmt.Sprintf("%s left the chat", conn.RemoteAddr().String()), conn)
}


// Broadcasts a message to all clients except the sender
func broadcastMessage(message string, sender net.Conn) {
	
}

