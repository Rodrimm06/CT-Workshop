package main

import (
	"crypto/tls"
	"fmt"
	"net"
)

func main() {

	cert, err := tls.LoadX509KeyPair("prueba_cert.crt", "server.key")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12}

	listener, err := tls.Listen("tcp", ":8080", config)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return
	}

	defer listener.Close()
	fmt.Printf("Listening on port 8080 \n")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Eror: %v", err)
		return
	}
	msg := string(buffer[:n])
	fmt.Printf("Message: %s \n", msg)
}
