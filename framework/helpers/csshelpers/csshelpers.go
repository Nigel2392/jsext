//go:build js && wasm
// +build js,wasm

package csshelpers

import (
	"fmt"
	"math"
)

var hexMap = map[string]string{
	"aliceblue":            "#f0f8ff",
	"antiquewhite":         "#faebd7",
	"aqua":                 "#00ffff",
	"aquamarine":           "#7fffd4",
	"azure":                "#f0ffff",
	"beige":                "#f5f5dc",
	"bisque":               "#ffe4c4",
	"black":                "#000000",
	"blanchedalmond":       "#ffebcd",
	"blue":                 "#0000ff",
	"blueviolet":           "#8a2be2",
	"brown":                "#a52a2a",
	"burlywood":            "#deb887",
	"cadetblue":            "#5f9ea0",
	"chartreuse":           "#7fff00",
	"chocolate":            "#d2691e",
	"coral":                "#ff7f50",
	"cornflowerblue":       "#6495ed",
	"cornsilk":             "#fff8dc",
	"crimson":              "#dc143c",
	"cyan":                 "#00ffff",
	"darkblue":             "#00008b",
	"darkcyan":             "#008b8b",
	"darkgoldenrod":        "#b8860b",
	"darkgrey":             "#a9a9a9",
	"darkgreen":            "#006400",
	"darkkhaki":            "#bdb76b",
	"darkmagenta":          "#8b008b",
	"darkolivegreen":       "#556b2f",
	"darkorange":           "#ff8c00",
	"darkorchid":           "#9932cc",
	"darkred":              "#8b0000",
	"darksalmon":           "#e9967a",
	"darkseagreen":         "#8fbc8f",
	"darkslateblue":        "#483d8b",
	"darkslategrey":        "#2f4f4f",
	"darkturquoise":        "#00ced1",
	"darkviolet":           "#9400d3",
	"deeppink":             "#ff1493",
	"deepskyblue":          "#00bfff",
	"dimgray":              "#696969",
	"dodgerblue":           "#1e90ff",
	"firebrick":            "#b22222",
	"floralwhite":          "#fffaf0",
	"forestgreen":          "#228b22",
	"fuchsia":              "#ff00ff",
	"gainsboro":            "#dcdcdc",
	"ghostwhite":           "#f8f8ff",
	"gold":                 "#ffd700",
	"goldenrod":            "#daa520",
	"grey":                 "#808080",
	"green":                "#008000",
	"greenyellow":          "#adff2f",
	"honeydew":             "#f0fff0",
	"hotpink":              "#ff69b4",
	"indianred":            "#cd5c5c",
	"indigo":               "#4b0082",
	"ivory":                "#fffff0",
	"khaki":                "#f0e68c",
	"lavender":             "#e6e6fa",
	"lavenderblush":        "#fff0f5",
	"lawngreen":            "#7cfc00",
	"lemonchiffon":         "#fffacd",
	"lightblue":            "#add8e6",
	"lightcoral":           "#f08080",
	"lightcyan":            "#e0ffff",
	"lightgoldenrodyellow": "#fafad2",
	"lightgrey":            "#d3d3d3",
	"lightgreen":           "#90ee90",
	"lightpink":            "#ffb6c1",
	"lightsalmon":          "#ffa07a",
	"lightseagreen":        "#20b2aa",
	"lightskyblue":         "#87cefa",
	"lightslategrey":       "#778899",
	"lightsteelblue":       "#b0c4de",
	"lightyellow":          "#ffffe0",
	"lime":                 "#00ff00",
	"limegreen":            "#32cd32",
	"linen":                "#faf0e6",
	"magenta":              "#ff00ff",
	"maroon":               "#800000",
	"mediumaquamarine":     "#66cdaa",
	"mediumblue":           "#0000cd",
	"mediumorchid":         "#ba55d3",
	"mediumpurple":         "#9370d8",
	"mediumseagreen":       "#3cb371",
	"mediumslateblue":      "#7b68ee",
	"mediumspringgreen":    "#00fa9a",
	"mediumturquoise":      "#48d1cc",
	"mediumvioletred":      "#c71585",
	"midnightblue":         "#191970",
	"mintcream":            "#f5fffa",
	"mistyrose":            "#ffe4e1",
	"moccasin":             "#ffe4b5",
	"navajowhite":          "#ffdead",
	"navy":                 "#000080",
	"oldlace":              "#fdf5e6",
	"olive":                "#808000",
	"olivedrab":            "#6b8e23",
	"orange":               "#ffa500",
	"orangered":            "#ff4500",
	"orchid":               "#da70d6",
	"palegoldenrod":        "#eee8aa",
	"palegreen":            "#98fb98",
	"paleturquoise":        "#afeeee",
	"palevioletred":        "#d87093",
	"papayawhip":           "#ffefd5",
	"peachpuff":            "#ffdab9",
	"peru":                 "#cd853f",
	"pink":                 "#ffc0cb",
	"plum":                 "#dda0dd",
	"powderblue":           "#b0e0e6",
	"purple":               "#800080",
	"red":                  "#ff0000",
	"rosybrown":            "#bc8f8f",
	"royalblue":            "#4169e1",
	"saddlebrown":          "#8b4513",
	"salmon":               "#fa8072",
	"sandybrown":           "#f4a460",
	"seagreen":             "#2e8b57",
	"seashell":             "#fff5ee",
	"sienna":               "#a0522d",
	"silver":               "#c0c0c0",
	"skyblue":              "#87ceeb",
	"slateblue":            "#6a5acd",
	"slategrey":            "#708090",
	"snow":                 "#fffafa",
	"springgreen":          "#00ff7f",
	"steelblue":            "#4682b4",
	"tan":                  "#d2b48c",
	"teal":                 "#008080",
	"thistle":              "#d8bfd8",
	"tomato":               "#ff6347",
	"turquoise":            "#40e0d0",
	"violet":               "#ee82ee",
	"wheat":                "#f5deb3",
	"white":                "#ffffff",
	"whitesmoke":           "#f5f5f5",
	"yellow":               "#ffff00",
	"yellowgreen":          "#9acd32",
}

