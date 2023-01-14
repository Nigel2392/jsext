//go:build js && wasm
// +build js,wasm

package charts

import (
	"strconv"
	"syscall/js"

	"github.com/Nigel2392/jsext/canvas"
	"github.com/Nigel2392/jsext/canvas/context"
	"github.com/Nigel2392/jsext/framework/graphs/labels"
	"github.com/Nigel2392/jsext/framework/graphs/options"
	"github.com/Nigel2392/jsext/framework/helpers"
	"github.com/Nigel2392/jsext/framework/helpers/convert"
)

// Line chart is bugged.
func Line(Canvas canvas.Canvas, opts options.GraphOptions) {
	// SETUP CANVAS
	///////////////////////////////////////
	var width = opts.WidthPX
	var height = opts.HeightPX
	Canvas.Width(width)
	Canvas.Height(height)
	var ctx = Canvas.Context2D()
	// Calculate the max height of the highest value.
	var maxHeight float64
	for _, value := range opts.Values {
		if convert.ToFloat(value) > maxHeight {
			maxHeight = convert.ToFloat(value)
		}
	}
	// Set up graph labelOptions
	var labelOptions = labels.LabelOpts{}
	labelOptions.SetLabels(opts.Labels, opts.Values)

	var topMargin float64 = float64(height) / 10
	var leftMargin float64 = float64(width) / 10
	var topOfChart = topMargin
	var bottomOfChart = float64(height) * 0.85
	var heightOfChart = bottomOfChart - topOfChart
	// Draw the labels across the bottom of the graph.
	var leftOfChart = leftMargin
	var rightOfChart = float64(width)
	var widthOfChart = rightOfChart - leftOfChart
	var labelWidth = widthOfChart / float64(len(labelOptions.Labels))
	var labelHeight = float64(height) * 0.85

	var drawCanvas = func() {
		writeTitle(ctx, width, int(topMargin/2), opts)

		// Delineate the graph.
		ctx.BeginPath()
		ctx.MoveTo(leftMargin, float64(height)*0.85)
		ctx.LineTo(float64(width), float64(height)*0.85)
		ctx.StrokeStyle("#000000")
		ctx.LineWidth(1)
		ctx.Stroke()

		ctx.BeginPath()
		ctx.MoveTo(leftMargin, topMargin)
		ctx.LineTo(leftMargin, float64(height)*0.85)
		ctx.StrokeStyle("#000000")
		ctx.LineWidth(1)
		ctx.Stroke()

		// Draw a dot at the zero point.
		ctx.BeginPath()
		ctx.Arc(leftMargin, float64(height)*0.85, 3, 0, 2*3.14159)
		ctx.FillStyle("#000000")
		ctx.Fill()

		//Draw values across the left side of the graph.
		if maxHeight == 0 || maxHeight/10 < 1 {
			maxHeight = 10
		}
		var valueIncrement = maxHeight / 10
		var valueHeight = heightOfChart / 10
		var value float64 = 0
		for i := 0; i < 11; i++ {
			valStr := strconv.Itoa(int(value))
			// Calculate the text size.
			var textSize = topMargin - (topMargin / 5)
			if textSize > labelWidth {
				textSize = labelWidth
			}
			textSize = textSize - (textSize / 2)

			ctx.BeginPath()
			ctx.Arc(leftMargin, bottomOfChart-(float64(i)*valueHeight), 3, 0, 2*3.14159)
			ctx.Fill()
			ctx.FillStyle("#000000")

			// Draw a gridline.
			ctx.BeginPath()
			ctx.MoveTo(leftMargin, bottomOfChart-(float64(i)*valueHeight))
			ctx.LineTo(float64(width), bottomOfChart-(float64(i)*valueHeight))
			ctx.StrokeStyle("rgba(0,0,0,0.1)")
			ctx.LineWidth(1)
			ctx.Stroke()

			ctx.Font(strconv.Itoa(int(textSize)) + "px Arial")
			ctx.TextAlign("right")
			ctx.TextBaseline("middle")
			ctx.FillText(valStr, leftMargin-5, bottomOfChart-(float64(i)*valueHeight))
			value += valueIncrement
		}

		for i, label := range labelOptions.Labels {
			// Calculate the text size.
			var textSize int = int(calcTextSize(label, topMargin, labelWidth))
			ctx.FillStyle("#000000")
			ctx.Font(strconv.Itoa(int(textSize)) + "px Arial")
			ctx.TextAlign("center")
			ctx.TextBaseline("top")
			ctx.FillText(label, leftOfChart+(float64(i)*labelWidth)+(labelWidth/2), labelHeight+5)
			// Draw a dot at the zero point.
			ctx.BeginPath()
			ctx.Arc(leftOfChart+(float64(i)*labelWidth)+(labelWidth/2), float64(height)*0.85, 3, 0, 2*3.14159)
			ctx.FillStyle("#000000")
			ctx.Fill()
		}
		// Get the old alpha value.
		var oldAlpha = ctx.GlobalAlpha()
		var opacity = opts.GraphBackgroundOpacity
		if opacity <= 0 {
			opacity = 1
		}
		// Set the new alpha value.
		ctx.GlobalAlpha(opacity)

		// Draw the values on the graph.
		var lastX float64
		var lastY float64
		var ct int
		var color string
		var lineColor string
		for i, value := range opts.Values {
			var x = leftOfChart + (float64(i) * labelWidth) + (labelWidth / 2)
			var y = bottomOfChart - ((convert.ToFloat(value) / maxHeight) * heightOfChart)
			if i == 0 {
				lastX = x
				lastY = y
			}

			lineColor, ct = helpers.GetColor(opts.BackgroundColors, ct, "#000000")
			drawLine(ctx, lastX, lastY, x, y, lineColor)
			lastX = x
			lastY = y

		}
		for i, value := range opts.Values {
			var x = leftOfChart + (float64(i) * labelWidth) + (labelWidth / 2)
			var y = bottomOfChart - ((convert.ToFloat(value) / maxHeight) * heightOfChart)
			if i == 0 {
				lastX = x
				lastY = y
			}

			color, ct = helpers.GetColor(opts.Colors, ct, "#9200ff")
			ctx.BeginPath()
			ctx.Arc(x, y, widthOfChart/75, 0, 2*3.14159)
			ctx.FillStyle(color)
			ctx.Fill()

			if opts.GraphBorder {
				ctx.BeginPath()
				ctx.Arc(x, y, widthOfChart/75, 0, 2*3.14159)
				ctx.StrokeStyle("#000000")
				ctx.Stroke()
			}
		}
		// Reset the alpha value.
		ctx.GlobalAlpha(oldAlpha)
	}

	drawCanvas()

	// Add eventlisteners to the canvas dots.
	Canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		var event = args[0]
		var x = event.Get("offsetX").Float()
		var y = event.Get("offsetY").Float()

		// Get all positions of the dots.
		var positions = make([]float64, 0)
		for i, value := range opts.Values {
			var x = leftOfChart + (float64(i) * labelWidth) + (labelWidth / 2)
			var y = bottomOfChart - ((convert.ToFloat(value) / maxHeight) * heightOfChart)
			positions = append(positions, x, y)
		}

		// Check if the mouse is over a dot.
		for i := 0; i < len(positions); i += 2 {
			var x1 = positions[i]
			var y1 = positions[i+1]
			var x2 = x1 + (widthOfChart / 100)
			var y2 = y1 + (widthOfChart / 100)
			if x > x1 && x < x2 && y > y1 && y < y2 {
				ctx.ClearRect(0, 0, float64(opts.WidthPX), float64(opts.HeightPX))
				drawCanvas()
				drawTooltip(ctx, width, height, x1-50, y1-20, labelOptions.UnderlyingLabel(i/2)+": "+convert.FormatNumber(opts.Values[i/2]))
			} else {
				for i := 0; i < len(opts.Values); i++ {
					var x1 = leftOfChart + (float64(i) * labelWidth) + (labelWidth / 2)
					var y1 = bottomOfChart - ((convert.ToFloat(opts.Values[i]) / maxHeight) * heightOfChart)
					var x2 = x1 + (widthOfChart / 100)
					var y2 = y1 + (widthOfChart / 100)
					if x > x1-(widthOfChart/75) && x < x2+(widthOfChart/75) && y > y1-(widthOfChart/75) && y < y2+(widthOfChart/75) {
						break
					}
					if i == len(opts.Values)-1 {
						ctx.ClearRect(0, 0, float64(width), float64(height))
						drawCanvas()
					}
				}

			}
		}
		return nil
	}))
	Canvas.Call("addEventListener", "mouseout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.ClearRect(0, 0, float64(width), float64(height))
		drawCanvas()
		return nil
	}))

}

func drawLine(ctx context.Context2D, x1, y1, x2, y2 float64, color string) {
	ctx.BeginPath()
	ctx.MoveTo(x1, y1)
	ctx.LineTo(x2, y2)
	ctx.StrokeStyle(color)
	ctx.LineWidth(1)
	ctx.Stroke()
}

func calcTextSize(text string, topMargin, labelWidth float64) float64 {
	var textSize = topMargin - (topMargin / 5)
	if textSize > labelWidth {
		textSize = labelWidth
	}
	if textSize*float64(len(text)) > labelWidth {
		textSize = labelWidth / float64(len(text)-1)
	}
	return textSize
}
