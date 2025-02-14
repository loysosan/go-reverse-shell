package main

import (
	"bufio"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	clients   = make(map[int]net.Conn)
	clientIDs = 0
	mu        sync.Mutex
)

func handleClient(conn net.Conn, clientID int) {
	defer conn.Close()

	clientRemoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Printf("Client %d connected: %s\n", clientID, clientRemoteAddr)

	for {
		// Читаем длину ответа
		var length int32
		err := binary.Read(conn, binary.LittleEndian, &length)
		if err != nil {
			fmt.Printf("Client %d disconnected\n", clientID)
			mu.Lock()
			delete(clients, clientID)
			mu.Unlock()
			return
		}

		// Читаем сам ответ
		response := make([]byte, length)
		_, err = conn.Read(response)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Printf("Client %d response:\n%s\n", clientID, string(response))
	}
}

func main() {
	// Загружаем TLS сертификат
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	// Настройка TLS сервера
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("TLS server listening on port 8080...")

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Connection error:", err)
				continue
			}

			mu.Lock()
			clientID := clientIDs
			clients[clientID] = conn
			clientIDs++
			mu.Unlock()

			go handleClient(conn, clientID)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	var activeClient int

	for {
		fmt.Println("\nConnected clients:")
		mu.Lock()
		for id := range clients {
			fmt.Printf("Client %d\n", id)
		}
		mu.Unlock()

		fmt.Print("\nEnter client ID to interact with: ")
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "%d", &activeClient)

		mu.Lock()
		client, exists := clients[activeClient]
		mu.Unlock()

		if !exists {
			fmt.Println("Invalid client ID.")
			continue
		}

		for {
			fmt.Printf("Enter command to send to Client %d (or type 'exit' to switch): ", activeClient)
			scanner.Scan()
			command := scanner.Text()

			if command == "exit" {
				break
			}

			_, err := client.Write([]byte(command + "\n"))
			if err != nil {
				fmt.Println("Error sending command:", err)
				break
			}
		}
	}
}