package shortcuts

import (
	"fmt"

	"github.com/Nigel2392/jsext/v2/css"
	"github.com/Nigel2392/jsext/v2/jse"
	"github.com/Nigel2392/jsext/v2/jsrand"
)

type HR struct {
	// Width and height of the HR
	Width, Height string
	// Main colors to use.
	BackCol, FadeCol string
	// Margin top and bottom.
	MarginTopBottom string
	// Border radius
	BorderRadius string
	// Fade direction, can be "left", "right", "top", "bottom"
	FadeDir string
	// Opacity of the hr
	Opacity float64
	// Format strings for the colors
	// Example: "rgba(%d, %d, %d, 1)"
	BackColFormat string
	// Format strings for the colors
	// Example: "rgba(%d, %d, %d, 0.1)"
	FadeColFormat string
	// Animate
	Animate bool
}

func (h *HR) setDefaults() {
	if h.Width == "" {
		h.Width = "100%"
	}
	if h.Height == "" {
		h.Height = "1px"
	}
	if h.BackCol == "" {
		h.BackCol = "#000000"
	}
	if h.BackColFormat == "" {
		h.BackColFormat = "rgba(%d, %d, %d, 1)"
	}
	if h.FadeColFormat == "" {
		h.FadeColFormat = "rgba(%d, %d, %d, 0.1)"
	}
	if h.FadeCol == "" {
		h.FadeCol = "#ffffff"
	}
	if h.MarginTopBottom == "" {
		h.MarginTopBottom = "0"
	}
	if h.FadeDir == "" {
		h.FadeDir = "right"
	}
	if h.BorderRadius == "" {
		h.BorderRadius = "0"
	}
	if h.Opacity == 0 {
		h.Opacity = 0.25
	}
	if h.Opacity > 1 && h.Opacity <= 100 {
		h.Opacity /= 100
	}
}

func FancyHR(opts *HR) *jse.Element {
	opts.setDefaults()
	var mainColor = css.FormatRGBA(opts.BackCol, opts.BackColFormat)
	var fadeColor = css.FormatRGBA(opts.FadeCol, opts.FadeColFormat)
	var hash = jsrand.XorBytes(opts.Width + opts.Height + opts.BackCol + opts.FadeCol + opts.MarginTopBottom + opts.FadeDir)
	var hr = jse.NewElement("hr")
	hr.ClassList("fancy-hr" + hash)
	hr.AttrID("fancy-hr" + hash)

	var animation string
	if opts.Animate {
		animation = `@keyframes gradientAnimation` + hash + ` {
			0% { background-position: 0% 50%; }
			50% { background-position: 100% 50%; }
			100% { background-position: 0% 50%; }
		}`
	}
	var animation_css string
	if opts.Animate {
		animation_css = `animation: gradientAnimation` + hash + ` 5s ease infinite;`
	}

	hr.StyleBlock(`
		.fancy-hr` + hash + ` {
			width: ` + opts.Width + `;
			border: 0;
			height: ` + opts.Height + `;
			background: linear-gradient(to ` + opts.FadeDir + `, ` + fadeColor + ` 10%, ` + mainColor + ` 50%, ` + fadeColor + ` 90%);
			margin: ` + opts.MarginTopBottom + ` 0;
			background-size: 200% 200%;
			opacity: ` + fmt.Sprintf("%.2f", opts.Opacity) + `;
			border-radius: ` + opts.BorderRadius + `;
			` + animation_css + `
		}` + animation)

	return hr
}
