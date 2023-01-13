# Simple example of a websocket connection.

This example shows how to use the websocket module to create a simple websocket client.

```go
Websocket, err := websocket.Open("ws://127.0.0.1:8000/ws")
if err != nil {
	messages.SendError(err.Error())
	os.Exit(1)
}

Websocket := websocket.New("ws://127.0.0.1:8000/ws")
if err != nil {
	messages.SendError(err.Error())
	os.Exit(1)
}


Websocket.OnMessage(func(event websocket.MessageEvent) {
	messages.SendInfo(event.Data().String())
})

Websocket.OnClose(func(event js.Value) {
	messages.SendWarning("Websocket closed")
})

go func() {
	for {
		Websocket.Send("Hello World!")
		time.Sleep(1 * time.Second)
	}
}()
```