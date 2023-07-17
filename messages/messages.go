//go:build js && wasm
// +build js,wasm

package messages

import (
	"github.com/Nigel2392/jsext/v2"
)

// Sends messages through eventlisteners.
// These can be accessed though window.runtime.eventOn("jsextMessages", callback(arg: function(type, message)))
// The type is a string and the message is a string.

type MessageType string

const (
	TypeInfo    MessageType = "info"
	TypeSuccess MessageType = "success"
	TypeWarning MessageType = "warning"
	TypeError   MessageType = "error"
)

func Emit(typ MessageType, message string) {
	jsext.EventEmit("jsextMessages", string(typ), message)
}

func Info(message string) {
	Emit(TypeInfo, message)
}

func Success(message string) {
	Emit(TypeSuccess, message)
}

func Warning(message string) {
	Emit(TypeWarning, message)
}

func Error(message string) {
	Emit(TypeError, message)
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
