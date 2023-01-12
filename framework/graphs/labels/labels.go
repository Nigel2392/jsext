//go:build js && wasm
// +build js,wasm

package labels

import (
	"strconv"

	"github.com/Nigel2392/jsext/canvas/context"
	"github.com/Nigel2392/jsext/framework/helpers/convert"
)

type LabelOpts struct {
	textPositionX float64
	textPositionY float64
	textSize      int
	Labels        []string
	UnderLying    []string
}

// Set the labels for a graph.
func (l *LabelOpts) SetLabels(labels []string, values []any) {
	l.Labels = make([]string, len(values))
	l.UnderLying = make([]string, len(values))
	if len(labels) > 0 && len(labels) == len(values) {
		l.Labels = labels
		l.UnderLying = labels
	} else {
		// If the labels are x times the length of the values, fill x spots in between the list of labels
		// with the string zero value.
		if len(labels) > 0 && len(values)%len(labels) == 0 {
			var labelMod = len(values) / len(labels)
			var labelIndex = 0
			var lastLabel = ""
			for i := 0; i < len(values); i++ {
				if i%labelMod == 0 {
					lastLabel = labels[labelIndex]
					l.Labels[i] = lastLabel
					l.UnderLying[i] = lastLabel
					labelIndex++
				} else {
					l.Labels[i] = ""
					l.UnderLying[i] = lastLabel
				}
			}
		} else {
			for i := 0; i < len(values); i++ {
				l.Labels[i] = strconv.FormatFloat(convert.ToFloat(values[i]), 'f', 2, 64)
				l.UnderLying[i] = labels[i]
			}
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

func (l *LabelOpts) UnderlyingLabel(i int) string {
	return l.UnderLying[i]
}
