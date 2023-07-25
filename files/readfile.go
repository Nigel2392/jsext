package files

import (
	"io"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2/reader"
)

func Read(f File) (io.ReadCloser, error) {
	var (
		jsFile = js.Value(f)
		fr     = fileReader.New()

		doneCh = make(chan js.Value)
		errCh  = make(chan error)

		onLoad js.Func
		onErr  js.Func

		done js.Value
		err  error
	)

	onLoad = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		doneCh <- fr.Get("result")
		return nil
	})
	onErr = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errCh <- js.Error{Value: args[0]}
		return nil
	})

	defer func() {
		onLoad.Release()
		onErr.Release()
		close(doneCh)
		close(errCh)
	}()

	fr.Set("onload", onLoad)
	fr.Set("onerror", onErr)

	fr.Call("readAsArrayBuffer", jsFile)

	select {
	case done = <-doneCh:
		return reader.NewArrayBufferReader(done), nil
	case err = <-errCh:
		return nil, err
	}
}
