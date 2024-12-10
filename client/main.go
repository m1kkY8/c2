package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatalf("failed getting conn %v", err)
		return
	}

	var numStrings int32
	var args []string

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

		args = append(args, string(strBuf))
	}

	fmt.Println(strings.Join(args, " "))

	execute(args)

	conn.Close()
}

func execute(args []string) {
	// Ensure there are arguments
	if len(args) == 0 {
		log.Println("no command provided")
		return
	}

	// Check if the command exists
	_, err := exec.LookPath(args[0])
	if err != nil {
		log.Printf("command not found: %v", err)
		return
	}

	// Create a command with all arguments
	cmd := exec.Command(args[0], args[1:]...) // Exclude the first argument from `args`

	// Set the command's output to standard output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr // It's a good idea to capture errors as well

	// Execute the command
	err = cmd.Run()
	if err != nil {
		log.Printf("error running command: %v", err)
		return
	}
}
