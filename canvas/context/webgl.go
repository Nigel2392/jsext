//go:build js && wasm
// +build js,wasm

package context

import "syscall/js"

type WebGLRenderingContext js.Value

func (c WebGLRenderingContext) Value() js.Value {
	return js.Value(c)
}

func (c WebGLRenderingContext) Get(key string) js.Value {
	return c.Value().Get(key)
}

func (c WebGLRenderingContext) Set(key string, value interface{}) {
	c.Value().Set(key, value)
}

func (c WebGLRenderingContext) Call(method string, args ...interface{}) js.Value {
	return c.Value().Call(method, args...)
}

func (c WebGLRenderingContext) Canvas() js.Value {
	return c.Get("canvas")
}

func (c WebGLRenderingContext) DrawingBufferColorSpace(cSpace ...string) string {
	if len(cSpace) > 0 {
		c.Set("drawingBufferColorSpace", cSpace[0])
	}
	return c.Get("drawingBufferColorSpace").String()
}

func (c WebGLRenderingContext) DrawingBufferHeight() int {
	return int(c.Get("drawingBufferHeight").Float())
}

func (c WebGLRenderingContext) DrawingBufferWidth() int {
	return int(c.Get("drawingBufferWidth").Float())
}

func (c WebGLRenderingContext) UnpackColorSpace(cSpace ...string) string {
	if len(cSpace) > 0 {
		c.Set("unpackColorSpace", cSpace[0])
	}
	return c.Get("unpackColorSpace").String()
}

func (c WebGLRenderingContext) ActiveTexture(texture ...int) int {
	if len(texture) > 0 {
		c.Call("activeTexture", texture[0])
	}
	return int(c.Get("activeTexture").Float())
}

func (c WebGLRenderingContext) AttachShader(program, shader js.Value) {
	c.Call("attachShader", program, shader)
}

func (c WebGLRenderingContext) BindAttribLocation(program, location js.Value, name string) {
	c.Call("bindAttribLocation", program, location, name)
}

func (c WebGLRenderingContext) BindBuffer(target, buffer js.Value) {
	c.Call("bindBuffer", target, buffer)
}

func (c WebGLRenderingContext) BindFramebuffer(target, framebuffer js.Value) {
	c.Call("bindFramebuffer", target, framebuffer)
}

func (c WebGLRenderingContext) BindRenderbuffer(target, renderbuffer js.Value) {
	c.Call("bindRenderbuffer", target, renderbuffer)
}

func (c WebGLRenderingContext) BindTexture(target, texture js.Value) {
	c.Call("bindTexture", target, texture)
}

func (c WebGLRenderingContext) BlendColor(red, green, blue, alpha float32) {
	c.Call("blendColor", red, green, blue, alpha)
}

func (c WebGLRenderingContext) BlendEquation(mode int) {
	c.Call("blendEquation", mode)
}

func (c WebGLRenderingContext) BlendEquationSeparate(modeRGB, modeAlpha int) {
	c.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

func (c WebGLRenderingContext) BlendFunc(sfactor, dfactor int) {
	c.Call("blendFunc", sfactor, dfactor)
}

func (c WebGLRenderingContext) BlendFuncSeparate(srcRGB, dstRGB, srcAlpha, dstAlpha int) {
	c.Call("blendFuncSeparate", srcRGB, dstRGB, srcAlpha, dstAlpha)
}

func (c WebGLRenderingContext) BufferData(target int, src []byte, usage int) {
	var buffer js.Value = js.Global().Get("ArrayBuffer").New(len(src))
	js.CopyBytesToJS(buffer, src)
	c.Call("bufferData", target, buffer, usage)
}

func (c WebGLRenderingContext) BufferData2(target int, size int, usage int) {
	c.Call("bufferData", target, size, usage)
}

func (c WebGLRenderingContext) BufferDataWebGL2(target int, usage int, srcOffset int) {
	c.Call("bufferData", target, usage, srcOffset)
}

func (c WebGLRenderingContext) BufferData2WebGL2(target int, srcData []byte, usage int, srcOffset int, length ...int) {
	var buffer js.Value = js.Global().Get("ArrayBuffer").New(len(srcData))
	js.CopyBytesToJS(buffer, srcData)
	if len(length) > 0 {
		c.Call("bufferData", target, buffer, usage, srcOffset, length[0])
	} else {
		c.Call("bufferData", target, buffer, usage, srcOffset)
	}
}

func (c WebGLRenderingContext) BufferSubData(target int, dstByteOffset int, src ...[]byte) {
	if len(src) > 0 {
		var buffer js.Value = js.Global().Get("ArrayBuffer").New(len(src[0]))
		js.CopyBytesToJS(buffer, src[0])
		c.Call("bufferSubData", target, dstByteOffset, buffer)
		return
	}
	c.Call("bufferSubData", target, dstByteOffset)
}

func (c WebGLRenderingContext) BufferSubDataWebGL2(target int, dstByteOffset int, rest ...any) {
	for i, r := range rest {
		switch r := r.(type) {
		case []byte:
			var buffer js.Value = js.Global().Get("ArrayBuffer").New(len(r))
			js.CopyBytesToJS(buffer, r)
			rest[i] = buffer
		}
	}
	c.Call("bufferSubData", append([]any{target, dstByteOffset}, rest...)...)
}

func (c WebGLRenderingContext) CheckFramebufferStatus(target int) int {
	return int(c.Call("checkFramebufferStatus", target).Float())
}

func (c WebGLRenderingContext) Clear(mask int) {
	c.Call("clear", mask)
}

func (c WebGLRenderingContext) ClearColor(red, green, blue, alpha float32) {
	c.Call("clearColor", red, green, blue, alpha)
}

func (c WebGLRenderingContext) ClearDepth(depth float32) {
	c.Call("clearDepth", depth)
}

func (c WebGLRenderingContext) ClearStencil(s int) {
	c.Call("clearStencil", s)
}

func (c WebGLRenderingContext) ColorMask(red, green, blue, alpha bool) {
	c.Call("colorMask", red, green, blue, alpha)
}

func (c WebGLRenderingContext) CompileShader(shader js.Value) {
	c.Call("compileShader", shader)
}

// NOT FINISHED
