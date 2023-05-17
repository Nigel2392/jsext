//go:build js && wasm
// +build js,wasm

package canvas

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/canvas/context"
)

const (
	E   = 2.71828182845904523536028747135266249775724709369995957496696763
	Pi  = 3.14159265358979323846264338327950288419716939937510582097494459
	Phi = 1.61803398874989484820458683436563811772030917980576286213544862
	Tau = 6.28318530717958647692528676655900576839433879875021164194988918
)

type Canvas js.Value

func NewCanvas(width, height int) Canvas {
	canvas := js.Global().Get("document").Call("createElement", "canvas")
	canvas.Set("width", width)
	canvas.Set("height", height)
	return Canvas(canvas)
}

func FromQuerySelector(querySelector string) Canvas {
	return Canvas(js.Global().Get("document").Call("querySelector", querySelector))
}

func (c Canvas) Value() js.Value {
	return js.Value(c)
}

func (c Canvas) Render() jsext.Element {
	return jsext.Element(c)
}

func (c Canvas) Get(key string) js.Value {
	return c.Value().Get(key)
}

func (c Canvas) Set(key string, value interface{}) {
	c.Value().Set(key, value)
}

func (c Canvas) Call(method string, args ...interface{}) js.Value {
	return c.Value().Call(method, args...)
}

func (c Canvas) Style() jsext.Style {
	return jsext.Style(c.Get("style"))
}

func (c Canvas) Height(h ...int) int {
	if len(h) > 0 {
		c.Set("height", h[0])
	}
	return c.Get("height").Int()
}

func (c Canvas) Width(w ...int) int {
	if len(w) > 0 {
		c.Set("width", w[0])
	}
	return c.Get("width").Int()
}

func (c Canvas) Context2D() context.Context2D {
	return context.Context2D(c.Call("getContext", "2d"))
}

func (c Canvas) InnerHTML(s ...string) string {
	if len(s) > 0 {
		c.Set("innerHTML", s[0])
	}
	return c.Get("innerHTML").String()
}
