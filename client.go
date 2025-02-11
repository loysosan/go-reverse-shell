package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
)

const serverAddress = "127.0.0.1:9000"

func main() {
	// Configure TLS (disable verification for self-signed certificates)
	config := &tls.Config{InsecureSkipVerify: true}

	// Connect to the TLS server
	conn, err := tls.Dial("tcp", serverAddress, config)
	if err != nil {
		fmt.Println("Server connection ERROR:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the TLS server. Please enter command:")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(command))
		if err != nil {
			fmt.Println("Error sending data:", err)
			break
		}

		// Read server response
		response := make([]byte, 4096)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error receiving response:", err)
			break
		}

		fmt.Println(string(response[:n]))
	}
}