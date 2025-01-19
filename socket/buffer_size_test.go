package socket_test

import (
	"testing"

	"github.com/shamxl/socket/socket"
)

func TestMain (t *testing.T) {
  bufferSizeSwitched := false
  bufferSizeToSwitch := 1024
  ev := socket.EventHandler{}
  ev.SetOnOpen(func (args socket.Args) {
    t.Log("Socket is open, try sending some data")
  })
  ev.SetOnError(func (args socket.Args) {
    t.Error (args.ErrorMsg)
  })
  ev.SetOnData(func (args socket.Args) {
    t.Log("Recieved Data", args.Data)
    if bufferSizeSwitched {
      if len(args.Data) == bufferSizeToSwitch {
        t.Log("Test Passed") 
      } else {
        t.Fail()
      }
    } else {
      t.Log("Switching buffer size to ", bufferSizeToSwitch)
      args.Socket.SwitchBufferSize(bufferSizeToSwitch)
      bufferSizeSwitched = true
      t.Log("Send some data")
    }
  })

  sock := socket.NewTCPServer(socket.WithEventHandler(ev), socket.WithBufferSize(3))
  sock.Listen()

}
