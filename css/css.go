//go:build js && wasm
// +build js,wasm

package css

import "strconv"

// General colors for use in CSS. Can be overridden.
var (
	COLOR_MAIN        = "white"           // Main color - white
	COLOR_ONE         = "#222e50"         // COLOR_ONE - dark blue
	COLOR_ONE_LIGHTER = "#334679"         // COLOR_ONE_LIGHTER - lighter dark blue
	COLOR_TWO         = "#e9d985"         // COLOR_TWO - yellow
	COLOR_THREE       = "#439a86"         // COLOR_THREE - green
	COLOR_FOUR        = "#DADFF7"         // COLOR_FOUR - light blue
	COLOR_FIVE        = "#007991"         // COLOR_FIVE - dark green
	COLOR_SIX         = "#ff0057"         // COLOR_SIX - red
	BACKGROUND_COLOR  = "rgba(0,0,0,0.5)" // Background color (for modals, etc.)
)

// Predefined colors
const (
	COLOR_BLACK      = "#000000" // Black
	COLOR_GREY       = "#575757" // Grey
	COLOR_DARK_GRAY  = "#252525" // Dark Grey
	COLOR_LIGHT_GRAY = "#9d9d9d" // Light Grey
	COLOR_RED        = "#ff0000" // Red
	COLOR_GREEN      = "#00ff00" // Green
	COLOR_BLUE       = "#0000ff" // Blue
	COLOR_YELLOW     = "#ffff00" // Yellow
	COLOR_PURPLE     = "#ff00ff" // Purple
	COLOR_CYAN       = "#00ffff" // Cyan
	COLOR_ORANGE     = "#ffa500" // Orange
	COLOR_BROWN      = "#a52a2a" // Brown
	COLOR_PINK       = "#ffc0cb" // Pink
	COLOR_WHITE      = "#ffffff" // White
	COLOR_LIME       = "#00ff00" // Lime
	COLOR_MAROON     = "#800000" // Maroon
	COLOR_NAVY       = "#000080" // Navy
	COLOR_OLIVE      = "#808000" // Olive
	COLOR_TEAL       = "#008080" // Teal
	COLOR_AQUA       = "#00ffff" // Aqua
	COLOR_FUCHSIA    = "#ff00ff" // Fuchsia
	COLOR_SILVER     = "#c0c0c0" // Silver
)

// Convert an interger and append "px".
func ToPX(px int, end_css_statement ...bool) string {
	var s = strconv.Itoa(px) + "px"
	var end_CSS = len(end_css_statement) > 0 && end_css_statement[0]
	if end_CSS {
		// god I fucking hate concatenating strings,
		// is a buffer worth?
		return s + ";"
	}
	return s
}

// Convert an interger and append "%".
func ToPercent(percent int, end_css_statement ...bool) string {
	var s = strconv.Itoa(percent) + "%"
	var end_CSS = len(end_css_statement) > 0 && end_css_statement[0]
	if end_CSS {
		// god I fucking hate concatenating strings,
		// is a buffer worth?
		return s + ";"
	}
	return s
}

// Convert an interger and append "rem".
func ToRem(rem int, end_css_statement ...bool) string {
	var s = strconv.Itoa(rem) + "rem"
	var end_CSS = len(end_css_statement) > 0 && end_css_statement[0]
	if end_CSS {
		// god I fucking hate concatenating strings,
		// is a buffer worth?
		return s + ";"
	}
	return s
}

// Convert an interger and append "vw".
func ToVW(vw int, end_css_statement ...bool) string {
	var s = strconv.Itoa(vw) + "vw"
	var end_CSS = len(end_css_statement) > 0 && end_css_statement[0]
	if end_CSS {
		// god I fucking hate concatenating strings,
		// is a buffer worth?
		return s + ";"
	}
	return s
}

// Convert an interger and append "vh".
func ToVH(vh int, end_css_statement ...bool) string {
	var s = strconv.Itoa(vh) + "vh"
	var end_CSS = len(end_css_statement) > 0 && end_css_statement[0]
	if end_CSS {
		// god I fucking hate concatenating strings,
		// is a buffer worth?
		return s + ";"
	}
	return s
}
