package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rautNishan/diskquery/connection"
)

/*
We will accept connection and do database thing per thread
Now we can later do async i/o per each thread like postgres does for each processes
Since we are not spinning up backend process (Like postgres does), for each thread we can have async io
Not for the concurrency (Which is alread handle whele calling connecion.HandleConnection) but rather to listen events
Such as when client exit its query press ctrl+c or want to exit middle of the query
*/
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
