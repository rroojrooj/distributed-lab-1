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

func handleError(err error) {
	if err == io.EOF {
		fmt.Println("Client disconnected :(")
	} else {
		fmt.Printf("Error occurred: %s\n", err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	conns <- conn
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
			handleError(err)
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
