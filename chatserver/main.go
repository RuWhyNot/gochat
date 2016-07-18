package main

import "fmt"
import "net"
import "bufio"
import "os"


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

func AcceptClients(ln net.Listener, clients *[]*Client) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error establishing new connection\n")
			continue
		}
		
		client := &Client{
			incoming: make(chan string),
			outgoing: make(chan string),
		}
		*clients = append(*clients, client)

		go client.Listen(clients)
		handleConnection(conn, client)
	}
}

func DoCommand(comm string) {
	switch comm {
	case "stop":
		os.Exit(0)
	default:
		fmt.Printf("error: Unknown command \"%s\"\n", comm)
	}
}

func main() {
	serverPort := "8081"
	if len(os.Args) == 2 {
		serverPort = os.Args[1]
	}

	ln, err := net.Listen("tcp", ":" + serverPort)
	if err != nil {
		fmt.Printf("Error startin listen server\n")
		os.Exit(1)
	}

	fmt.Printf("Server runned on port %s\n", serverPort)
	
	clients := make([]*Client, 0)
	
	go AcceptClients(ln, &clients)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		DoCommand(text[:len(text)-2])
	}
}
