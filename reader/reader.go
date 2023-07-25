package reader

import (
	"errors"
	"io"
	"syscall/js"
)

// Taken from go/src/net/http/roundtrip_js.go

var errClosed = errors.New("jsext/reader: reader is closed")
var uint8Array = js.Global().Get("Uint8Array")

// streamReader implements an io.ReadCloser wrapper for ReadableStream.
// See https://fetch.spec.whatwg.org/#readablestream for more information.
type streamReader struct {
	pending []byte
	stream  js.Value
	err     error // sticky read error
}

func NewStreamReader(stream js.Value) io.ReadCloser {
	return &streamReader{stream: stream}
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if len(r.pending) == 0 {
		var (
			bCh   = make(chan []byte, 1)
			errCh = make(chan error, 1)
		)
		success := js.FuncOf(func(this js.Value, args []js.Value) any {
			result := args[0]
			if result.Get("done").Bool() {
				errCh <- io.EOF
				return nil
			}
			value := make([]byte, result.Get("value").Get("byteLength").Int())
			js.CopyBytesToGo(value, result.Get("value"))
			bCh <- value
			return nil
		})
		defer success.Release()
		failure := js.FuncOf(func(this js.Value, args []js.Value) any {
			// Assumes it's a TypeError. See
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/TypeError
			// for more information on this type. See
			// https://streams.spec.whatwg.org/#byob-reader-read for the spec on
			// the read method.
			errCh <- errors.New(args[0].Get("message").String())
			return nil
		})
		defer failure.Release()
		r.stream.Call("read").Call("then", success, failure)
		select {
		case b := <-bCh:
			r.pending = b
		case err := <-errCh:
			r.err = err
			return 0, err
		}
	}
	n = copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}

func (r *streamReader) Close() error {
	// This ignores any error returned from cancel method. So far, I did not encounter any concrete
	// situation where reporting the error is meaningful. Most users ignore error from resp.Body.Close().
	// If there's a need to report error here, it can be implemented and tested when that need comes up.
	r.stream.Call("cancel")
	if r.err == nil {
		r.err = errClosed
	}
	return nil
}

// arrayReader implements an io.ReadCloser wrapper for ArrayBuffer.
// https://developer.mozilla.org/en-US/docs/Web/API/Body/arrayBuffer.
type arrayReader struct {
	buf       js.Value
	pending   []byte
	read      bool
	err       error // sticky read error
	readFunc  func(js.Value) ([]byte, error)
	closeFunc func(js.Value)
}

func NewArrayBufferReader(arrayBuffer js.Value) io.ReadCloser {
	return &arrayReader{buf: arrayBuffer, readFunc: readArrayBuffer}
}

func NewArrayPromiseReader(arrayPromise js.Value) io.ReadCloser {
	return &arrayReader{buf: arrayPromise, readFunc: readArrayPromise}
}

func (r *arrayReader) Read(p []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	if !r.read {
		r.read = true
		var b, err = r.readFunc(r.buf)
		if err != nil {
			return 0, err
		}
		r.pending = b
	}
	if len(r.pending) == 0 {
		return 0, io.EOF
	}
	n = copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}

func (r *arrayReader) Close() error {
	if r.closeFunc != nil {
		r.closeFunc(r.buf)
	}
	if r.err == nil {
		r.err = errClosed
	}
	return nil
}

func readArrayBuffer(arrayBuffer js.Value) ([]byte, error) {
	// Wrap the input ArrayBuffer with a Uint8Array
	if arrayBuffer.IsUndefined() || arrayBuffer.IsNull() {
		return nil, io.EOF
	}
	uint8arrayWrapper := uint8Array.New(arrayBuffer)
	value := make([]byte, uint8arrayWrapper.Get("byteLength").Int())
	js.CopyBytesToGo(value, uint8arrayWrapper)
	return value, nil
}

func readArrayPromise(arrayPromise js.Value) ([]byte, error) {
	var (
		bCh   = make(chan []byte, 1)
		errCh = make(chan error, 1)
	)
	success := js.FuncOf(func(this js.Value, args []js.Value) any {
		// Wrap the input ArrayBuffer with a Uint8Array
		uint8arrayWrapper := uint8Array.New(args[0])
		value := make([]byte, uint8arrayWrapper.Get("byteLength").Int())
		js.CopyBytesToGo(value, uint8arrayWrapper)
		bCh <- value
		return nil
	})
	defer success.Release()
	failure := js.FuncOf(func(this js.Value, args []js.Value) any {
		// Assumes it's a TypeError. See
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/TypeError
		// for more information on this type.
		// See https://fetch.spec.whatwg.org/#concept-body-consume-body for reasons this might error.
		errCh <- errors.New(args[0].Get("message").String())
		return nil
	})
	defer failure.Release()
	arrayPromise.Call("then", success, failure)
	select {
	case b := <-bCh:
		return b, nil
	case err := <-errCh:
		return nil, err
	}
}
