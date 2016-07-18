package main

import "fmt"
import "net"
import "bufio"
import "os"

func Read(conn net.Conn) {
	for {
		str, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Connection lost\n")
			os.Exit(1)
		}
		fmt.Printf(">")
		fmt.Printf(str)
	}
}

func DoCommand(comm string) {
	switch comm {
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("error: Unknown command \"%s\"\n", comm)
	}
}

func main() {
	serverAddress := "localhost:8081"
	if len(os.Args) == 2 {
		serverAddress = os.Args[1]
	}


	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Printf("Can't find server\n")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	go Read(conn)

	for {
		text, _ := reader.ReadString('\n')

		if (text[:1] == "/") {
			// without "/" and "\n"
			DoCommand(text[1:len(text)-2])
			continue
		}

		conn.Write([]byte(text))
	}
}