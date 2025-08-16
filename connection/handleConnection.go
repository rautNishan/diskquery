package connection

import (
	"fmt"
	"log"
	"net"
	"time"
)

const SEND_BUFFER_SIZE = 8192
const RECEVE_BUFFER_SIZE = 8192

type Port struct {
	Conn       net.Conn
	RemoteAddr net.Addr
	LocalAddr  net.Addr
	SendBuffer []byte
	RecvBuffer []byte

	//Buffer state tracking
	SendPointer int //Current position in send buffer
	SendStart   int //Start of unsent data
	RecvPointer int //Current position in receive buffer
	RecvLength  int //Amount of data in receive buffer

	//Communication state
	CommBusy       bool //Busy sending data
	ReadingMessage bool //In middle of reading a message

	//Connection Properties
	IsNonBlocking bool
	TCPNoDelay    bool
	TCPKeepAlive  bool

	//Protocol version /state
	ProtocolVersion int
	SSLActive       bool
	GSSActive       bool

	//Session State
	DatabaseName string
	UserName     string

	ReadDeadLine  time.Time
	WriteDeadLine time.Time

	BackendPID int
}

/*
*

	This function will handle connection

*
*/
func HandelConnection(conn net.Conn) {
	// Initialize the port
	port := init_port(conn)
	if port.RemoteAddr == nil {
		log.Fatal("Failed to get remote address")
	}
	fmt.Print("Port is initialize")
}

func init_port(conn net.Conn) *Port {
	port := &Port{
		Conn:           conn,
		RemoteAddr:     conn.RemoteAddr(),
		LocalAddr:      conn.LocalAddr(),
		SendBuffer:     make([]byte, SEND_BUFFER_SIZE),
		RecvBuffer:     make([]byte, RECEVE_BUFFER_SIZE),
		SSLActive:      false,
		GSSActive:      false,
		BackendPID:     0,
		SendPointer:    0,
		SendStart:      0,
		RecvPointer:    0,
		RecvLength:     0,
		CommBusy:       false,
		ReadingMessage: false,
		IsNonBlocking:  false,
	}
	err := configureTCPSocket(port)
	if err != nil {
		log.Printf("Warning: failed to configure TCP socket: %v", err)
	}
	return port
}

func configureTCPSocket(port *Port) error {
	tcpConn, ok := port.Conn.(*net.TCPConn)
	if !ok {
		//Not a tcp connection (Maybe Unix socket), skip
		return nil
	}

	//Enable TCP_NODELAY (Disable Nagle Algorith)
	if err := tcpConn.SetNoDelay(true); err != nil {
		return fmt.Errorf("failed to set TCP keepalive: %v", err)
	}

	port.TCPNoDelay = true

	//Enable TCP Keep Alive
	if err := tcpConn.SetKeepAlivePeriod(2 * time.Hour); err != nil {
		return fmt.Errorf("failed to set keepalive period: %v", err)
	}
	port.TCPKeepAlive = true

	log.Printf("TCP socket configured: NoDelay=%v KeepAlive=%v", port.TCPNoDelay, port.TCPKeepAlive)
	return nil
}
