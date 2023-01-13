//go:build js && wasm
// +build js,wasm

package misc

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
	"github.com/Nigel2392/jsext/framework/helpers"
	"github.com/Nigel2392/jsext/framework/helpers/csshelpers"
)

func Join(css ...[]string) []string {
	var count int
	for _, v := range css {
		count += len(v)
	}
	var ret = make([]string, 0, count)
	for _, v := range css {
		ret = append(ret, v...)
	}
	return ret
}

var Container = []string{
	"width: 80%",
	"margin: 1rem auto",
	"padding: 0 2%",
	"padding-bottom: 1rem",
}

var Box = []string{
	"border-radius: 0 0 0.5rem 0.5rem",
	"box-shadow: 0 1rem 1rem 0 rgba(0,0,0,0.2)",
}

// Search bar element with button
// Returns slice of elements
// [0]: search container, [1]: search bar, [2]: search bar submit button
func SearchBar(classPrefix, foregroundHex, background, text string) []*elements.Element {
	classPrefix = classPrefix + helpers.FNVHashString(classPrefix+background+foregroundHex+text) + "-"
	var searchContainer = elements.Div().AttrClass(classPrefix + "search-container")
	var searchbar = searchContainer.Input("text", "search", text).AttrClass(classPrefix + "searchbar")
	var searchBarSubmit = searchContainer.Button(text).AttrClass(classPrefix + "searchbar-submit")
	var borderColor, err = csshelpers.Hex(foregroundHex)
	if err != nil {
		panic(err)
	}

	var b_A, b_G, b_B = borderColor.RGB255()
	var b_A_str = strconv.Itoa(int(b_A))
	var b_G_str = strconv.Itoa(int(b_G))
	var b_B_str = strconv.Itoa(int(b_B))

	jsext.StyleBlock(classPrefix, `
		.`+classPrefix+`search-container {
			display: grid;
			grid-template-columns: 3fr 1fr;
			grid-template-areas: "searchbar submit";
			grid-template-rows: 1fr;
			column-gap: 3px;
		}
		.`+classPrefix+`searchbar {
			height: 35px;
			margin: 6px 0;
			padding: 0 5px;
			background-color: `+background+`;
			color: `+foregroundHex+`;
			border: 1px solid rgba(`+b_A_str+`, `+b_G_str+`, `+b_B_str+`, 0.5);
			border-radius: 5px;
			font-size: 20px;
		}
		.`+classPrefix+`searchbar:focus {
			outline: none;
		}
		.`+classPrefix+`searchbar-submit {
			grid-area: submit;
			width: 100%;
			height: 37px;
			margin: 6px 0;
			padding: 0 5px;
			background-color: `+background+`;
			color: `+foregroundHex+`;
			border: 1px solid rgba(`+b_A_str+`, `+b_G_str+`, `+b_B_str+`, 0.5);
			border-radius: 5px;
			cursor: pointer;
			font-size: 20px;
		}
	`)
	return []*elements.Element{searchContainer, searchbar, searchBarSubmit}
}

var spacing = " "
var jsextPrefix = "jsext-"

// Provide a grid based on a pattern.
// Example: "$$$ ## $214" will create a grid with 3 columns.
// Column 1 will be 3fr, column 2 will be 2fr, column 3 will be 4fr.
func Grid(gridPattern string) (*elements.Element, []*elements.Element, error) {
	var className = jsextPrefix + helpers.FNVHashString(gridPattern)
	var grid = elements.Div().AttrClass(className + "-grid")
	var splitGrid = strings.Split(gridPattern, spacing)
	var fractionSlice = make([]string, len(splitGrid))
	if len(splitGrid) == 0 {
		return nil, nil, errors.New("Grid pattern is empty")
	}
	for i, v := range splitGrid {
		fractionSlice[i] = strconv.Itoa(len(v)) + "fr"
	}
	var gridItems = make([]*elements.Element, len(splitGrid))
	for i := 0; i < len(splitGrid); i++ {
		var gridItem = grid.Div().AttrClass(className + "-grid-item")
		gridItems[i] = gridItem
	}
	var fractionString = strings.Join(fractionSlice, " ")
	var css = `
		.` + className + `-grid {
			display: grid;
			grid-template-columns: ` + fractionString + `;
			grid-template-rows: auto;
			grid-gap: 1%;
			width: 100%;
		}
		.` + className + `-grid-item {
			width: 100%;
		}
		`
	jsext.StyleBlock(className, css)
	return grid, gridItems, nil
}
