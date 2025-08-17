package connection

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const SEND_BUFFER_SIZE = 8192
const RECEVE_BUFFER_SIZE = 8192

const (
	Msg_Query   = 'Q'
	Msg_Parse   = 'P'
	Msg_Bind    = 'B'
	Msg_Execute = 'E'
)

type InputMessage struct {
	msgType byte
	data    []byte
	pos     int
}
type Connection struct {
	conn       net.Conn
	remoteAddr net.Addr
	localAddr  net.Addr
	sendBuffer []byte
	recvBuffer []byte
	reader     *bufio.Reader
	writer     *bufio.Writer
	// ctx          context.Context
	// cancel       context.CancelFunc
	tcpNoDelay   bool
	tcpKeepAlive bool
	wg           sync.WaitGroup
	backendId    int
}

/*
*

	This function will handle connection

*
*/
func HandelConnection(conn net.Conn) {
	// Initialize the port
	connection := initConnection(conn)
	if connection.remoteAddr == nil {
		log.Fatal("Failed to get remote address")
	}

	defer func() {
		log.Printf("Cleaning up connection from %s", connection.remoteAddr.String())
		connection.conn.Close()
	}()
	connection.messageLoop()
}

func initConnection(conn net.Conn) *Connection {

	port := &Connection{
		conn:       conn,
		remoteAddr: conn.RemoteAddr(),
		localAddr:  conn.LocalAddr(),
		sendBuffer: make([]byte, SEND_BUFFER_SIZE),
		recvBuffer: make([]byte, RECEVE_BUFFER_SIZE),
		reader:     bufio.NewReaderSize(conn, RECEVE_BUFFER_SIZE),
		writer:     bufio.NewWriterSize(conn, SEND_BUFFER_SIZE),
		backendId:  os.Getpid(),
	}
	err := configureTCPSocket(port)
	if err != nil {
		log.Printf("Warning: failed to configure TCP socket: %v", err)
	}

	log.Printf("Connection initialize: %s -> %s", port.remoteAddr.String(), port.localAddr.String())
	return port
}

func configureTCPSocket(connection *Connection) error {
	tcpConn, ok := connection.conn.(*net.TCPConn)
	if !ok {
		//Not a tcp connection (Maybe Unix socket), skip
		return nil
	}

	//Enable TCP_NODELAY (Disable Nagle Algorith)
	if err := tcpConn.SetNoDelay(true); err != nil {
		return fmt.Errorf("failed to set TCP keepalive: %v", err)
	}

	connection.tcpNoDelay = true

	//Enable TCP Keep Alive
	if err := tcpConn.SetKeepAlivePeriod(2 * time.Hour); err != nil {
		return fmt.Errorf("failed to set keepalive period: %v", err)
	}
	connection.tcpKeepAlive = true

	/*
		Set connection timeouts
		If there is noting to read/write in the socket then close the connection
	*/
	connection.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	connection.conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

	log.Printf("TCP socket configured: NoDelay=%v KeepAlive=%v", connection.tcpNoDelay, connection.tcpKeepAlive)
	return nil
}

func (connection *Connection) messageLoop() {
	var inputMessage InputMessage

	/*
		TODO
		1) Handle startup
		2) Send Ready for query (For client)
	*/
	for {
		firstChar, err := connection.readCommand(&inputMessage)
		if err != nil {
			if err == io.EOF {
				log.Printf("Client disconnected")
				return
			}
			log.Printf("Error reading command: %v", err)
			return
		}
		fmt.Println(firstChar == Msg_Query)
		fmt.Printf("This is query: %v", string(inputMessage.data))
		fmt.Println()
		fmt.Printf("This is pos: %d", inputMessage.pos)
	}
}

func (connection *Connection) readCommand(inputMessage *InputMessage) (int, error) {
	//Read the first byte to know the message type (Q,P and so on)
	msgType, err := connection.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return -1, err //End of file
		}
		return 0, err
	}

	//Read the message lenght From where actual query start
	var length uint32
	if err := binary.Read(connection.reader, binary.BigEndian, &length); err != nil {
		return 0, err
	}

	fmt.Println()

	//Length includes the 4 byte length field itself
	payloadSize := length - 4

	fmt.Println()

	if payloadSize > 0 {
		data := make([]byte, payloadSize)
		if _, err := io.ReadFull(connection.reader, data); err != nil {
			return 0, err
		}
		inputMessage.data = data
	} else {
		inputMessage.data = nil
	}
	inputMessage.msgType = msgType
	inputMessage.pos = 0
	return int(msgType), nil
}
