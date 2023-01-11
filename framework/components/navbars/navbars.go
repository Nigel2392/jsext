package navbars

import (
	"strconv"

	"github.com/Nigel2392/jsext/framework/components"
	"github.com/Nigel2392/jsext/framework/elements"
)

var prefix = "jsext-navbar-"

type Logo struct {
	ResourceURL string
	Alt         string
	URL         string
}

// Simple reusable navbar css
func getCSS(prefix, bg, fg string) string {
	return `.` + prefix + `logo {
			height: 50px;
			width: auto;
		}
		.` + prefix + `url {
			display: inline-block;
			padding: 15px 10px;
			text-decoration: none;
			font-weight: bold;
			color: ` + fg + `;
			height: 100%;
			transition: 0.3s;
			font-size: 20px;
		}
		.` + prefix + `url:hover {
			color: ` + bg + `;
			background-color: ` + fg + `;
		}
		.` + prefix + `url.active {
			color: ` + bg + `;
			background-color: ` + fg + `;
		}
	}`
}

// Simple navbar base.
func plainNav(foreground, background string, middle bool) []*elements.Element {
	var navbarMain = elements.Div().AttrClass(prefix + "main")
	var navbarLeft = navbarMain.Div().AttrClass(prefix + "left")
	var navbarMiddle *elements.Element
	var navbarRight = navbarMain.Div().AttrClass(prefix + "right")
	// Create middle element if specified.
	// Set the grid-template-columns to 2 or 3 depending on if middle is specified.
	var repeat = 2
	var templateAreas = `"left right"`
	if middle {
		navbarMiddle = navbarMain.Div().AttrClass(prefix + "middle")
		repeat = 3
		templateAreas = `"left middle right"`
	}
	navbarMain.StyleBlock(
		`.` + prefix + `main {
			width: 100%;
			height: 50px;
			background-color: ` + background + `;
			color: ` + foreground + `;
			overflow: hidden;
			display: grid;
			grid-template-columns: repeat(` + strconv.Itoa(repeat) + `, 1fr);
			grid-template-areas: ` + templateAreas + `;
			grid-template-rows: 1fr;
			column-gap: 0px;
		}
		.` + prefix + `left { grid-area: left; margin-left: 1%; }
		.` + prefix + `middle { grid-area: middle; }
		.` + prefix + `right { grid-area: right; text-align: right; margin-right: 1%; }`)

	// Also return middle element if specified.
	if middle {
		return []*elements.Element{navbarMain, navbarLeft, navbarMiddle, navbarRight}
	}
	return []*elements.Element{navbarMain, navbarLeft, navbarRight}
}

// Variables for styling the Official navbar's colors
var OfficialForeground = "#ffffff"
var OfficialBackground = "#142836"

// The navbar element.
func Official(logo *Logo, urls *elements.URLs) *elements.Element {
	var items = plainNav(OfficialForeground, OfficialBackground, false)
	var navbarMain, navbarLeft, navbarRight = items[0], items[1], items[2]

	navbarMain.StyleBlock(getCSS(prefix, OfficialBackground, OfficialForeground))

	navbarLeft.A(logo.URL).AttrClass(prefix + "logo-url").Append(
		elements.Img(logo.ResourceURL).AttrAlt(logo.Alt).AttrClass(prefix + "logo"),
	)

	urls.ForEach(func(k string, elem *elements.Element) {
		navbarRight.Append(elem.AttrClass(prefix + "url"))
	})

	return navbarMain
}

var SearchForeground = "#333333"
var SearchBackground = "#ffffff"
var SearchText = "Search"

// Returns the main navbar element, and the searchbar element, with the submit button in an array.
func Search(logo *Logo, urls *elements.URLs) (*elements.Element, []*elements.Element) {
	var items = components.SearchBar("search-", "#333333", "#ffffff", "Search")
	var searchBarContainer = items[0]
	var navbar = Custom(searchBarContainer, logo, urls, SearchBackground, SearchForeground)
	return navbar, items[1:]
}

func Custom(middle *elements.Element, logo *Logo, urls *elements.URLs, bg, fg string) *elements.Element {
	var items = plainNav(fg, bg, true)
	var navbarMain, navbarLeft, navbarMiddle, navbarRight = items[0], items[1], items[2], items[3]

	navbarMain.StyleBlock(getCSS(prefix, bg, fg))

	navbarLeft.A(logo.URL).AttrClass(prefix + "logo-url").Append(
		elements.Img(logo.ResourceURL).AttrAlt(logo.Alt).AttrClass(prefix + "logo"),
	)

	navbarMiddle.Append(middle)

	urls.ForEach(func(k string, elem *elements.Element) {
		navbarRight.Append(elem.AttrClass(prefix + "url"))
	})

	return navbarMain
}
