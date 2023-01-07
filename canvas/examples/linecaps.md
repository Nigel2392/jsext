# Canvas Linecaps

Here is an example of what canvas linecaps look like.

```go
// SETUP CANVAS
///////////////////////////////////////
var Canvas = canvas.NewCanvas(640, 640)
// Append element to document here:

// 
// Get the canvas style to set it.
var style = Canvas.Style()
style.Border("1px solid black")
var ctx = Canvas.Context2D()
Canvas.Width(640)
Canvas.Height(640)
///////////////////////////////////////

for i := 0; i < 20; i++ {
	time.Sleep(100 * time.Millisecond)
	ctx.LineWidth(1 + float64(i*1))
	ctx.BeginPath()
	ctx.MoveTo(50, float64(10+(i*30)))
	ctx.LineTo(500, float64(10+(i*30)))
	ctx.Stroke()
}

ctx.ClearRect(0, 0, 640, 640)

// Line caps
for i := 0; i < 20; i++ {
	time.Sleep(100 * time.Millisecond)
	ctx.LineWidth(1 + float64(i*1))

	switch i % 3 {
	case 0:
		ctx.LineCap("butt") // butt, round, square
	case 1:
		ctx.LineCap("round") // butt, round, square
	case 2:
		ctx.LineCap("square") // butt, round, square
	}

	ctx.BeginPath()
	ctx.MoveTo(float64(10+(i*30)), 50)
	ctx.LineTo(float64(10+(i*30)), 500)
	ctx.Stroke()
}
```