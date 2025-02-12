# Go Reverse SHELL with TLS Encryption

This project consists of a **server** and a **client** written in Go that allow remote command execution over a **TLS-encrypted** connection. The server sends commands to the client, and the client executes them locally before sending back the results.

## Features
- Secure communication using **TLS encryption**.
- Executes **commands remotely** on the client machine.
- Sends **multi-line command output** back to the server.
- Supports **Linux/macOS** (for Windows, adjust shell execution accordingly).

## Prerequisites
Before running the server and client, you need to generate a **TLS certificate**.

### Generate a TLS Certificate
Run the following command to generate a self-signed certificate:
```sh
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```
This will generate:
- `cert.pem` - TLS certificate.
- `key.pem` - Private key.

Ensure that **both the server and client have access to these files**.

---

## Installation & Usage

### 1. Clone the Repository
```sh
git clone https://github.com/your-repo/remote-command-tls.git
cd remote-command-tls
```

### 2. Build the Server and Client
```sh
go build -o server server.go
go build -o client client.go
```

### 3. Start the Server
Run the server on the **host machine**:
```sh
./server
```

### 4. Start the Client
Run the client on the **remote machine** (or another terminal):
```sh
./client
```

### 5. Execute Commands
Once the client is connected, type any command in the **server terminal**, and it will execute on the **client machine**. The result will be sent back to the server.

Example commands:
```sh
ls -la
ps aux
df -h
```

---

## Code Overview

### Server (`server.go`)
- Listens on port **8080** over a **TLS-encrypted** connection.
- Waits for client connection.
- Sends commands to the client.
- Receives and prints command output from the client.

### Client (`client.go`)
- Connects to the server over **TLS**.
- Waits for commands.
- Executes commands locally.
- Sends the output (including multi-line responses) back to the server.

---

## Security Notes
- This implementation **disables certificate verification on the client** (`InsecureSkipVerify: true`).
  - This should be replaced with proper **certificate verification** in a production environment.
- Using self-signed certificates is acceptable for local use, but **a valid CA-signed certificate is recommended** for real-world deployment.

---

## Troubleshooting

### Error: `Error loading certificate`
- Ensure `cert.pem` and `key.pem` are in the same directory as `server.go`.

### Error: `connection refused`
- Ensure the server is **running** before starting the client.
- Check **firewall settings** to allow connections on port **8080**.

### Error: `command not found`
- Ensure the command exists on the client machine.
- Try running it manually on the client before executing remotely.

---

## License
This project is licensed under the MIT License. Feel free to modify and use it in your own projects.
