package socket

import (
	"net"
	"strconv"
	"sync"
)

type Option func (*Socket)

type Socket struct {
  mu sync.Mutex
  host string
  port int
  Listener net.Listener
  Conn net.Conn
  Connections int
  EventHandler EventHandler
  isServer bool
  bufferSize int
}

func NewTCPServer (opts ...Option) *Socket {
  s := &Socket{
    host: "localhost",
    port: 7687,
    isServer: true,
    bufferSize: 1024,
  }
  for _, fn := range opts {
    fn(s)
  }
  return s
}

func WithHost (host string) Option {
  return func(s *Socket) {
    s.host = host
  }
}

func WithPort (port int) Option {
  return func(s *Socket) {
    s.port = port
  }
}

func WithEventHandler (ev EventHandler) Option {
  return func (s *Socket) {
    s.EventHandler = ev
  }
}

func WithBufferSize (size int) Option {
  return func (s *Socket) {
    s.bufferSize = size
  }
}


func (s *Socket) Listen () {
  addr := s.host + ":" + strconv.Itoa(s.port)
  listener, err := net.Listen("tcp", addr)

  if err != nil {
    s.EventHandler.err(err)
  }

  s.Listener = listener
  s.EventHandler.open(s)

  for {
    conn, err := listener.Accept()
    if err != nil {
      s.EventHandler.err(err)
    }
    s.EventHandler.connection()

    go s.handleConn(conn)
  }
}

// Handles the incoming connection
func (s *Socket) handleConn (conn net.Conn) {
  defer conn.Close()
  // initalizing buffer with the initial buffer size set by the user/default
  s.mu.Lock()
  var buffer = make ([]byte, s.bufferSize)
  s.mu.Unlock()
  for {
    // reading buffer size to check the changes in buffer size
    s.mu.Lock()
    bufferSize := s.bufferSize
    s.mu.Unlock()

    // checks and re-initializes array with the specified array size if changes are detected
    if len(buffer) != bufferSize {
      buffer = make([]byte, bufferSize)
    }

    n, err := conn.Read(buffer)
    if err != nil {
      if err.Error() != "EOF" {
        s.EventHandler.err(err)
      }
      s.EventHandler.close(s, "Client Disconnected")
      break
    }
    
    // calls the data handler only if theres any data
    // to avoid calling the data handler forever
    if n > 0 {
      s.EventHandler.data(s, buffer)
    }
    
  }
}

// Switches the buffer size to specified length
// to recieve the data in full piece instead of chunks.
// Use this to avoid breaking your protocol.
func (s *Socket) SwitchBufferSize (size int) {
  s.mu.Lock()
  s.bufferSize = size
  s.mu.Unlock()
}


