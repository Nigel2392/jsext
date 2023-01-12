//go:build js && wasm
// +build js,wasm

package charts

import (
	"syscall/js"

	"github.com/Nigel2392/jsext/canvas"
	"github.com/Nigel2392/jsext/canvas/context"
	"github.com/Nigel2392/jsext/framework/graphs/labels"
	"github.com/Nigel2392/jsext/framework/graphs/options"
	"github.com/Nigel2392/jsext/framework/helpers/convert"
)

func Bar(Canvas canvas.Canvas, opts options.GraphOptions) {
	// SETUP CANVAS
	///////////////////////////////////////
	var width = opts.WidthPX
	var height = opts.HeightPX
	Canvas.Width(opts.WidthPX)
	Canvas.Height(opts.HeightPX)
	var ctx = Canvas.Context2D()
	///////////////////////////////////////

	//if len(opts.Labels) != len(opts.Values) {
	//	if len(opts.Labels) != 0 {
	//		panic("Length of labels must match length of values")
	//	}
	//}

	// SETUP GRAPH
	///////////////////////////////////////
	var barWidth = float64(width / (len(opts.Values)))
	var paddingRight = barWidth / 5
	barWidth -= paddingRight

	// Set up graph border width
	var borderWidth float64
	if opts.GraphBorder {
		borderWidth = paddingRight / 8
		if borderWidth < 1 {
			borderWidth = 1
		}
	}

	// Set up bottom margin for labels
	var bottomMargin int
	tenths := height / 10
	if height-tenths < 0 {
		bottomMargin = 0
	} else {
		bottomMargin = tenths
	}
	bottomMargin /= 2

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

	var drawCanvas = func() {
		// Set up graph title
		writeTitle(ctx, width, bottomMargin, opts)
		///////////////////////////////////////
		// DRAW GRAPH
		///////////////////////////////////////
		var ct = 0
		for i, value := range opts.Values {
			var color, borderColor string
			color, ct = getColor(opts.Colors, ct, "#5555ff")
			borderColor, ct = getColor(opts.BackgroundColors, ct, "#000000")

			var barHeight = (convert.ToFloat(value) / maxHeight) * float64(height-bottomMargin*2)
			var barY = float64(height) - float64(barHeight)
			var barPositionY = barY - float64(bottomMargin)
			var barPositionX = barWidth*float64(i) + paddingRight*float64(i) + float64(paddingRight/2)

			createBar(ctx, barPositionX, barPositionY, float64(barWidth), barHeight, color, borderColor, int(borderWidth), opts.GraphBackgroundOpacity)

			// Calculate the text size.
			var textSize = float64(bottomMargin - (bottomMargin / 5))
			if textSize > barWidth {
				textSize = barWidth
			}
			if textSize*float64(len(labelOptions.Label(i))) > barWidth {
				textSize = float64(barWidth) / float64(len(labelOptions.Label(i))-1)
			}

			if opts.ShowResults {
				// Write the value on the bar.
				var textSize = float64(bottomMargin - (bottomMargin / 5))
				if textSize > barWidth {
					textSize = barWidth
				}

				ctx.Font(convert.FormatNumber(textSize) + "px Arial")
				ctx.TextAlign("center")
				ctx.TextBaseline("middle")
				ctx.FillStyle("#000000")
				ctx.FillText(convert.FormatNumber(value), barPositionX+float64(barWidth/2), float64(height-bottomMargin)-textSize)
			}

			// Set the label position
			labelOptions.SetPosition(barPositionX+float64(barWidth/2), float64(height)-float64(bottomMargin/2), int(textSize))
			// Write out the label
			labelOptions.WriteLabel(ctx, i)
		}
	}
	drawCanvas()
	///////////////////////////////////////
	Canvas.Call("addEventListener", "mousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var event = args[0]
		var x = event.Get("offsetX").Float()
		var y = event.Get("offsetY").Float()

		// Check if the mouse is over a bar.
		for i, value := range opts.Values {
			var barHeight = (convert.ToFloat(value) / maxHeight) * float64(height-bottomMargin*2)
			var barY = float64(height) - float64(barHeight)
			var barPositionY = barY - float64(bottomMargin)
			var barPositionX = barWidth*float64(i) + paddingRight*float64(i) + float64(paddingRight/2)

			if x >= barPositionX && x <= barPositionX+barWidth && y >= barPositionY && y <= barPositionY+barHeight {
				// Draw the tooltip.
				ctx.ClearRect(0, 0, float64(width), float64(height))
				drawCanvas()
				drawTooltip(ctx, width, height, x, y, labelOptions.UnderlyingLabel(i)+": "+convert.FormatNumber(opts.Values[i]))
				// Check if the value is not on any bar, but maybe between them. If the mouse is not over a bar on the canvas, clear the canvas.
			} else {
				for i := 0; i < len(opts.Values); i++ {
					var barHeight = (convert.ToFloat(opts.Values[i]) / maxHeight) * float64(height-bottomMargin*2)
					var barY = float64(height) - float64(barHeight)
					var barPositionY = barY - float64(bottomMargin)
					var barPositionX = barWidth*float64(i) + paddingRight*float64(i) + float64(paddingRight/2)
					if x >= barPositionX && x <= barPositionX+barWidth && y >= barPositionY && y <= barPositionY+barHeight {
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

func createBar(ctx context.Context2D, x, y, width, height float64, color, borderColor string, borderWidth int, opacity float64) {
	// Get the old alpha value.
	var oldAlpha = ctx.GlobalAlpha()
	if opacity <= 0 {
		opacity = 1
	}
	// Set the new alpha value.
	ctx.GlobalAlpha(opacity)
	// Draw the bar.
	ctx.BeginPath()
	ctx.FillStyle(color)
	ctx.FillRect(x, y, float64(width), height)
	if borderWidth > 0 {
		ctx.FillStyle(borderColor)
		ctx.LineWidth(float64(borderWidth))
		var borderX = x + float64(borderWidth/2)
		var borderY = y + float64(borderWidth/2)
		ctx.StrokeRect(borderX, borderY, float64(width)-float64(borderWidth), height-float64(borderWidth))
	}
	ctx.Stroke()
	// Set the alpha back to what it was.
	ctx.GlobalAlpha(oldAlpha)
}
