package graphs

import (
	"github.com/Nigel2392/jsext/canvas"
	"github.com/Nigel2392/jsext/graphs/charts"
	"github.com/Nigel2392/jsext/graphs/options"
)

// CreateGraph creates a graph based on the options provided.
func CreateGraph(Canvas canvas.Canvas, opts options.GraphOptions) {
	switch opts.Type {
	case options.Bar:
		charts.Bar(Canvas, opts)
	case options.Line:
		charts.Line(Canvas, opts)
	case options.Pie:
		charts.Pie(Canvas, opts, false)
	case options.Donut:
		charts.Pie(Canvas, opts, true)
	default:
		panic("Invalid Graph Type")
	}
}
