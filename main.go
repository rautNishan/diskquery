package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rautNishan/diskquery/connection"
)

// Accept tcp connection
func main() {
	listner, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("Error while starting server")
	}
	fmt.Println("Listning on port 3000")
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("Error while accepting connection: %v", err) //We do not want to shut down our program
		}
		go connection.HandelConnection(conn)
	}
}
