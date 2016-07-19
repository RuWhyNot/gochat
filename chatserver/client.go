package main

import (
	"fmt"
	"net"
	"bufio"
)


type Client struct {
	incoming chan string
	outgoing chan string
	connection net.Conn
	connected bool
}

func (thisClient *Client) Read() {
	for {
		str, err := bufio.NewReader(thisClient.connection).ReadString('\n')
		if (err != nil) {
			fmt.Printf("User was disconnected\n")
			thisClient.connected = false
			close(thisClient.incoming)
			close(thisClient.outgoing)
			return
		}
		thisClient.incoming <- str
	}
}

func (thisClient *Client) Write() {
	for {
		str, ok := <- thisClient.outgoing
		if (!ok) {
			return
		}
		thisClient.connection.Write([]byte(str))
	}
}

func (thisClient *Client) Listen(clients *[]*Client) {
	for {
		str, ok := <- thisClient.incoming
		if (!ok) {
			return
		}

		// downward loop for safe deleting clients
		for i := len(*clients) - 1; i >= 0; i-- {
			otherClient := (*clients)[i]
			
			if (otherClient.connected) {
				otherClient.outgoing <- str
			} else {
				*clients = append((*clients)[:i], (*clients)[i+1:]...)
			}
		}
	}
}

func (thisClient *Client) HandleConnection() {
	thisClient.connected = true
	go thisClient.Read()
	go thisClient.Write()
}

func AcceptClients(ln net.Listener, clients *[]*Client) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error establishing new connection\n")
			continue
		}
		
		fmt.Printf("New user is connected!\n")

		client := &Client{
			incoming: make(chan string),
			outgoing: make(chan string),
			connection: conn,
		}
		*clients = append(*clients, client)

		go client.Listen(clients)

		client.HandleConnection()
	}
}
