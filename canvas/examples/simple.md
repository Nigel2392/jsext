# Simple canvas example

Simple example to show how to manipulate the canvas.

```go
var Canvas = canvas.NewCanvas(500, 500)
var style = Canvas.Style()
style.Border("1px solid black")
var ctx = Canvas.Context2D()

// Set the canvas size.
Canvas.Width(640)
Canvas.Height(640)

// Draw a filled rectangle.
ctx.FillStyle("#00ff00")
ctx.FillRect(100, 100, 450, 300)

// Draw an outline.
ctx.StrokeStyle("#ff0000")
ctx.StrokeRect(90, 90, 470, 320)

// Clear the rectangle. (Make it transparent)
ctx.ClearRect(125, 125, 400, 250)

// Draw another square
ctx.FillStyle("#0000ff")
ctx.FillRect(150, 150, 350, 200)

js.Global().Get("document").Get("body").Call("appendChild", Canvas.Value())
```
