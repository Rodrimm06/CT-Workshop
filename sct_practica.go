package main

import (
	"crypto/tls"
	"encoding/asn1"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

// Represent a Certificate log
type CTLog struct {
	Name string
	URL  string
}

// Hash table with the hash sha-256 of the log's key and relevant information
var knownLogs = map[string]CTLog{
	"d809553b944f7affc816196f944f85abb0f8fc5e8755260f15d12e72bb454b14": {"Google 'Argon2026' Log", "ct.googleapis.com/logs/argon2026"},
	"cb38f715897c84a1445f5bc1ddfbc96ef29a59cd470a690585b0cb14c31458e7": {"DigiCert Nessie Log", "ct.digicert.com/log/nessie"},
	"41c8ca980744e2b5f6b64f9b8c3d8c1c302fa2d9d1091b6a0f4439c368249080": {"Cloudflare 'Nimbus' Log", "ct.cloudflare.com/logs/nimbus"},
}

func main() {
	// Declare the target
	target := "google.com:443"
	fmt.Printf("Target: %s\n", target)

	// Call tls to target
	conf := &tls.Config{InsecureSkipVerify: false}
	conn, err := tls.Dial("tcp", target, conf)
	if err != nil {
		log.Fatalf("Error dialing target: %v", err)
	}
	defer conn.Close() // Close communication

	// Retrieve certificate
	cert := conn.ConnectionState().PeerCertificates[0]

	// Extention for the id of sct
	sctOID := "1.3.6.1.4.1.11129.2.4.2"
	foundSCTs := 0

	for _, ext := range cert.Extensions {
		if ext.Id.String() == sctOID { // Search the extention
			var sctListDer []byte
			if _, err := asn1.Unmarshal(ext.Value, &sctListDer); err != nil { // get the raw bytes of the sct list
				continue
			}

			offset := 2 // first two bytes are the length of the list
			for offset < len(sctListDer) {
				// 	Length of each sct
				sctLen := int(binary.BigEndian.Uint16(sctListDer[offset : offset+2]))
				//  Move two units because we know the length of the sct
				offset += 2
				// Slice the bytes to get exactly 1 sct
				sctBytes := sctListDer[offset : offset+sctLen]
				// Call function with the sct
				processSCT(sctBytes, cert.NotBefore)

				offset += sctLen // Advance to process the next SCT
				foundSCTs++      // Increase the count for the SCT
			}
		}
	}

	if foundSCTs == 0 { // If there is no SCT found
		fmt.Println("\nNo SCT found")
	}
}

// Turn raw bytes to information
func processSCT(b []byte, certNotBefore time.Time) {
	// Structure:
	// Version: 1 byte (b[0])
	// Log ID: 32 bytes (b[1:33])
	// Timestamp: 8 bytes (b[33:41])

	logIDHex := fmt.Sprintf("%x", b[1:33])                       // Get the SHA-256 of the public key log & convert it to hex
	timestampMS := int64(binary.BigEndian.Uint64(b[33:41]))      // Get the timestamp & build a 64 bit number
	sctTime := time.Unix(0, timestampMS*int64(time.Millisecond)) // Turn the timestamp into a date

	fmt.Printf("\nSCT\n")

	// Seasrch in the hash table
	if info, ok := knownLogs[logIDHex]; ok {
		fmt.Printf("Log Recognized\n")
		fmt.Printf("  Server: %s\n", info.Name)
		fmt.Printf("  URL:      %s\n", info.URL)

	} else {
		fmt.Printf("Log Unrecognized\n")
		fmt.Printf("LogID: %s\n", logIDHex[:16])
	}

	fmt.Printf("Date: %s\n", sctTime.Format(time.RFC1123))

}
