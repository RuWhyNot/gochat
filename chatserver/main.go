package main

import "fmt"
import "net"
import "bufio"


type Client struct {
	incoming chan string
	outgoing chan string
}

func Read(incoming chan string, con net.Conn) {
	for {
		str, err := bufio.NewReader(con).ReadString('\n')
		if (err != nil) {
			fmt.Printf("Error Read\n")
			break
		}
		fmt.Printf(str)
		incoming <- str
	}
}

func Write(outcoming chan string, con net.Conn) {
	for {
		str := <- outcoming
		con.Write([]byte(str))
	}
}

func handleConnection(con net.Conn, client *Client) {
	fmt.Printf("New user is connected!\n")
	go Read(client.incoming, con)
	go Write(client.outgoing, con)
}

func (thisClient *Client) Listen(clients *[]*Client) {
	for {
		str := <- thisClient.incoming
		for _, otherClient := range *clients {
			otherClient.outgoing <- str
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Printf("Error startin listen server\n")
	}
	
	clients := make([]*Client, 0)
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error establishing new connection\n")
		}
		
		client := &Client{
			incoming: make(chan string),
			outgoing: make(chan string),
		}
		clients = append(clients, client)

		go client.Listen(&clients)
		handleConnection(conn, client)
	}
}
