package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Message represents a chat message in the database.
type Message struct {
	gorm.Model
	Sender  string
	Message string
}

var (
	clients   = make(map[net.Conn]string) // List of clients and their names
	clientMux sync.Mutex                  // Mutex to protect clients map
	db        *gorm.DB                    // GORM database connection
)

func main() {
	var err error

	// Initialize GORM with SQLite
	db, err = gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

	// Automatically migrate the schema for the Message struct
	db.AutoMigrate(&Message{})

	// Start listening for incoming TCP connections
	ln, err := net.Listen("tcp", ":8080")
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

	// Fetch and send previous chat messages to the new client
	sendPreviousMessages(conn)

	// Broadcast that the client has joined
	broadcastMessage(fmt.Sprintf("%s joined the chat", conn.RemoteAddr().String()), conn)

	// Read incoming messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.TrimSpace(msg) == "" {
			continue
		}

		// Save message to the database using GORM
		saveMessage(clients[conn], msg)

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

// Save the message to the SQLite database using GORM
func saveMessage(sender, message string) {
	msg := Message{
		Sender:  sender,
		Message: message,
	}
	result := db.Create(&msg) // GORM's Create method
	if result.Error != nil {
		fmt.Println("Error saving message to database:", result.Error)
	}
}

// Send previous messages to the new client using GORM
func sendPreviousMessages(conn net.Conn) {
	var messages []Message
	result := db.Order("created_at").Find(&messages) // Fetch all messages sorted by creation time
	if result.Error != nil {
		fmt.Println("Error fetching previous messages:", result.Error)
		return
	}

	for _, msg := range messages {
		fmt.Fprintf(conn, "%s: %s\n", msg.Sender, msg.Message)
	}
}

