# Message Support from WASM to Javascript.

This example shows how to send messages from WASM to Javascript, and how to listen for these messages.

The raw way to interop between JSExt and Javascript is with eventListeners.
We The arguments passed to the eventListener, are passed on to the `event.args` as a list.

We will show this using our embedded message api as an example.
We do provide a wrapper for this message api, to make it easier to use, but we provide the full information on how this wrapper is made to give an insight on how the message api works.

Here is how to listen for events on Javascript:
```js
jsext.runtime.eventOn("jsextMessages", function(e){console.log("JSEvent, Type:", e.args[0], " Message:", e.args[1])})
``` 

And the equivalent in GO:
```go
jsext.EventOn("jsextMessages", func(args ...interface{}) {
		var typStr string = args[0].(string)
		var message string = args[1].(string)
		println("WASMEvent, Type:", typStr, " Message:", message)
})
```

To emit events, messages in this case, we can use the following code in javascript.
```js
jsext.runtime.eventEmit("jsextMessages", "error", "My custom error message!")
```
And in GO:
```go
jsext.EventEmit("jsextMessages", "error", "My custom error message!")
```