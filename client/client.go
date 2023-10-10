package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func read(conn net.Conn) {
	//TODO In a continuous loop, read a message from the server and display it.

	// In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Print("Received: ", message)
	}
}

func write(conn net.Conn) {
	//TODO Continually get input from the user and send messages to the server.

	// Create a new buffered reader to read input from the standard input (typically the keyboard).
	reader := bufio.NewReader(os.Stdin)

	for { // for loop to make waiting for input forever. Not only one message and gone

		fmt.Printf("Enter your message (or 'exit' to quit): ")
		// Read a string from the user until a newline character (\n) is encountered.
		// The `ReadString` function will return the entered string, including the newline character.
		message, _ := reader.ReadString('\n')

		// Check for exit command
		if strings.TrimSpace(message) == "exit" {
			break
		}

		// Attempt to write (send) the user's message to the server via the provided connection.
		// The `Write` function sends data over the connection and returns the number of bytes written and any error encountered.
		_, err := conn.Write([]byte(message))

		// Check if there was an error while sending the message.
		if err != nil {
			// If there was an error, print it to the console.
			fmt.Println("Error sending message:", err)
			// Exit the function since we encountered an error.
			return
		}

	}

}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	// Try to connect to the server
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start asynchronously reading and displaying messages
	go read(conn)

	// Start getting and sending user messages.
	write(conn)

}
