package charts

import (
	"strconv"

	"github.com/Nigel2392/jsext/canvas/context"
	"github.com/Nigel2392/jsext/framework/graphs/options"
)

//	var ct = 0
//	for i, value := range opts.Values {
//		var color, borderColor string
//		color, ct = getColor(opts.Colors, ct, "#5555ff")
//
// Get the color for a given chart index.
func getColor(colors []string, index int, def string) (string, int) {
	if len(colors) > 0 {
		if index >= len(colors) {
			index = 0
		}
		return colors[index], index + 1
	}
	return def, index
}

func writeTitle(ctx context.Context2D, width, topMargin int, opts options.GraphOptions) {
	if opts.GraphTitle != "" {
		// Draw Title
		var titleSize = topMargin - (topMargin / 5)
		ctx.FillStyle("#000000")
		ctx.Font(strconv.Itoa(int(titleSize)) + "px Arial")
		ctx.TextAlign("center")
		ctx.TextBaseline("middle")
		ctx.FillText(opts.GraphTitle, float64(width/2), float64(topMargin/2))
	}

}

func drawTooltip(ctx context.Context2D, width, height int, x, y float64, text string) {
	var tooltipWidth = float64(width / 4)
	var tooltipHeight = float64(height / 10)
	var tooltipX = x - float64(tooltipWidth/2)
	var tooltipY = y - float64(tooltipHeight*2)
	if tooltipX < 0 {
		tooltipX = 0
	}
	if tooltipX+tooltipWidth > float64(width) {
		tooltipX = float64(width) - tooltipWidth
	}
	if tooltipY < 0 {
		tooltipY = 0
	}
	if tooltipY+tooltipHeight > float64(height) {
		tooltipY = float64(height) - tooltipHeight
	}

	var fontSize = 40
	if fontSize*len(text) > int(tooltipWidth) {
		fontSize = int(tooltipWidth) / len(text)
	}

	ctx.BeginPath()
	ctx.FillStyle("rgba(0,0,0, 0.8)")
	ctx.FillRect(tooltipX, tooltipY, tooltipWidth, tooltipHeight)
	ctx.FillStyle("#ffffff")
	ctx.Font(strconv.Itoa(fontSize) + "px Arial")
	ctx.TextAlign("center")
	ctx.TextBaseline("middle")
	ctx.FillText(text, tooltipX+float64(tooltipWidth/2), tooltipY+float64(tooltipHeight/2))
}

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64
}

func Min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
