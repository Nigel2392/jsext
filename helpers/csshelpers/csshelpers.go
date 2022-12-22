package csshelpers

import (
	"fmt"
	"math"
)

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
func Hex(scol string) (Color, error) {
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

func RGB255(r, g, b uint8) Color {
	return Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
}

func (col *Color) Complementary() Color {
	var col2 = HueOffset(*col, 180)
	return col2
}

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

func ConvertColor(col Color) string {
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
