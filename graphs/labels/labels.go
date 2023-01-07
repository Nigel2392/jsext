package labels

import (
	"strconv"

	"github.com/Nigel2392/jsext/canvas/context"
	"github.com/Nigel2392/jsext/graphs/convert"
)

type LabelOpts struct {
	textPositionX float64
	textPositionY float64
	textSize      int
	Labels        []string
}

func (l *LabelOpts) SetLabels(labels []string, values []any) {
	l.Labels = make([]string, len(values))
	if len(labels) > 0 && len(labels) == len(values) {
		l.Labels = labels
	} else {
		for i := 0; i < len(values); i++ {
			l.Labels[i] = strconv.FormatFloat(convert.ToFloat(values[i]), 'f', 2, 64)
		}
	}
}

func (l *LabelOpts) SetPosition(x, y float64, textSize int) {
	l.textPositionX = x
	l.textPositionY = y
	l.textSize = textSize
}

func (l *LabelOpts) WriteLabel(ctx context.Context2D, i int) {
	// Draw Text
	ctx.FillStyle("black")
	// Calculate text size
	ctx.Font(strconv.Itoa(l.textSize) + "px Arial")
	ctx.TextAlign("center")
	ctx.TextBaseline("middle")
	ctx.FillText(l.Labels[i], l.textPositionX, l.textPositionY)
}

func (l *LabelOpts) Label(i int) string {
	return l.Labels[i]
}
