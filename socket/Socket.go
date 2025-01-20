package socket

import (
	"net"
	"strconv"
	"sync"
	"time"
)

type Option func (*Socket)

type Socket struct {
  mu sync.Mutex
  host string
  port int
  listener net.Listener
  Conn net.Conn
  Connections int
  EventHandler EventHandler
  IsServer bool
  bufferSize int
  Wg *sync.WaitGroup
  quit chan interface{}
  addr string
}

func NewTCPSocket (opts ...Option) *Socket {
  s := &Socket{
    host: "localhost",
    port: 7687,
    bufferSize: 1024,
    Wg: &sync.WaitGroup{},
  }
  for _, fn := range opts {
    fn(s)
  }
  return s
}

// set the host to use
func WithHost (host string) Option {
  return func(s *Socket) {
    s.host = host
  }
}

// set the port to listen/dail
func WithPort (port int) Option {
  return func(s *Socket) {
    s.port = port
  }
}

// set event handler
func WithEventHandler (ev EventHandler) Option {
  return func (s *Socket) {
    s.EventHandler = ev
  }
}

// set initial buffer size
func WithBufferSize (size int) Option {
  return func (s *Socket) {
    s.bufferSize = size
  }
}

// set wait group 
func SetWaitGroup (wg *sync.WaitGroup) Option {
  return func (s *Socket) {
    s.Wg = wg
  }
}


func (s *Socket) Listen () {
  s.IsServer = true
  addr := s.host + ":" + strconv.Itoa(s.port)
  s.addr = addr
  s.quit = make(chan interface{})
  listener, err := net.Listen("tcp", addr)

  if err != nil {
    s.EventHandler.err(err)
  }

  s.listener = listener
  s.EventHandler.open(s)

  for {
    conn, err := listener.Accept()
    if err != nil {
      // check if the error occurred because the listener was closed
      // and terminates the loop
      select {
      case <- s.quit:
        return
      default:
        s.EventHandler.err(err)
      }
    } else {
      s.Conn = conn
      s.EventHandler.connection(s)
      s.Wg.Add(1)
      go s.handleConn(conn)
    }
  }
}

// Handles the incoming connection
func (s *Socket) handleConn (conn net.Conn) {
  defer s.Wg.Done()
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
      s.EventHandler.close(s, "Connection Dropped")
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

// stops the listener from accepting new connections and 
// waits for the goroutines to end
func (s *Socket) Close () {
  close(s.quit)
  if s.IsServer {
    s.listener.Close()
  } else {
    s.Conn.Close()
  }
  s.Wg.Wait()
}

// returns the address of the listener
func (s *Socket) Address () string {
  return s.addr
}

// Connects to a TCP server.
func (s *Socket) Dial () {
  s.IsServer = false
  addr := s.host + ":" + strconv.Itoa(s.port)
  conn, err := net.Dial("tcp", addr)
  if err != nil {
    s.EventHandler.err(err)
  }
  s.Conn = conn
  s.EventHandler.open(s)

  s.quit = make(chan interface{})
  s.Wg.Add(1)
  go s.handleConn(conn)

  // waits for the quit signal
  for {
    select {
    case <- s.quit:
      return // exit

    default:
      time.Sleep (100 * time.Millisecond)
    }
  }
}


func (s *Socket) Write (data []byte) {
  _, err := s.Conn.Write(data)
  if err != nil {
    s.EventHandler.err(err)
  }
}
