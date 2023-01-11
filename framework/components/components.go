//go:build js && wasm
// +build js,wasm

package components

import (
	"strconv"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
	"github.com/Nigel2392/jsext/framework/helpers/csshelpers"
)

// Component interface
type Component interface {
	Render() jsext.Element
}

// Component with eventlistener interface
type ComponentWithEventListener interface {
	Component
	AddEventListener(event string, fn func(this jsext.Value, event jsext.Event))
}

// Loader component
type Loader interface {
	Stop()        // Stop the loader.
	Show()        // Show the loader.
	Run(f func()) // Run the function, finalize loader automatically.
	Finalize()    // Finalize loader.
}

// URL for use in components
type URL struct {
	Name string
	Url  string
}

// Search bar element with button
// Returns slice of elements
// [0]: search container, [1]: search bar, [2]: search bar submit button
func SearchBar(classPrefix, foregroundHex, background, text string) []*elements.Element {
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

	searchContainer.StyleBlock(`
		.` + classPrefix + `search-container {
			display: grid;
			grid-template-columns: 3fr 1fr;
			grid-template-areas: "searchbar submit";
			grid-template-rows: 1fr;
			column-gap: 3px;
		}
		.` + classPrefix + `searchbar {
			height: 35px;
			margin: 6px 0;
			padding: 0 5px;
			background-color: ` + background + `;
			color: ` + foregroundHex + `;
			border: 1px solid rgba(` + b_A_str + `, ` + b_G_str + `, ` + b_B_str + `, 0.5);
			border-radius: 5px;
			font-size: 20px;
		}
		.` + classPrefix + `searchbar:focus {
			outline: none;
		}
		.` + classPrefix + `searchbar-submit {
			grid-area: submit;
			width: 100%;
			height: 37px;
			margin: 6px 0;
			padding: 0 5px;
			background-color: ` + background + `;
			color: ` + foregroundHex + `;
			border: 1px solid rgba(` + b_A_str + `, ` + b_G_str + `, ` + b_B_str + `, 0.5);
			border-radius: 5px;
			cursor: pointer;
			font-size: 20px;
		}
	`)
	return []*elements.Element{searchContainer, searchbar, searchBarSubmit}
}
