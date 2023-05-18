//go:build js && wasm
// +build js,wasm

package console

import (
	"fmt"
	"syscall/js"
)

// MISSING:
// console.table()
//
// Most functions call js.Call() with a variadic argument list.
// Thus, the arguments get mapped to JavaScript values according to the js.ValueOf function.

type Stringer interface {
	String() string
}

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
	args = replaceArgs(args)
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
	args = replaceArgs(args)
	js.Global().Get("console").Call("error", args...)
}

// Javascript Console.exception() method
func Group(args ...interface{}) {
	args = replaceArgs(args)
	js.Global().Get("console").Call("group", args...)
}

// Javascript Console.groupCollapsed() method
func GroupCollapsed(args ...interface{}) {
	args = replaceArgs(args)
	js.Global().Get("console").Call("groupCollapsed", args...)
}

// Javascript Console.groupEnd() method
func GroupEnd() {
	js.Global().Get("console").Call("groupEnd")
}

// Javascript Console.info() method
func Info(args ...interface{}) {
	args = replaceArgs(args)
	js.Global().Get("console").Call("info", args...)
}

// Javascript Console.log() method
func Log(args ...interface{}) {
	args = replaceArgs(args)
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
	args = replaceArgs(args)
	js.Global().Get("console").Call("trace", args...)
}

// Javascript Console.warn() method
func Warn(args ...interface{}) {
	args = replaceArgs(args)
	js.Global().Get("console").Call("warn", args...)
}

func replaceArgs(args []interface{}) []interface{} {
	if len(args) > 0 {
		for i, arg := range args {
			if err, ok := arg.(error); ok {
				args[i] = err.Error()
			} else if err, ok := arg.(Stringer); ok {
				args[i] = err.String()
			}
		}
	}
	return args
}

// Javascript Console.Log with a format string
func Logf(format string, args ...interface{}) {
	Log(fmt.Sprintf(format, args...))
}

// Javascript Console.Error with a format string
func Errorf(format string, args ...interface{}) {
	Error(fmt.Sprintf(format, args...))
}

// Javascript Console.Info with a format string
func Infof(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

// Javascript Console.Warn with a format string
func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args...))
}

// Javascript Console.Debug with a format string
func Debugf(format string, args ...interface{}) {
	Debug(fmt.Sprintf(format, args...))
}

// Javascript Console.Trace with a format string
func Tracef(format string, args ...interface{}) {
	Trace(fmt.Sprintf(format, args...))
}

// Javascript Console.Assert with a format string
func Assertf(condition bool, format string, args ...interface{}) {
	Assert(condition, fmt.Sprintf(format, args...))
}
