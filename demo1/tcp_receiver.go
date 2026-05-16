package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error: ", err)
	}
	defer listener.Close()
	fmt.Printf("Listening on port 8080 \n")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %v \n", err)
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
		fmt.Printf("Eror: %v \n", err)
		return
	}
	msg := string(buffer[:n])
	fmt.Printf("Message: %s \n", msg)
}
