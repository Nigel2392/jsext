//go:build js && wasm
// +build js,wasm

package context

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/canvas/path"
)

type Context2D js.Value

func NewContext2D(canvas js.Value) Context2D {
	return Context2D(canvas.Call("getContext", "2d"))
}

func (c Context2D) Value() js.Value {
	return js.Value(c)
}

func (c Context2D) Get(key string) js.Value {
	return c.Value().Get(key)
}

func (c Context2D) Set(key string, value interface{}) {
	c.Value().Set(key, value)
}

func (c Context2D) Call(method string, args ...interface{}) js.Value {
	return c.Value().Call(method, args...)
}

func (c Context2D) Canvas() js.Value {
	return c.Get("canvas")
}

func (c Context2D) Direction(d ...string) string {
	if len(d) > 0 {
		c.Set("direction", d[0])
	}
	return c.Get("direction").String()
}

func (c Context2D) FillStyle(s ...string) string {
	if len(s) > 0 {
		c.Set("fillStyle", s[0])
	}
	return c.Get("fillStyle").String()
}

func (c Context2D) Filter(f ...string) string {
	if len(f) > 0 {
		c.Set("filter", f[0])
	}
	return c.Get("filter").String()
}

func (c Context2D) Font(f ...string) string {
	if len(f) > 0 {
		c.Set("font", f[0])
	}
	return c.Get("font").String()
}

func (c Context2D) FontKerning(f ...string) string {
	if len(f) > 0 {
		c.Set("fontKerning", f[0])
	}
	return c.Get("fontKerning").String()
}

func (c Context2D) FontStretch(f ...string) string {
	if len(f) > 0 {
		c.Set("fontStretch", f[0])
	}
	return c.Get("fontStretch").String()
}

func (c Context2D) FontVariantCaps(f ...string) string {
	if len(f) > 0 {
		c.Set("fontVariantCaps", f[0])
	}
	return c.Get("fontVariantCaps").String()
}

func (c Context2D) GlobalAlpha(a ...float64) float64 {
	if len(a) > 0 {
		c.Set("globalAlpha", a[0])
	}
	return c.Get("globalAlpha").Float()
}

func (c Context2D) GlobalCompositeOperation(o ...string) string {
	if len(o) > 0 {
		c.Set("globalCompositeOperation", o[0])
	}
	return c.Get("globalCompositeOperation").String()
}

func (c Context2D) ImageSmoothingEnabled(e ...bool) bool {
	if len(e) > 0 {
		c.Set("imageSmoothingEnabled", e[0])
	}
	return c.Get("imageSmoothingEnabled").Bool()
}

func (c Context2D) ImageSmoothingQuality(q ...string) string {
	if len(q) > 0 {
		c.Set("imageSmoothingQuality", q[0])
	}
	return c.Get("imageSmoothingQuality").String()
}

func (c Context2D) LetterSpacing(l ...string) string {
	if len(l) > 0 {
		c.Set("letterSpacing", l[0])
	}
	return c.Get("letterSpacing").String()
}

func (c Context2D) LineCap(col ...string) string {
	if len(col) > 0 {
		c.Set("lineCap", col[0])
	}
	return c.Get("lineCap").String()
}

func (c Context2D) LineDashOffset(o ...float64) float64 {
	if len(o) > 0 {
		c.Set("lineDashOffset", o[0])
	}
	return c.Get("lineDashOffset").Float()
}

func (c Context2D) LineJoin(j ...string) string {
	if len(j) > 0 {
		c.Set("lineJoin", j[0])
	}
	return c.Get("lineJoin").String()
}

func (c Context2D) LineWidth(w ...float64) float64 {
	if len(w) > 0 {
		c.Set("lineWidth", w[0])
	}
	return c.Get("lineWidth").Float()
}

func (c Context2D) MiterLimit(l ...float64) float64 {
	if len(l) > 0 {
		c.Set("miterLimit", l[0])
	}
	return c.Get("miterLimit").Float()
}

func (c Context2D) ShadowBlur(b ...float64) float64 {
	if len(b) > 0 {
		c.Set("shadowBlur", b[0])
	}
	return c.Get("shadowBlur").Float()
}

func (c Context2D) ShadowColor(col ...string) string {
	if len(col) > 0 {
		c.Set("shadowColor", col[0])
	}
	return c.Get("shadowColor").String()
}

func (c Context2D) ShadowOffsetX(x ...float64) float64 {
	if len(x) > 0 {
		c.Set("shadowOffsetX", x[0])
	}
	return c.Get("shadowOffsetX").Float()
}

func (c Context2D) ShadowOffsetY(y ...float64) float64 {
	if len(y) > 0 {
		c.Set("shadowOffsetY", y[0])
	}
	return c.Get("shadowOffsetY").Float()
}

func (c Context2D) StrokeStyle(s ...string) string {
	if len(s) > 0 {
		c.Set("strokeStyle", s[0])
	}
	return c.Get("strokeStyle").String()
}

