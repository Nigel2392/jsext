//go:build js && wasm
// +build js,wasm

package path

import (
	"syscall/js"
)

type Path2D js.Value

func NewPath2D() Path2D {
	return Path2D(js.Global().Get("Path2D").New())
}

func (p Path2D) Value() js.Value {
	return js.Value(p)
}

func (p Path2D) Get(key string) js.Value {
	return p.Value().Get(key)
}

func (p Path2D) Set(key string, value interface{}) {
	p.Value().Set(key, value)
}

func (p Path2D) Call(method string, args ...interface{}) js.Value {
	return p.Value().Call(method, args...)
}

func (p Path2D) AddPath(path Path2D, transform ...interface{}) {
	p.Call("addPath", path, transform)
}

func (p Path2D) ClosePath() {
	p.Call("closePath")
}

func (p Path2D) MoveTo(x, y float64) {
	p.Call("moveTo", x, y)
}

func (p Path2D) LineTo(x, y float64) {
	p.Call("lineTo", x, y)
}

func (p Path2D) QuadraticCurveTo(cpx, cpy, x, y float64) {
	p.Call("quadraticCurveTo", cpx, cpy, x, y)
}

func (p Path2D) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	p.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x, y)
}

func (p Path2D) Arc(x, y, radius, startAngle, endAngle float64, anticlockwise ...bool) {
	if len(anticlockwise) > 0 {
		p.Call("arc", x, y, radius, startAngle, endAngle, anticlockwise[0])
	} else {
		p.Call("arc", x, y, radius, startAngle, endAngle)
	}
}

func (p Path2D) ArcTo(x1, y1, x2, y2, radius float64) {
	p.Call("arcTo", x1, y1, x2, y2, radius)
}

func (p Path2D) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, anticlockwise ...bool) {
	if len(anticlockwise) > 0 {
		p.Call("ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle, anticlockwise[0])
	} else {
		p.Call("ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle)
	}
}

func (p Path2D) Rect(x, y, width, height float64) {
	p.Call("rect", x, y, width, height)
}

func (p Path2D) RoundRect(x, y, width, height, radius float64) {
	p.Call("roundRect", x, y, width, height, radius)
}
