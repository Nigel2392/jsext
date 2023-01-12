package navbars

import (
	"strconv"
	"strings"

	"github.com/Nigel2392/jsext/framework/components"
	"github.com/Nigel2392/jsext/framework/elements"
)

var prefix = "jsext-navbar-"

type Logo struct {
	ResourceURL string
	Alt         string
	URL         string
}

func NewLogo(resourceURL, url string) *Logo {
	var alt = "Logo"
	if resourceURL != "" && strings.Contains(resourceURL, "/") && strings.Contains(resourceURL, ".") {
		alt = resourceURL[strings.LastIndex(resourceURL, "/")+1:]
		alt = alt[:strings.LastIndex(alt, ".")]
	}
	return &Logo{
		ResourceURL: resourceURL,
		Alt:         alt,
		URL:         url,
	}
}

// Simple reusable navbar css
func getCSS(prefix, bg, fg string) string {
	return `.` + prefix + `logo {
			height: 50px;
			width: auto;
		}
		.` + prefix + `url {
			display: inline-block;
			padding: 0 10px;
			text-decoration: none;
			font-weight: bold;
			color: ` + fg + `;
			height: 50px;
			line-height: 50px;
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
func plainNav(foreground, background string, middle int) []*elements.Element {
	var navbarMain = elements.Div().AttrClass(prefix + "main")
	var navbarLeft = navbarMain.Div().AttrClass(prefix + "left")
	var navbarRight = navbarMain.Div().AttrClass(prefix + "right")
	var navItemSlice = make([]*elements.Element, 0, 3+middle)
	navItemSlice = append(navItemSlice, navbarMain, navbarLeft)

	// Create middle element if specified.
	// Set the grid-template-columns to 2 or 3 depending on if middle is specified.
	var templateAreas = `"left right"`
	var middleCSS = strings.Builder{}
	var columnBuilder = strings.Builder{}
	columnBuilder.WriteString("1.5fr")
	if middle > 0 {
		templateAreas = `"left`
		for i := 0; i < middle; i++ {
			templateAreas += ` middle-` + strconv.Itoa(i) + ``
			navItemSlice = append(navItemSlice, navbarMain.Div().AttrClass(prefix+"middle-"+strconv.Itoa(i)))
			middleCSS.WriteString(`.` + prefix + `middle-` + strconv.Itoa(i) + ` { grid-area: middle-` + strconv.Itoa(i) + `; height: 50px; }`)
			columnBuilder.WriteString(" 1fr")
		}
		templateAreas += ` right"`
	}
	columnBuilder.WriteString(" 1.5fr")
	middleCSS.WriteString(
		`.` + prefix + `main {
			width: 100%;
			height: 50px;
			background-color: ` + background + `;
			color: ` + foreground + `;
			display: grid;
			grid-template-columns: ` + columnBuilder.String() + `;
			grid-template-areas: ` + templateAreas + `;
			grid-template-rows: 1fr;
			column-gap: 3px;
		}
		.` + prefix + `left { grid-area: left; margin-left: 1%; height: 50px; }
		.` + prefix + `right { grid-area: right; text-align: right; margin-right: 1%; height: 50px; }`)

	navbarMain.StyleBlock(middleCSS.String())

	navItemSlice = append(navItemSlice, navbarRight)

	// Also return middle element if specified.
	return navItemSlice
}

// Variables for styling the Official navbar's colors
var OfficialForeground = "#ffffff"
var OfficialBackground = "#142836"

// The navbar element.
func Official(logo *Logo, urls *elements.URLs) *elements.Element {
	var items = plainNav(OfficialForeground, OfficialBackground, 0)
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
	var navbar = Custom(logo, urls, SearchBackground, SearchForeground, searchBarContainer)
	return navbar, items[1:]
}

func Custom(logo *Logo, urls *elements.URLs, bg, fg string, middle ...*elements.Element) *elements.Element {
	var items = plainNav(fg, bg, len(middle))

	items[0].StyleBlock(getCSS(prefix, bg, fg))

	items[1].A(logo.URL).AttrClass(prefix + "logo-url").Append(
		elements.Img(logo.ResourceURL).AttrAlt(logo.Alt).AttrClass(prefix + "logo"),
	)

	for i := 0; i < len(middle); i++ {
		items[i+2].Append(middle[i])
	}

	urls.ForEach(func(k string, elem *elements.Element) {
		items[len(items)-1].Append(elem.AttrClass(prefix + "url"))
	})

	return items[0]
}
