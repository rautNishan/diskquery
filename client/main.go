package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Running client")
	conn, _ := net.Dial("tcp", "localhost:3000")

	queries := []string{"SELECT * FROM users"}
	buf := bytes.NewBuffer([]byte{})

	for _, q := range queries {
		fmt.Println("This is q: ", q)
		fmt.Println()
		payload := append([]byte(q), 0)
		fmt.Printf("Thi is payload %v", string(payload))
		fmt.Println()
		length := uint32(len(payload) + 4)
		buf.WriteByte('Q')
		binary.Write(buf, binary.BigEndian, length)
		buf.Write(payload)
	}
	conn.Write(buf.Bytes())
}
