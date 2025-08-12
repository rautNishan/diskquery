package connection

import (
	"fmt"
	"net"
)

/*
*

	This function will handle connection

*
*/
func HandelConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error while reading from connection: %v", err)
	}
	arrayBuffer := buffer[:n]
	fmt.Println("This is incoming data: ", string(arrayBuffer))
}
