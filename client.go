package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const serverAddress = "127.0.0.1:9000"

func main() {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Server connection ERROR:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Please enter command:")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(command))
		if err != nil {
			fmt.Println("Error data sending:", err)
			break
		}

		// Red server responce
		response := make([]byte, 4096)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error reciving responce", err)
			break
		}

		fmt.Println(string(response[:n]))
	}
}