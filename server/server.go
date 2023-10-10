package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error, clientid int) {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		fmt.Printf("Error with client %d: Connection timed out.\n", clientid)
	} else if opErr, ok := err.(*net.OpError); ok {
		if opErr.Op == "read" {
			fmt.Printf("Error with client %d: Connection was reset by peer.\n", clientid)
		} else if opErr.Op == "write" {
			fmt.Printf("Error with client %d: Failed to write data.\n", clientid)
		}
	} else {
		fmt.Printf("Error with client %d: %s\n", clientid, err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue // If there's an error, skip the rest of the loop and try to accept the next connection.
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			// Handling the EOF error separately as it's a common and expected case
			if err == io.EOF {
				fmt.Printf("Client %d disconnected :(\n", clientid)
			} else {
				handleError(err, clientid)
			}
			return
		}

		// Check for exit command
		if strings.TrimSpace(message) == "exit" {
			fmt.Printf("Client %d exited\n", clientid)
			return
		}

		// Tidy up each message (trimming newline and spaces)
		tidyMessage := strings.TrimSpace(message)

		// Create a Message struct with the tidied message and client ID
		msg := Message{sender: clientid,
			message: tidyMessage}

		fmt.Printf("Client %d says: %s\n", clientid, msg.message)

		// Send this Message to the msgs channel
		msgs <- msg

		// Sending acknowledgment back to the client
		reply := tidyMessage + "\n"
		_, err = client.Write([]byte(reply))
		if err != nil {
			handleError(err, clientid)
			return
		}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above. Check

	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	var clientIDCounter int = 0 // This will be used to give each client a unique ID

	//Start accepting connections
	go acceptConns(ln, conns)

	for {
		select {
		case conn := <-conns:

			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client

			clientIDCounter++                            // Increment the counter to get a new unique ID
			clients[clientIDCounter] = conn              // Add the new client to the clients map
			go handleClient(conn, clientIDCounter, msgs) // Start handling messages from the new client
		case msg := <-msgs:

			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for id, client := range clients {
				if id != msg.sender { // Do not send the message back to the sender
					fmt.Fprintf(client, "%s\n", msg.message)
				}
			}
		}
	}
}
