//go:build js && wasm
// +build js,wasm

package csshelpers

import (
	"fmt"
	"math"
)

var hexMap = map[string]string{
	"AliceBlue":            "#F0F8FF",
	"AntiqueWhite":         "#FAEBD7",
	"Aqua":                 "#00FFFF",
	"Aquamarine":           "#7FFFD4",
	"Azure":                "#F0FFFF",
	"Beige":                "#F5F5DC",
	"Bisque":               "#FFE4C4",
	"Black":                "#000000",
	"BlanchedAlmond":       "#FFEBCD",
	"Blue":                 "#0000FF",
	"BlueViolet":           "#8A2BE2",
	"Brown":                "#A52A2A",
	"BurlyWood":            "#DEB887",
	"CadetBlue":            "#5F9EA0",
	"Chartreuse":           "#7FFF00",
	"Chocolate":            "#D2691E",
	"Coral":                "#FF7F50",
	"CornflowerBlue":       "#6495ED",
	"Cornsilk":             "#FFF8DC",
	"Crimson":              "#DC143C",
	"Cyan":                 "#00FFFF",
	"DarkBlue":             "#00008B",
	"DarkCyan":             "#008B8B",
	"DarkGoldenRod":        "#B8860B",
	"DarkGrey":             "#A9A9A9",
	"DarkGreen":            "#006400",
	"DarkKhaki":            "#BDB76B",
	"DarkMagenta":          "#8B008B",
	"DarkOliveGreen":       "#556B2F",
	"Darkorange":           "#FF8C00",
	"DarkOrchid":           "#9932CC",
	"DarkRed":              "#8B0000",
	"DarkSalmon":           "#E9967A",
	"DarkSeaGreen":         "#8FBC8F",
	"DarkSlateBlue":        "#483D8B",
	"DarkSlateGrey":        "#2F4F4F",
	"DarkTurquoise":        "#00CED1",
	"DarkViolet":           "#9400D3",
	"DeepPink":             "#FF1493",
	"DeepSkyBlue":          "#00BFFF",
	"DimGray":              "#696969",
	"DodgerBlue":           "#1E90FF",
	"FireBrick":            "#B22222",
	"FloralWhite":          "#FFFAF0",
	"ForestGreen":          "#228B22",
	"Fuchsia":              "#FF00FF",
	"Gainsboro":            "#DCDCDC",
	"GhostWhite":           "#F8F8FF",
	"Gold":                 "#FFD700",
	"GoldenRod":            "#DAA520",
	"Grey":                 "#808080",
	"Green":                "#008000",
	"GreenYellow":          "#ADFF2F",
	"HoneyDew":             "#F0FFF0",
	"HotPink":              "#FF69B4",
	"IndianRed":            "#CD5C5C",
	"Indigo":               "#4B0082",
	"Ivory":                "#FFFFF0",
	"Khaki":                "#F0E68C",
	"Lavender":             "#E6E6FA",
	"LavenderBlush":        "#FFF0F5",
	"LawnGreen":            "#7CFC00",
	"LemonChiffon":         "#FFFACD",
	"LightBlue":            "#ADD8E6",
	"LightCoral":           "#F08080",
	"LightCyan":            "#E0FFFF",
	"LightGoldenRodYellow": "#FAFAD2",
	"LightGrey":            "#D3D3D3",
	"LightGreen":           "#90EE90",
	"LightPink":            "#FFB6C1",
	"LightSalmon":          "#FFA07A",
	"LightSeaGreen":        "#20B2AA",
	"LightSkyBlue":         "#87CEFA",
	"LightSlateGrey":       "#778899",
	"LightSteelBlue":       "#B0C4DE",
	"LightYellow":          "#FFFFE0",
	"Lime":                 "#00FF00",
	"LimeGreen":            "#32CD32",
	"Linen":                "#FAF0E6",
	"Magenta":              "#FF00FF",
	"Maroon":               "#800000",
	"MediumAquaMarine":     "#66CDAA",
	"MediumBlue":           "#0000CD",
	"MediumOrchid":         "#BA55D3",
	"MediumPurple":         "#9370D8",
	"MediumSeaGreen":       "#3CB371",
	"MediumSlateBlue":      "#7B68EE",
	"MediumSpringGreen":    "#00FA9A",
	"MediumTurquoise":      "#48D1CC",
	"MediumVioletRed":      "#C71585",
	"MidnightBlue":         "#191970",
	"MintCream":            "#F5FFFA",
	"MistyRose":            "#FFE4E1",
	"Moccasin":             "#FFE4B5",
	"NavajoWhite":          "#FFDEAD",
	"Navy":                 "#000080",
	"OldLace":              "#FDF5E6",
	"Olive":                "#808000",
	"OliveDrab":            "#6B8E23",
	"Orange":               "#FFA500",
	"OrangeRed":            "#FF4500",
	"Orchid":               "#DA70D6",
	"PaleGoldenRod":        "#EEE8AA",
	"PaleGreen":            "#98FB98",
	"PaleTurquoise":        "#AFEEEE",
	"PaleVioletRed":        "#D87093",
	"PapayaWhip":           "#FFEFD5",
	"PeachPuff":            "#FFDAB9",
	"Peru":                 "#CD853F",
	"Pink":                 "#FFC0CB",
	"Plum":                 "#DDA0DD",
	"PowderBlue":           "#B0E0E6",
	"Purple":               "#800080",
	"Red":                  "#FF0000",
	"RosyBrown":            "#BC8F8F",
	"RoyalBlue":            "#4169E1",
	"SaddleBrown":          "#8B4513",
	"Salmon":               "#FA8072",
	"SandyBrown":           "#F4A460",
	"SeaGreen":             "#2E8B57",
	"SeaShell":             "#FFF5EE",
	"Sienna":               "#A0522D",
	"Silver":               "#C0C0C0",
	"SkyBlue":              "#87CEEB",
	"SlateBlue":            "#6A5ACD",
	"SlateGrey":            "#708090",
	"Snow":                 "#FFFAFA",
	"SpringGreen":          "#00FF7F",
	"SteelBlue":            "#4682B4",
	"Tan":                  "#D2B48C",
	"Teal":                 "#008080",
	"Thistle":              "#D8BFD8",
	"Tomato":               "#FF6347",
	"Turquoise":            "#40E0D0",
	"Violet":               "#EE82EE",
	"Wheat":                "#F5DEB3",
	"White":                "#FFFFFF",
	"WhiteSmoke":           "#F5F5F5",
	"Yellow":               "#FFFF00",
	"YellowGreen":          "#9ACD32",
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
	if val, ok := hexMap[scol]; ok {
		scol = val
	}
	format := "#%02x%02x%02x"
	factor := 1.0 / 255.0
	if len(scol) == 4 {
		format = "#%1x%1x%1x"
		factor = 1.0 / 15.0
	}

	var r, g, b uint8
	n, err := fmt.Sscanf(scol, format, &r, &g, &b)
	if err != nil {
		return Color{}, err
	}
	if n != 3 {
		return Color{}, fmt.Errorf("color: %v is not a hex-color", scol)
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
