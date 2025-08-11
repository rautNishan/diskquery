package main

import (
	"fmt"
	"log"
	"net"
)

// Accept tcp connection

func main() {
	listner, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Error while starting server")
	}
	fmt.Println("Listning on port 3000")
	conn, err := listner.Accept()
	if err != nil {
		log.Fatal("Errpr while accepting connection")
	}
	fmt.Println("This is connection file descriptor: ", &conn)
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal("Error while reading from the connection")
	}
	parsedData := string(buffer[:n])
	fmt.Println(parsedData)

	_, err = conn.Write([]byte(parsedData))
	if err != nil {
		log.Fatal("Error while writing to the connection")
	}
	fmt.Println("Closing the server")
}
