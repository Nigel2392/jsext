# Canvas Stickman

Here is an example that shows how you can draw a simple stickman with the canvas.

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
ctx.Arc(320, 240, 100, 0, canvas.Pi*2)
ctx.FillStyle("yellow")
ctx.Stroke()
ctx.Fill()

// Left Eye
ctx.BeginPath()
ctx.Arc(280, 190, 20, 0, canvas.Pi*2)
ctx.FillStyle("black")
ctx.Stroke()
ctx.Fill()

// Right Eye
ctx.BeginPath()
ctx.Arc(360, 190, 20, 0, canvas.Pi*2)
ctx.FillStyle("black")
ctx.Stroke()
ctx.Fill()

// Nose
ctx.MoveTo(320, 205)
ctx.LineTo(320, 250)
ctx.Stroke()

// Mouth
ctx.MoveTo(320, 260)
ctx.Arc(320, 260, 65, 0, canvas.Pi)
ctx.Fill()

// Body
ctx.BeginPath()

// Body
ctx.MoveTo(320, 340)
ctx.LineTo(320, 500)

// Reset to armheight
ctx.LineTo(320, 390)

// Left arm
ctx.LineTo(250, 480)
ctx.LineTo(320, 390)

// Right arm
ctx.LineTo(390, 480)
ctx.LineTo(320, 390)

// Reset to legheight
ctx.LineTo(320, 500)

// Left leg
ctx.LineTo(270, 700)
ctx.LineTo(320, 500)

// Right leg
ctx.LineTo(370, 700)
ctx.LineTo(320, 500)

ctx.Stroke()

ctx.MoveTo(190, 160)
ctx.LineTo(450, 160)
ctx.LineTo(320, 100)

ctx.FillStyle("black")
ctx.Fill()
```