# Simple example of a graph

JSExt also has support for graphs. The following example shows how to create a graph.

Supports the following types of graphs:
```go
options.Bar
options.Line
options.Pie
options.Donut
```

```go
var Application = app.App("#app")
var canvasDiv = js.Global().Get("document").Call("createElement", "div")
var Canvas = canvas.NewCanvas(640, 640)
// Canvas.Style().Border("1px solid black")
Canvas.Style().BackgroundColor("rgba(100,100,255, 0.2)")
Canvas.Style().Set("borderRadius", "10px")
canvasDiv.Set("style", "width: 100%; height: 100%;")
canvasDiv.Call("appendChild", Canvas.Value())
Application.RenderValue(canvasDiv)

var opts = options.GraphOptions{
	WidthPX:  900,
	HeightPX: 800,
	// Labels:          []string{"C", "V"},
	// Values:          []any{1, 20},
	Labels:                 []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	Values:                 []any{31, 28, 31, 305, 31, 30, 31, 31, 30, 31, 30, 31},
	Colors:                 []string{"green", "blue", "red", "yellow", "orange", "purple", "pink", "brown", "black", "teal"},
	GraphTitle:             "Graph Title",
	Type:                   options.Bar,
	GraphBackgroundOpacity: 0.4,
	GraphBorder:            true,
	ShowResults:            true,
}
graphs.CreateGraph(Canvas, opts)
```