func (c Context2D) TextAlign(t ...string) string {
	if len(t) > 0 {
		c.Set("textAlign", t[0])
	}
	return c.Get("textAlign").String()
}

func (c Context2D) TextBaseline(t ...string) string {
	if len(t) > 0 {
		c.Set("textBaseline", t[0])
	}
	return c.Get("textBaseline").String()
}

func (c Context2D) TextRendering(t ...string) string {
	if len(t) > 0 {
		c.Set("textRendering", t[0])
	}
	return c.Get("textRendering").String()
}

func (c Context2D) WordSpacing(w ...string) string {
	if len(w) > 0 {
		c.Set("wordSpacing", w[0])
	}
	return c.Get("wordSpacing").String()
}

func (c Context2D) Arc(x, y, radius, startAngle, endAngle float64, anticlockwise ...bool) {
	if len(anticlockwise) > 0 {
		c.Call("arc", x, y, radius, startAngle, endAngle, anticlockwise[0])
	} else {
		c.Call("arc", x, y, radius, startAngle, endAngle)
	}
}

func (c Context2D) ArcTo(x1, y1, x2, y2, radius float64) {
	c.Call("arcTo", x1, y1, x2, y2, radius)
}

func (c Context2D) BeginPath() {
	c.Call("beginPath")
}

func (c Context2D) BeginPath2(startX, startY float64) {
	c.Call("beginPath", startX, startY)
}

func (c Context2D) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	c.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x, y)
}

func (c Context2D) ClearRect(x, y, width, height float64) {
	c.Call("clearRect", x, y, width, height)
}

func (c Context2D) Clip(fillRule ...string) {
	if len(fillRule) > 0 {
		c.Call("clip", fillRule[0])
	} else {
		c.Call("clip")
	}
}

func (c Context2D) ClosePath() {
	c.Call("closePath")
}

func (c Context2D) CreateConicGradient(startAngle, x, y float64) js.Value {
	return c.Call("createConicGradient", startAngle, x, y)
}

func (c Context2D) CreateImageData(width, height float64, settings map[string]any) js.Value {
	if len(settings) > 0 {
		var obj js.Value = js.Global().Get("Object").New()
		for k, v := range settings {
			obj.Set(k, v)
		}
		return c.Call("createImageData", width, height, obj)
	}
	return c.Call("createImageData", width, height)
}

func (c Context2D) CreateImageDataFrom(data js.Value) js.Value {
	return c.Call("createImageData", data)
}

func (c Context2D) CreateLinearGradient(x0, y0, x1, y1 float64) js.Value {
	return c.Call("createLinearGradient", x0, y0, x1, y1)
}

func (c Context2D) CreatePattern(image js.Value, repetition string) js.Value {
	return c.Call("createPattern", image, repetition)
}

func (c Context2D) CreateRadialGradient(x0, y0, r0, x1, y1, r1 float64) js.Value {
	return c.Call("createRadialGradient", x0, y0, r0, x1, y1, r1)
}

func (c Context2D) DrawFocusIfNeeded(element js.Value, path ...path.Path2D) {
	if len(path) > 0 {
		c.Call("drawFocusIfNeeded", element, path[0])
	} else {
		c.Call("drawFocusIfNeeded", element)
	}
}

func (c Context2D) DrawImage(image js.Value, dx, dy float64, settings map[string]any) {
	if len(settings) > 0 {
		var obj js.Value = js.Global().Get("Object").New()
		for k, v := range settings {
			obj.Set(k, v)
		}
		c.Call("drawImage", image, dx, dy, obj)
	} else {
		c.Call("drawImage", image, dx, dy)
	}
}

func (c Context2D) DrawImage2(image js.Value, dx, dy, dw, dh float64, settings map[string]any) {
	if len(settings) > 0 {
		var obj js.Value = js.Global().Get("Object").New()
		for k, v := range settings {
			obj.Set(k, v)
		}
		c.Call("drawImage", image, dx, dy, dw, dh, obj)
	} else {
		c.Call("drawImage", image, dx, dy, dw, dh)
	}
}

func (c Context2D) DrawImage3(image js.Value, sx, sy, sw, sh, dx, dy, dw, dh float64, settings map[string]any) {
	if len(settings) > 0 {
		var obj js.Value = js.Global().Get("Object").New()
		for k, v := range settings {
			obj.Set(k, v)
		}
		c.Call("drawImage", image, sx, sy, sw, sh, dx, dy, dw, dh, obj)
	} else {
		c.Call("drawImage", image, sx, sy, sw, sh, dx, dy, dw, dh)
	}
}