// Credit where credit is due.
//
// Needed a simple package to get a color's complementary color.

// https://github.com/lucasb-eyer/go-colorful
type Color struct {
	R float64
	G float64
	B float64
}

// https://github.com/lucasb-eyer/go-colorful
// Hex returns the hex "html" representation of the color, as in #ff0080.
func (col Color) Hex() string {
	// Add 0.5 for rounding
	return fmt.Sprintf("#%02x%02x%02x", uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
}

// https://github.com/lucasb-eyer/go-colorful
// Hex parses a "html" hex color-string, either in the 3 "#f0c" or 6 "#ff1034" digits form.
// It also parses the named colors from the CSS3 specification.
func Hex(scol string) (Color, error) {
	var c string
	if val, ok := hexMap[scol]; ok {
		c = val
	} else {
		c = scol
	}
	format := "#%02x%02x%02x"
	factor := 1.0 / 255.0
	if len(c) == 4 {
		format = "#%1x%1x%1x"
		factor = 1.0 / 15.0
	}

	var r, g, b uint8
	n, err := fmt.Sscanf(c, format, &r, &g, &b)
	if err != nil {
		return Color{}, err
	}
	if n != 3 {
		return Color{}, fmt.Errorf("color: %v is not a hex-color", c)
	}
	return Color{float64(r) * factor, float64(g) * factor, float64(b) * factor}, nil
}

// !WARNING! Only use this if you are sure the color is valid.
func UnsafeHex(scol string) Color {
	col, err := Hex(scol)
	if err != nil {
		panic(err)
	}
	return col
}

// https://github.com/lucasb-eyer/go-colorful
func (col Color) Hsv() (h, s, v float64) {
	min := math.Min(math.Min(col.R, col.G), col.B)
	v = math.Max(math.Max(col.R, col.G), col.B)
	C := v - min

	s = 0.0
	if v != 0.0 {
		s = C / v
	}

	h = 0.0 // We use 0 instead of undefined as in wp.
	if min != v {
		if v == col.R {
			h = math.Mod((col.G-col.B)/C, 6.0)
		}
		if v == col.G {
			h = (col.B-col.R)/C + 2.0
		}
		if v == col.B {
			h = (col.R-col.G)/C + 4.0
		}
		h *= 60.0
		if h < 0.0 {
			h += 360.0
		}
	}
	return
}

// https://github.com/lucasb-eyer/go-colorful
// Implement the Go color.Color interface.
func (col Color) RGBA() (r, g, b, a uint32) {
	r = uint32(col.R*65535.0 + 0.5)
	g = uint32(col.G*65535.0 + 0.5)
	b = uint32(col.B*65535.0 + 0.5)
	a = 0xFFFF
	return
}

// https://github.com/lucasb-eyer/go-colorful
// Might come in handy sometimes to reduce boilerplate code.
func (col Color) RGB255() (r, g, b uint8) {
	r = uint8(col.R*255.0 + 0.5)
	g = uint8(col.G*255.0 + 0.5)
	b = uint8(col.B*255.0 + 0.5)
	return
}

// Convert RGB values to a Color struct.
func RGB255(r, g, b uint8) Color {
	return Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
}

// Get the complementary color struct.
func (col *Color) Complementary() Color {
	var col2 = HueOffset(*col, 180)
	return col2
}

// Get the complementary color of a string.
func Complementary(col string) (string, error) {
	c, err := Hex(col)
	if err != nil {
		return "", err
	}
	return c.Complementary().Hex(), nil
}

// Triadic returns the triadic values for any given color
func (c *Color) Triadic() []Color {
	var cc []Color
	cc = append(cc, HueOffset(*c, -120))
	cc = append(cc, HueOffset(*c, 120))
	return cc
}

// https://github.com/lucasb-eyer/go-colorful
// Hsv creates a new Color given a Hue in [0..360], a Saturation and a Value in [0..1]
func Hsv(H, S, V float64) Color {
	Hp := H / 60.0
	C := V * S
	X := C * (1.0 - math.Abs(math.Mod(Hp, 2.0)-1.0))

	m := V - C
	r, g, b := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= Hp && Hp < 1.0:
		r = C
		g = X
	case 1.0 <= Hp && Hp < 2.0:
		r = X
		g = C
	case 2.0 <= Hp && Hp < 3.0:
		g = C
		b = X
	case 3.0 <= Hp && Hp < 4.0:
		g = X
		b = C
	case 4.0 <= Hp && Hp < 5.0:
		r = X
		b = C
	case 5.0 <= Hp && Hp < 6.0:
		r = C
		b = X
	}

	return Color{m + r, m + g, m + b}
}

// https://github.com/muesli/gamut
func HueOffset(col Color, degrees int) Color {
	h, s, v := col.Hsv()
	h += float64(degrees)
	if h < 0 {
		h += 360
	} else if h >= 360 {
		h -= 360
	}

	return Hsv(h, s, v)
}

func SeenColor(col Color) string {
	if r, g, b := col.RGB255(); r == 0 && g == 0 && b == 0 {
		return "#FFFFFF"
	} else {
		if r > 127 && g > 127 && b > 127 {
			return "#000000"
		} else {
			return "#FFFFFF"
		}
	}
}

// If the right format is not given, we will return the color as is.
func FormatRGBA(color, format string) string {
	var ret string
	var col, err = Hex(color)
	if err == nil {
		var r, g, b = col.RGB255()
		ret = fmt.Sprintf(format, r, g, b)
	} else {
		ret = color
	}
	return ret
}
