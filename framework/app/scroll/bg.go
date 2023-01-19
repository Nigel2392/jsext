package scroll

import (
	"strings"
)

// Supported background types for the scrollable app.
type BackgroundType int8

// Type of gradient to use.
type GradientType int8

// Supported background types for the scrollable app.
const (
	BackgroundTypeImage BackgroundType = 1
	BackgroundTypeColor BackgroundType = 2
	BackgroundTypeStyle BackgroundType = 3
)

// Type of gradient to use.
const (
	GradientTypeLinear GradientType = 1
	GradientTypeRadial GradientType = 2
)

// Gradient is a struct that holds the type of gradient and the gradients.
type Gradient struct {
	GradientType GradientType
	Direction    string
	Gradients    []string
}

func (g *Gradient) gradientType() string {
	switch g.GradientType {
	case GradientTypeLinear:
		return "linear"
	case GradientTypeRadial:
		return "radial"
	}
	return "linear"
}

func (g *Gradient) direction() string {
	if g.Direction == "" {
		switch g.GradientType {
		case GradientTypeLinear:
			return "to bottom"
		case GradientTypeRadial:
			return "circle at center"
		}
		return "to bottom"
	}
	return g.Direction
}

func (g *Gradient) String() string {
	var b strings.Builder
	typ := g.gradientType()
	grad := "-gradient("
	dir := g.direction()
	comma := ", "
	fallback := "rgba(0,0,0,0), rgba(0,0,0,0))"
	var totalLen int = len(typ) + len(grad) + len(dir) + len(comma)
	switch l := len(g.Gradients); l {
	case 0, 1:
		totalLen += len(fallback)
	default:
		totalLen += (len(g.Gradients) - 1) * len(", ")
		for _, gradient := range g.Gradients {
			totalLen += len(gradient)
		}
		totalLen += len(")")
	}
	b.Grow(totalLen)
	b.WriteString(g.gradientType())
	b.WriteString("-gradient(")
	b.WriteString(g.direction())
	b.WriteString(", ")
	if len(g.Gradients) == 0 {
		b.WriteString("rgba(0,0,0,0), rgba(0,0,0,0))")
		return b.String()
	}
	for i, gradient := range g.Gradients {
		b.WriteString(gradient)
		if i < len(g.Gradients)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(")")
	return b.String()
}

// Background is a struct that holds the type of background and the background styling.
type Background struct {
	BackgroundType BackgroundType
	Background     string
	ExtraCSS       string
	Gradient       *Gradient
}

// Backgrounds is a slice of Backgrounds.
type Backgrounds []*Background

func (b Backgrounds) AddGradient(gradientType GradientType, direction string, g ...string) {
	for _, bg := range b {
		bg.AddGradient(gradientType, direction, g...)
	}
}

// Gradient type can only be set once!
// Gradients will be passed along in the following css string:
// background: linear-gradient(to boottom, gradient1, gradient2, ...);
func (b *Background) AddGradient(gradientType GradientType, direction string, g ...string) {
	if b.Gradient == nil {
		b.Gradient = &Gradient{
			Gradients: make([]string, 0),
		}
	}
	b.Gradient.Gradients = append(b.Gradient.Gradients, g...)
	b.Gradient.GradientType = gradientType
	b.Gradient.Direction = direction
}

func (bg *Background) CSS(selector string) string {
	var css string
	switch bg.BackgroundType {
	case BackgroundTypeImage:
		css += selector + ` {
			background-image: ` + bg.Gradient.String() + `, url('` + bg.Background + `');
			background-repeat: no-repeat;
			background-position: center;
			background-size: cover;
			` + bg.ExtraCSS + `
		}`
	case BackgroundTypeColor:
		if bg.Background == "" && len(bg.Gradient.Gradients) > 0 {
			bg.Background = "rgba(0,0,0,0)"
		}
		css += selector + ` {
				background:` + bg.Gradient.String() + `, ` + bg.Background + `;
				` + bg.ExtraCSS + `
			}`
	case BackgroundTypeStyle:
		css += selector + ` {
				` + bg.Background + `;
				` + bg.ExtraCSS + `
			}`
	}
	return css
}
