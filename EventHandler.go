package socket

type Args struct {
  Socket *Socket
  Data []byte
  Reason string
  ErrorMsg error
}


type OnOpenFn func (Args)
type OnDataFn func (Args)
type OnCloseFn func (Args)
type OnErrorFn func (Args)
type OnConnectionFn func ()
type OnDisconnectFn func ()

type EventHandler struct {
  OnConnection OnConnectionFn
  OnOpen OnOpenFn
  OnData OnDataFn
  OnClose OnCloseFn
  OnError OnErrorFn
}

func (ev *EventHandler) SetOnConnection (fn OnConnectionFn) {
  ev.OnConnection = fn
}

func (ev *EventHandler) SetOnOpen (fn OnOpenFn) {
  ev.OnOpen = fn
}

func (ev *EventHandler) SetOnData (fn OnDataFn) {
  ev.OnData = fn
}

func (ev *EventHandler) SetOnClose (fn OnCloseFn) {
  ev.OnClose = fn
}

func (ev *EventHandler) SetOnError (fn OnErrorFn) {
  ev.OnError = fn
}


func (ev *EventHandler) Connection () {
  if ev.OnConnection != nil {
    ev.OnConnection()
  }
}

func (ev *EventHandler) Open (sock *Socket) {
  if ev.OnOpen != nil {
    ev.OnOpen(Args{
      Socket: sock,
    })
  }
}

func (ev *EventHandler) Data (sock *Socket, data []byte) {
  if ev.OnData != nil {
    ev.OnData(Args{
      Socket: sock,
      Data: data,
    })
  }
}

func (ev *EventHandler) Close (sock *Socket, reason string) {
  if ev.OnClose != nil {
    ev.OnClose(Args{
      Socket: sock,
      Reason: reason,
    })
  }
}

func (ev *EventHandler) Error (err error) {
  if ev.OnError != nil {
    ev.OnError(Args{
      ErrorMsg: err,
    })
  }
}
