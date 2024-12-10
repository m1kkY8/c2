package main

import (
	"encoding/binary"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatalf("failed getting conn %v", err)
		return
	}

	var numStrings int32

	// Read the number of strings
	if err := binary.Read(conn, binary.LittleEndian, &numStrings); err != nil {
		log.Printf("error reading number of strings: %v", err)
		conn.Close()
		return
	}

	// Read each string
	for i := 0; i < int(numStrings); i++ {
		var strLen int32
		// Read the length of the string
		if err := binary.Read(conn, binary.LittleEndian, &strLen); err != nil {
			log.Printf("error reading string length: %v", err)
			conn.Close()
			return
		}

		// Read the actual string data
		strBuf := make([]byte, strLen)
		if err := binary.Read(conn, binary.LittleEndian, &strBuf); err != nil {
			log.Printf("error reading string: %v", err)
			conn.Close()
			return
		}

		// Convert to string and print
		log.Println(string(strBuf))
	}

	conn.Close()
}
