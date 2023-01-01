//go:build js && wasm
// +build js,wasm

package console

import "syscall/js"

// MISSING:
// console.table()

// Javascript Console.assert() method
func Assert(condition bool, args ...interface{}) {
	if condition {
		js.Global().Get("console").Call("assert", append([]any{condition}, args...)...)
		return
	}
}

// Javascript Console.clear() method
func Clear() {
	js.Global().Get("console").Call("clear")
}

// Javascript Console.count() method
func Count(label string) {
	js.Global().Get("console").Call("count", label)
}

// Javascript Console.countReset() method
func CountReset(label string) {
	js.Global().Get("console").Call("countReset", label)
}

// Javascript Console.debug() method
func Debug(args ...interface{}) {
	js.Global().Get("console").Call("debug", args...)
}

// Javascript Console.dir() method
func Dir(item interface{}) {
	js.Global().Get("console").Call("dir", item)
}

// Javascript Console.dirxml() method
func DirXML(item interface{}) {
	js.Global().Get("console").Call("dirxml", item)
}

// Javascript Console.error() method
func Error(args ...interface{}) {
	js.Global().Get("console").Call("error", args...)
}

// Javascript Console.exception() method
func Group(args ...interface{}) {
	js.Global().Get("console").Call("group", args...)
}

// Javascript Console.groupCollapsed() method
func GroupCollapsed(args ...interface{}) {
	js.Global().Get("console").Call("groupCollapsed", args...)
}

// Javascript Console.groupEnd() method
func GroupEnd() {
	js.Global().Get("console").Call("groupEnd")
}

// Javascript Console.info() method
func Info(args ...interface{}) {
	js.Global().Get("console").Call("info", args...)
}

// Javascript Console.log() method
func Log(args ...interface{}) {
	js.Global().Get("console").Call("log", args...)
}

// Javascript Console.time() method
func Time(label string) {
	js.Global().Get("console").Call("time", label)
}

// Javascript Console.timeEnd() method
func TimeEnd(label string) {
	js.Global().Get("console").Call("timeEnd", label)
}

// Javascript Console.timeLog() method
func TimeLog(label string, args ...interface{}) {
	js.Global().Get("console").Call("timeLog", append([]any{label}, args...)...)
}

// Javascript Console.trace() method
func Trace(args ...interface{}) {
	js.Global().Get("console").Call("trace", args...)
}

// Javascript Console.warn() method
func Warn(args ...interface{}) {
	js.Global().Get("console").Call("warn", args...)
}
