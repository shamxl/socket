package socket

type Args struct {
  Socket *Socket
  Data []byte
  Reason string
  ErrorMsg error
}


type onOpenFn        func (Args)
type onDataFn        func (Args)
type onCloseFn       func (Args)
type onErrorFn       func (Args)
type onConnectionFn  func (Args)
type onDisconnectFn  func ()

type EventHandler struct {
  OnConnection onConnectionFn
  OnOpen onOpenFn
  OnData onDataFn
  OnClose onCloseFn
  OnError onErrorFn
}

func (ev *EventHandler) SetOnConnection (fn onConnectionFn) {
  ev.OnConnection = fn
}

func (ev *EventHandler) SetOnOpen (fn onOpenFn) {
  ev.OnOpen = fn
}

func (ev *EventHandler) SetOnData (fn onDataFn) {
  ev.OnData = fn
}

func (ev *EventHandler) SetOnClose (fn onCloseFn) {
  ev.OnClose = fn
}

func (ev *EventHandler) SetOnError (fn onErrorFn) {
  ev.OnError = fn
}


func (ev *EventHandler) connection (sock *Socket) {
  if ev.OnConnection != nil {
    ev.OnConnection(Args{
      Socket: sock,
    })
  }
}

func (ev *EventHandler) open (sock *Socket) {
  if ev.OnOpen != nil {
    ev.OnOpen(Args{
      Socket: sock,
    })
  }
}

func (ev *EventHandler) data (sock *Socket, data []byte) {
  if ev.OnData != nil {
    ev.OnData(Args{
      Socket: sock,
      Data: data,
    })
  }
}

func (ev *EventHandler) close (sock *Socket, reason string) {
  if ev.OnClose != nil {
    ev.OnClose(Args{
      Socket: sock,
      Reason: reason,
    })
  }
}

func (ev *EventHandler) err (err error) {
  if ev.OnError != nil {
    ev.OnError(Args{
      ErrorMsg: err,
    })
  }
}
