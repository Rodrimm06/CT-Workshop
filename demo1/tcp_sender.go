package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	msg := "Hello TCP"

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error: %v", err)
		os.Exit(1)
	}

	fmt.Println("Mensaje enviado")
}
