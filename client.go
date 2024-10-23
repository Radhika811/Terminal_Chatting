package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the chat server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to chat server.")

	// Goroutine to handle receiving messages from the server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			msg := scanner.Text()
			fmt.Println(msg)
		}
	}()

	// Send messages to the server
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		message, _ := input.ReadString('\n')
		fmt.Fprint(conn, message)
	}
}
