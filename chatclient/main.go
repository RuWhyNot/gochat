package main

import "fmt"
import "net"
import "bufio"
import "os"

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Printf("Can't find server\n")
		os.Exit(0)
	}

	conn.Write([]byte("Test message\n"))

	for {
		str, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Connection lost\n")
			os.Exit(0)
		}
		fmt.Printf(str)
	}
}