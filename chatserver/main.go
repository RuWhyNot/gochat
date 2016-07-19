package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"strings"
)

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
		fmt.Printf("Error starting listen server\n")
		os.Exit(1)
	}

	fmt.Printf("Server runned on port %s\n", serverPort)
	
	clients := make([]*Client, 0)
	
	go AcceptClients(ln, &clients)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		DoCommand(strings.TrimSpace(text))
	}
}
