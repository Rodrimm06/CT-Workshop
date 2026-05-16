package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func main() {
	Cert, _ := os.ReadFile("prueba_cert.crt")
	CertPool := x509.NewCertPool()
	CertPool.AppendCertsFromPEM(Cert)

	config := &tls.Config{

		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
		RootCAs:    CertPool,
		ServerName: "localhost",
	}

	conn, err := tls.Dial("tcp", "localhost:8080", config)
	if err != nil {
		fmt.Println("Error: %v \n", err)
		os.Exit(1)
	}
	defer conn.Close()
	msg := "Hello TLS"

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}

	fmt.Println("Message Sent")
}
