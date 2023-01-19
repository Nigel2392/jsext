package scroll

import (
	"fmt"
	"strings"
)

type Background struct {
	BackgroundType BackgroundType
	Background     string
	ExtraCSS       string
	gradients      []string
}

type Backgrounds []*Background

func (b Backgrounds) AddGradient(g ...string) {
	for _, bg := range b {
		bg.AddGradient(g...)
	}
}

// Gradients can only be added once!
// Gradients will be passed along in the following css string:
// background: linear-gradient(to boottom, gradient1, gradient2, ...);
func (b *Background) AddGradient(g ...string) {
	b.gradients = append(b.gradients, g...)
}

func (bg *Background) CSS(selector string, gradientTo string) string {
	var css string
	switch bg.BackgroundType {
	case BackgroundTypeImage:
		css += selector + ` {
			background-image: linear-gradient(` + gradientTo + `, ` + bg.gradient() + `), url('` + bg.Background + `');
			background-repeat: no-repeat;
			background-position: center;
			background-size: cover;
			` + bg.ExtraCSS + `
		}`
	case BackgroundTypeColor:
		if bg.Background == "" && len(bg.gradients) > 0 {
			bg.Background = "rgba(0,0,0,0)"
		}
		css += selector + ` {
				background: linear-gradient(` + gradientTo + `, ` + bg.gradient() + `), ` + bg.Background + `;
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
	if len(b.gradients) < 2 {
		return "rgba(0,0,0,0), rgba(0,0,0,0)"
	}
	return strings.Join(b.gradients, ", ")
}
