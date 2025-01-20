package socket_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/shamxl/socket/socket"
)

func onOpen (args socket.Args) {
  fmt.Println("Connection open!")

  if !args.Socket.IsServer {
    args.Socket.Conn.Write([]byte("Hello, Server"))
  }
}
func onClose (args socket.Args) {
  fmt.Println("Connection closed")
}
func onError (args socket.Args) {
  fmt.Println("Error", args.ErrorMsg)
  os.Exit(1)
}
func onData (args socket.Args) {
  fmt.Println("Got data", args.Data)

  if args.Socket.IsServer {
    fmt.Println("Server")
    if string(args.Data) == "Goodbye!\n" {
      args.Socket.Write([]byte("Goodbye!\n"))
      args.Socket.Close()
    }
  } else {
    fmt.Println("Server")
    if string(args.Data) == "Goodbye!\n" {
      fmt.Println("Closing server...")
      args.Socket.Write([]byte("Goodbye!\n"))
      args.Socket.Close()
    }
  }
}
func onConn (args socket.Args) {
  fmt.Println("New connnnn")
  args.Socket.Conn.Write([]byte("Hello, Client"))
}

func TestMain (t *testing.T) {
  eventHandler := socket.EventHandler{}
  eventHandler.SetOnOpen(onOpen)
  eventHandler.SetOnClose(onClose)
  eventHandler.SetOnConnection(onConn)
  eventHandler.SetOnData(onData)


  serverSock := socket.NewTCPSocket (
    socket.WithEventHandler(eventHandler),
    socket.WithPort(8080),
    socket.WithBufferSize(1024),
  )
  clientSock := socket.NewTCPSocket(
    socket.WithEventHandler(eventHandler),
    socket.WithPort(8080),
    socket.WithBufferSize(1024),
  )
  
  var wg sync.WaitGroup
  wg.Add (2)
  go func (serverSock *socket.Socket) {
    defer wg.Done()
    serverSock.Listen()
  }(serverSock)


  go func (clientSock *socket.Socket) {
    defer wg.Done()
    clientSock.Dial()
  }(clientSock)

  wg.Wait()
}
