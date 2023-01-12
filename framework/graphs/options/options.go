//go:build js && wasm
// +build js,wasm

package options

type GraphTypes int

const (
	Bar GraphTypes = iota
	Line
	Pie
	Donut
)

type GraphOptions struct {
	WidthPX  int
	HeightPX int
	// Graph Type
	Type GraphTypes
	// Title of the graph
	GraphTitle string
	// Labels for each bar
	Labels []string
	// Values must be a number type!
	Values []any
	// Colors for each bar
	Colors                 []string
	BackgroundColors       []string
	GraphBackgroundOpacity float64 // 0.0 - 1.0, background opacity of the bars/lines/dots/...
	GraphBorder            bool    // Border around the bars/lines/dots/...
	ShowResults            bool    // Show the results on the graph
}