func (c Context2D) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, anticlockwise ...bool) {
	if len(anticlockwise) > 0 {
		c.Call("ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle, anticlockwise[0])
	} else {
		c.Call("ellipse", x, y, radiusX, radiusY, rotation, startAngle, endAngle)
	}
}

func (c Context2D) Fill(fillRule ...string) {
	if len(fillRule) > 0 {
		c.Call("fill", fillRule[0])
	} else {
		c.Call("fill")
	}
}

func (c Context2D) FillRect(x, y, width, height float64) {
	c.Call("fillRect", x, y, width, height)
}

func (c Context2D) FillText(text string, x, y float64, maxWidth ...float64) {
	if len(maxWidth) > 0 {
		c.Call("fillText", text, x, y, maxWidth[0])
	} else {
		c.Call("fillText", text, x, y)
	}
}

func (c Context2D) GetContextAttributes() js.Value {
	return c.Call("getContextAttributes")
}

func (c Context2D) GetImageData(sx, sy, sw, sh float64, settings ...map[string]any) js.Value {
	if len(settings) > 0 {
		var obj js.Value = js.Global().Get("Object").New()
		for k, v := range settings[0] {
			obj.Set(k, v)
		}
		return c.Call("getImageData", sx, sy, sw, sh, obj)
	}
	return c.Call("getImageData", sx, sy, sw, sh)
}

func (c Context2D) GetLineDash() js.Value {
	return c.Call("getLineDash")
}

func (c Context2D) GetTransform() js.Value {
	return c.Call("getTransform")
}

func (c Context2D) IsContextLost() bool {
	return c.Call("isContextLost").Bool()
}

func (c Context2D) IsPointInPath(x, y float64, fillRule ...string) bool {
	if len(fillRule) > 0 {
		return c.Call("isPointInPath", x, y, fillRule[0]).Bool()
	}
	return c.Call("isPointInPath", x, y).Bool()
}

func (c Context2D) IsPointInPath2(path path.Path2D, x, y float64, fillRule ...string) bool {
	if len(fillRule) > 0 {
		return c.Call("isPointInPath", path, x, y, fillRule[0]).Bool()
	}
	return c.Call("isPointInPath", path, x, y).Bool()
}

func (c Context2D) IsPointInStroke(x, y float64, path ...string) bool {
	if len(path) > 0 {
		return c.Call("isPointInStroke", path[0], x, y).Bool()
	}
	return c.Call("isPointInStroke", x, y).Bool()
}

func (c Context2D) LineTo(x, y float64) {
	c.Call("lineTo", x, y)
}

func (c Context2D) MeasureText(text string) js.Value {
	return c.Call("measureText", text)
}

func (c Context2D) MoveTo(x, y float64) {
	c.Call("moveTo", x, y)
}

func (c Context2D) PutImageData(imageData js.Value, dx, dy float64) {
	c.Call("putImageData", imageData, dx, dy)
}

func (c Context2D) PutImageData2(imageData js.Value, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight float64) {
	c.Call("putImageData", imageData, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight)
}

func (c Context2D) QuadraticCurveTo(cpx, cpy, x, y float64) {
	c.Call("quadraticCurveTo", cpx, cpy, x, y)
}

func (c Context2D) Rect(x, y, width, height float64) {
	c.Call("rect", x, y, width, height)
}

func (c Context2D) Reset() {
	c.Call("reset")
}

func (c Context2D) ResetTransform() {
	c.Call("resetTransform")
}

func (c Context2D) Restore() {
	c.Call("restore")
}

func (c Context2D) Rotate(angle float64) {
	c.Call("rotate", angle)
}

func (c Context2D) RoundRect(x, y, width, height, radius float64) {
	c.Call("roundRect", x, y, width, height, radius)
}

func (c Context2D) Save() {
	c.Call("save")
}

func (c Context2D) Scale(x, y float64) {
	c.Call("scale", x, y)
}

func (c Context2D) ScrollPathIntoView(path ...string) {
	if len(path) > 0 {
		c.Call("scrollPathIntoView", path[0])
	} else {
		c.Call("scrollPathIntoView")
	}
}

func (c Context2D) SetLineDash(segments []float64) {
	var arr js.Value = js.Global().Get("Array").New(len(segments))
	for i, v := range segments {
		arr.SetIndex(i, v)
	}
	c.Call("setLineDash", segments)
}

func (ctx Context2D) SetTransform(a, b, c, d, e, f float64) {
	ctx.Call("setTransform", a, b, c, d, e, f)
}

func (c Context2D) SetTransformMatrix(matrix js.Value) {
	c.Call("setTransform", matrix)
}

func (c Context2D) Stroke(path ...string) {
	if len(path) > 0 {
		c.Call("stroke", path[0])
	} else {
		c.Call("stroke")
	}
}

func (c Context2D) StrokeRect(x, y, width, height float64) {
	c.Call("strokeRect", x, y, width, height)
}

func (c Context2D) StrokeText(text string, x, y float64, maxWidth ...float64) {
	if len(maxWidth) > 0 {
		c.Call("strokeText", text, x, y, maxWidth[0])
	} else {
		c.Call("strokeText", text, x, y)
	}
}

func (ctx Context2D) Transform(a, b, c, d, e, f float64) {
	ctx.Call("transform", a, b, c, d, e, f)
}

func (ctx Context2D) Translate(x, y float64) {
	ctx.Call("translate", x, y)
}
