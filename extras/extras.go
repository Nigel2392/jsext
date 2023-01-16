//go:build js && wasm
// +build js,wasm

package extras

//	type MousePosition struct {
//		X chan int
//		Y chan int
//	}
//
//	func (m *MousePosition) String() string {
//		return fmt.Sprintf("X: %d, Y: %d", m.X, m.Y)
//	}
//
//	// GetMousePosition returns the current mouse position relative to the document.
//	func GetMousePosition() MousePosition {
//		var mousePosition MousePosition = MousePosition{
//			X: make(chan int),
//			Y: make(chan int),
//		}
//		var mouseMove = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
//			mousePosition.X <- args[0].Get("clientX").Int()
//			mousePosition.Y <- args[0].Get("clientY").Int()
//			return nil
//		})
//		js.Global().Get("document").Call("addEventListener", "mousemove", mouseMove)
//		return mousePosition
//	}
//
