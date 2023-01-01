//go:build js && wasm
// +build js,wasm

package messages

import (
	"github.com/Nigel2392/jsext"
)

// Sends messages through eventlisteners.
// These can be accessed though window.runtime.eventOn("jsextMessages", callback(arg: function(type, message)))
// Or optionally when using the jsext tool, you can use "window.runtime.onMessage(callback(arg: function(type, message)))"
// go get github.com/Nigel2392/jsexttool/

type MessageType string

const (
	Info    MessageType = "info"
	Success MessageType = "success"
	Warning MessageType = "warning"
	Error   MessageType = "error"
)

func SendMessage(typ MessageType, message string) {
	jsext.EventEmit("jsextMessages", string(typ), message)
}

func SendInfo(message string) {
	SendMessage(Info, message)
}

func SendSuccess(message string) {
	SendMessage(Success, message)
}

func SendWarning(message string) {
	SendMessage(Warning, message)
}

func SendError(message string) {
	SendMessage(Error, message)
}

func Listen(callback func(typ string, message string)) {
	jsext.EventOn("jsextMessages", func(args ...interface{}) {
		var typStr string = args[0].(string)
		var message string = args[1].(string)
		callback(typStr, message)
	})
}

func ListenFor(typ MessageType, callback func(message string)) {
	Listen(func(typ, message string) {
		if typ == string(typ) {
			callback(message)
		}
	})
}

func ListenForInfo(callback func(message string)) {
	ListenFor(Info, callback)
}

func ListenForSuccess(callback func(message string)) {
	ListenFor(Success, callback)
}

func ListenForWarning(callback func(message string)) {
	ListenFor(Warning, callback)
}

func ListenForError(callback func(message string)) {
	ListenFor(Error, callback)
}
