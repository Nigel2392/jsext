# Canvas Smiley

Here is an example that shows how you can draw a simple smiley with the canvas.

```go
var Canvas = canvas.NewCanvas(640, 640)
var style = Canvas.Style()
style.Border("1px solid black")
var ctx = Canvas.Context2D()
// Set the canvas size.
Canvas.Width(640)
Canvas.Height(640)

// Head
ctx.BeginPath()
ctx.Arc(320, 320, 100, 0, 2*canvas.Pi)
ctx.FillStyle("#f0f007")
ctx.Fill()
ctx.ClosePath()

// Left eye
ctx.BeginPath()
ctx.Arc(280, 300, 30, 0, 2*canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("black")
ctx.Fill()

ctx.BeginPath()
ctx.Arc(280, 300, 27, 0, 2*canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("aqua")
ctx.Fill()

ctx.BeginPath()
ctx.Arc(285, 305, 20, 0, 2*canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("white")
ctx.Fill()

// Right Eye
ctx.BeginPath()
ctx.Arc(360, 300, 30, 0, 2*canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("black")
ctx.Fill()

ctx.BeginPath()
ctx.Arc(360, 300, 27, 0, 2*canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("aqua")
ctx.Fill()

ctx.BeginPath()
ctx.Arc(355, 305, 20, 0, 2*canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("white")
ctx.Fill()

// Mouth
ctx.BeginPath()
ctx.Arc(320, 360, 40, 0, canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("black")
ctx.Fill()

ctx.BeginPath()
ctx.Arc(320, 360, 37, 0, canvas.Pi)
ctx.ClosePath()
ctx.FillStyle("white")
ctx.Fill()
ctx.StrokeStyle("black")
ctx.Stroke()
```