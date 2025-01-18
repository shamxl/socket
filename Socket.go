package socket

import (
	"bytes"
	"io"
	"net"
	"strconv"
)

type Option func (*Socket)

type Socket struct {
  host string
  port int
  Listener net.Listener
  Conn net.Conn
  Connections int
  EventHandler EventHandler
  isServer bool
}

func NewTCPServer (opts ...Option) *Socket {
  s := &Socket{
    host: "localhost",
    port: 7687,
    isServer: true,
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


func (s *Socket) Listen () {
  addr := s.host + ":" + strconv.Itoa(s.port)
  listener, err := net.Listen("tcp", addr)

  if err != nil {
    s.EventHandler.Error(err)
  }

  s.Listener = listener
  s.EventHandler.Open(s)

  for {
    conn, err := listener.Accept()
    defer conn.Close()
    if err != nil {
      s.EventHandler.Error(err)
    }
    s.EventHandler.Connection()

    go s.handleConn(conn)
  }
}

func (s *Socket) handleConn (conn net.Conn) {
  var buffer bytes.Buffer

  for {
    _, err := io.Copy(&buffer, conn) 
    s.EventHandler.Data(s, buffer.Bytes())
    if err != nil {
      if err.Error() != "EOF" {
        s.EventHandler.Error(err)
      }
      s.EventHandler.Close(s, "Client Disconnected")
      break
    }

  }
}
