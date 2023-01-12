//go:build js && wasm
// +build js,wasm

package charts

import (
	"math"

	"github.com/Nigel2392/jsext/canvas"
	"github.com/Nigel2392/jsext/framework/graphs/options"
	"github.com/Nigel2392/jsext/framework/helpers/convert"
)

func Pie(Canvas canvas.Canvas, opts options.GraphOptions, donut bool) {
	var width, height = float64(opts.WidthPX), float64(opts.HeightPX)
	var ctx = Canvas.Context2D()
	var marginTop = height * 0.1

	Canvas.Width(int(width))
	Canvas.Height(int(height))

	var radius = Min(width, height)/2 - marginTop
	var centerX = width / 2
	var centerY = height/2 + marginTop/2

	marginTop = height * 0.1
	var total = 0.0
	for _, v := range opts.Values {
		total += convert.ToFloat(v)
	}

	var drawPieChart = func() {

		writeTitle(ctx, int(width), int(marginTop), opts)

		var startAngle = 0.0
		var lastColor = ""
		for i, v := range opts.Values {
			var value = convert.ToFloat(v)
			var sliceAngle = 2 * math.Pi * value / total
			// Get the old alpha value.

			var opacity float64 = opts.GraphBackgroundOpacity
			var oldAlpha = ctx.GlobalAlpha()
			var color = ""
			if len(opts.Colors) > 0 {
				color = opts.Colors[(i)%len(opts.Colors)]
				if color == lastColor {
					color = opts.Colors[(i+1)%len(opts.Colors)]
				}
				lastColor = color
			}
			if opacity <= 0 {
				opacity = 1
			}
			// Set the new alpha value.
			ctx.GlobalAlpha(opacity)

			ctx.BeginPath()
			ctx.MoveTo(centerX, centerY)
			ctx.Arc(centerX, centerY, radius, startAngle, startAngle+sliceAngle)
			ctx.ClosePath()

			ctx.FillStyle(color)
			ctx.Fill()

			if opts.GraphBorder {
				// Draw the border.
				ctx.LineWidth(2)
				ctx.StrokeStyle(color)
				ctx.Stroke()
			}

			if opts.ShowResults {
				// Draw the value of the slice.
				ctx.FillStyle("white")
				ctx.TextAlign("center")
				ctx.TextBaseline("middle")
				ctx.Font(convert.FormatNumber(math.Cos(startAngle+sliceAngle/2)/2) + " Arial")
				// Write the text at outer edge of the slice, but not too close to the edge.
				var textRadius = radius * 0.8
				var textX = centerX + textRadius*math.Cos(startAngle+sliceAngle/2)
				var textY = centerY + textRadius*math.Sin(startAngle+sliceAngle/2)
				ctx.FillText(convert.FormatNumber(opts.Values[i]), textX, textY)
			}

			ctx.GlobalAlpha(oldAlpha)

			startAngle += sliceAngle
		}

		if donut {
			// Draw the center of the pie chart.
			ctx.BeginPath()
			ctx.GlobalCompositeOperation("destination-out")
			ctx.Arc(centerX, centerY, radius*0.35, 0, 2*math.Pi)
			ctx.ClosePath()
			ctx.Fill()
			ctx.GlobalCompositeOperation("source-over")
		}
	}

	drawPieChart()

}
