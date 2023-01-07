# Message Support from WASM to Javascript.

This example shows how to send messages from WASM to Javascript, and how to listen for these messages.

The raw way to interop between JSExt and Javascript is with eventListeners.
We The arguments passed to the eventListener, are passed on to the `event.args` as a list.

We will show this using our embedded message api as an example.
We do provide a wrapper for this message api, to make it easier to use, but we provide the full information on how this wrapper is made to give an insight on how the message api works.

**Warning:** The Javascript examples only works using the jsexttool, without the tool, the function `window.jsextLoaded.On()` is not available. The tool provides an initialization script to simplify loading, and do some setup for the application.

Here is how to listen for events on Javascript:
```js
// Wait for wasm module to be initialized.
window.jsextLoaded.On(function(){
	jsext.runtime.eventOn("jsextMessages", function(e){console.log("JSEvent, Type:", e.args[0], " Message:", e.args[1])})
})
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
// Wait for wasm module to be initialized.
window.jsextLoaded.On(function(){
	jsext.runtime.eventEmit("jsextMessages", "error", "My custom error message!")
})
```
And in GO:
```go
jsext.EventEmit("jsextMessages", "error", "My custom error message!")
```

Now that we know how to listen for events, and how to emit events, we can use the message api to make it easier to use.
As stated before, this is just a wrapper around the events api to abstract away some of the complexity.

First, again, we will show you how to send messages.
To use the message api from javascript, we can use the following code:
```js
// Wait for wasm module to be initialized.
window.jsextLoaded.On(function(){
	jsext.runtime.sendMessage("error","My custom error message!")
})
```

And in GO:
```go
messages.SendMessage(messages.Error, "My custom error message!")
// Optionally, use: messages.SendError("My custom error message!")
```

To receive messages, we can use the following code:
```js
// Wait for wasm module to be initialized.
window.jsextLoaded.On(function(){
    // Define what happens when messages are sent from the WASM module.
    // Function needs to be set on the window object and take two arguments.
    jsext.runtime.onMessage(function(typ, message) {
		console.log("JSEvent, Type:", typ, " Message:", message)
    });
})
```
And in GO:
```go
messages.Listen(func(typ string, message string) {
		println("WASMEvent, Type:", typ, " Message:", message)
})
```
