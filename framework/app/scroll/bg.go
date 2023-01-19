package scroll

import (
	"fmt"
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
	Gradients    []string
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

func (b Backgrounds) AddGradient(gradientType GradientType, g ...string) {
	for _, bg := range b {
		bg.AddGradient(gradientType, g...)
	}
}

// Gradient type can only be set once!
// Gradients will be passed along in the following css string:
// background: linear-gradient(to boottom, gradient1, gradient2, ...);
func (b *Background) AddGradient(gradientType GradientType, g ...string) {
	if b.Gradient == nil {
		b.Gradient = &Gradient{
			GradientType: gradientType,
		}
	}
	b.Gradient.Gradients = append(b.Gradient.Gradients, g...)
}

func (bg *Background) CSS(selector string, gradientTo string) string {
	var css string
	var gradientTyp string = "linear"
	switch bg.Gradient.GradientType {
	case GradientTypeLinear:
		gradientTyp = "linear"
	case GradientTypeRadial:
		gradientTyp = "radial"
	}
	switch bg.BackgroundType {
	case BackgroundTypeImage:
		css += selector + ` {
			background-image: ` + gradientTyp + `-gradient(` + gradientTo + `, ` + bg.gradient() + `), url('` + bg.Background + `');
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
				background: ` + gradientTyp + `-gradient(` + gradientTo + `, ` + bg.gradient() + `), ` + bg.Background + `;
				` + bg.ExtraCSS + `
			}`
	case BackgroundTypeStyle:
		css += selector + ` {
				` + bg.Background + `;
				` + bg.ExtraCSS + `
			}`
	}
	fmt.Println(css)
	return css
}

func (b *Background) gradient() string {
	if len(b.Gradient.Gradients) < 2 {
		return "rgba(0,0,0,0), rgba(0,0,0,0)"
	}
	return strings.Join(b.Gradient.Gradients, ", ")
}
