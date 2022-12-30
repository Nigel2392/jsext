//go:build js && wasm
// +build js,wasm

package menus

import (
	"strconv"
	"strings"

	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/elements"
	"github.com/Nigel2392/jsext/helpers/csshelpers"
)

// Provide a fullscreen menu with a curtain drop effect.
// Style is included in the returned element.
// CLASSES AND IDS
//
//	".jsext-menu-container"
//	"#jsext-menu-open-btn"
//	"#jsext-menu"
//	"#jsext-menu-close-btn"
type MenuCurtainDrop struct {
	BtnWidth     int
	CurtainColor string
	CompColor    string
}

func (m *MenuCurtainDrop) Nav(urls *components.URLs) *elements.Element {
	return m.Menu(urls, m.BtnWidth, m.CurtainColor, m.CompColor)
}

func (m *MenuCurtainDrop) Menu(urls *components.URLs, btnWidth int, curtainColor string, compColor ...string) *elements.Element {
	// Menu container
	var menu_container = elements.Div().AttrClass("jsext-menu-container")
	// Open button
	menu_container.Div().InnerHTML("&#9776;").AttrID("jsext-menu-open-btn").Set("onClick", "document.getElementById('jsext-menu').style.transform = 'translateX(0)';document.getElementById('jsext-menu-open-btn').style.display = 'none'")
	// Menu
	var menu = menu_container.Div().AttrID("jsext-menu")
	// Close button
	var closeBtn = menu.Div().AttrID("jsext-menu-close-btn").Set("onClick", "document.getElementById('jsext-menu').style.transform = 'translateX(-100%)';document.getElementById('jsext-menu-open-btn').style.display = 'block'")
	closeBtn.Span()
	closeBtn.Span()
	// Urls for the navigation menu
	if len(urls.URLs) > 0 {
		var ul = menu.Ul()
		for _, url := range urls.URLs {
			var urlItem *elements.Element
			if strings.HasPrefix(url.URL, "external:") {
				urlItem = ul.Li().A(strings.TrimPrefix(url.URL, "external:"), url.Name)
			} else {
				urlItem = ul.Li().A("router:"+url.URL, url.Name)
			}
			urlItem.AttrID("URL-" + url.Name)
			urlItem.Span()
			urlItem.Span()
			urlItem.Span()
			urlItem.Span()

			if url.Hidden {
				urlItem.AttrStyle("display: none")
			}
			url.Run()
			url.Elem = urlItem
		}
	}
	// Get the main color
	colr, err := csshelpers.Hex(curtainColor)
	if err != nil {
		panic(err)
	}
	// Get the complementary color
	var complementaryColor = colr.Complementary().Hex()
	var r, g, b = colr.RGB255()
	// Check if the color is valid, and not the same.
	if strings.EqualFold(complementaryColor, curtainColor) {
		if r == 0 && g == 0 && b == 0 {
			complementaryColor = "#FFFFFF"
		} else {
			if r > 127 && g > 127 && b > 127 {
				complementaryColor = "#000000"
			} else {
				complementaryColor = "#FFFFFF"
			}
		}
	}
	// If a complementary color is provided, use that instead.
	if len(compColor) > 0 {
		complementaryColor = compColor[0]
	}
	// Calculate the height of the menu buttons
	var btnHeight = btnWidth / 3

	// Convert the RGB values to strings for use in CSS rgba(r, g, b, a)
	var r_str, g_str, b_str = strconv.Itoa(int(r)), strconv.Itoa(int(g)), strconv.Itoa(int(b))

	// Add the style to the menu container
	menu_container.StyleBlock(`
		#jsext-menu-open-btn,
		#jsext-menu-close-btn{
			position: fixed;
			top: 10px;
			left: 10px;
			font-size: ` + strconv.Itoa(btnWidth/4) + `px;
			cursor: pointer;
		}
		#jsext-menu-open-btn {
			color: ` + curtainColor + `;
		}
		#jsext-menu-close-btn span {
			background-color: ` + complementaryColor + `;
		}
		#jsext-menu-close-btn {
			top: calc(10px + ` + strconv.Itoa(btnHeight/8) + `px);
			background-color:  ` + complementaryColor + `;
			border-radius: 50%;
			transition: all 0.3s;
			width: ` + strconv.Itoa(btnWidth/4) + `px;
			height: ` + strconv.Itoa(btnWidth/4) + `px;
		}
		#jsext-menu-close-btn span {
			position: absolute;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
			transition: all 0.3s;
		}
		#jsext-menu-close-btn span:nth-child(1),
		#jsext-menu-close-btn span:nth-child(2) {
			content: "";
			position: absolute;
			top: 50%;
			left: 50%;
			width: 60%;
			height: ` + strconv.Itoa(btnHeight/8) + `px;
			background-color: ` + curtainColor + `;
			z-index: 1;
			transition: all 0.3s;
		}
		#jsext-menu-close-btn span:nth-child(1) {
			transform: translate(-50%, -50%) rotate(45deg);
		}
		#jsext-menu-close-btn span:nth-child(2) {
			transform: translate(-50%, -50%) rotate(-45deg);
		}
		#jsext-menu-close-btn:hover {
			cursor: pointer;
			background-color: ` + curtainColor + `;
		}
		#jsext-menu-close-btn:hover span:nth-child(1),
		#jsext-menu-close-btn:hover span:nth-child(2) {
			background-color: ` + complementaryColor + `;
		}
		#jsext-menu {
			position: fixed;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background-color: rgba(` + r_str + `, ` + g_str + `, ` + b_str + `, 0.8);
			display: flex;
			justify-content: center;
			align-items: center;
			transform: translateX(-100%);
			transition: all 0.4s;
			z-index: 1000;
		}
		#jsext-menu ul {
			margin: 0;
			padding: 0;
			display: flex;
			list-style: none;
		}
		#jsext-menu ul li a {
			display: block;
			width: ` + strconv.Itoa(btnWidth) + `px;
			height: ` + strconv.Itoa(btnHeight) + `px;
			line-height: ` + strconv.Itoa(btnHeight) + `px;
			text-align: center;
			color: ` + complementaryColor + `;
			text-decoration: none;
			text-transform: uppercase;
			position: relative;
			transition: all 0.4s;
			border-top: 1px solid ` + complementaryColor + `;
			border-bottom: 1px solid ` + complementaryColor + `;
			letter-spacing: ` + strconv.Itoa(btnWidth*2/120) + `px;
    		font-weight: 800;
		}
		#jsext-menu ul li:first-child a {
			border-left: 1px solid ` + complementaryColor + `;
		}
		#jsext-menu ul li:last-child a {
			border-right: 1px solid ` + complementaryColor + `;
		}
		#jsext-menu ul li a:hover {
			color: ` + curtainColor + `;
		}
		#jsext-menu ul li a:hover span {
			transform: scaleY(1);
		}
		#jsext-menu ul li span {
			background-color: ` + complementaryColor + `;
			position: absolute;
			left: 0;
			top: 0;
			width: 25%;
			height: 100%;
			z-index: -1;
			transform: scaleY(0);
			transition: all 0.5s;
			transform-origin: top;
		}
		#jsext-menu ul li span:nth-child(1) {
			transition-delay: 0.2s;
		}
		#jsext-menu ul li span:nth-child(2) {
			left: 25%;
			transition-delay: 0.1s;
		}
		#jsext-menu ul li span:nth-child(3) {
			left: 50%;
		}
		#jsext-menu ul li span:nth-child(4) {
			left: 75%;
			transition-delay: 0.3s;
		}
		@media screen and (min-width: 1367px) {
			#jsext-menu ul li a {
				font-size: ` + strconv.Itoa(btnWidth/8) + `px;
				width: ` + strconv.Itoa(btnWidth) + `px;
				height: ` + strconv.Itoa(btnHeight) + `px;
				line-height: ` + strconv.Itoa(btnHeight) + `px;
				letter-spacing: ` + strconv.Itoa(btnWidth*2/120) + `px;
			}
		}
		@media screen and (max-width: 1366px) {
			#jsext-menu ul li a {
				font-size: ` + strconv.Itoa(btnWidth/14) + `px;
				width: ` + strconv.Itoa(int(float64(btnWidth)/1.5)) + `px;
				height: ` + strconv.Itoa(int(float64(btnHeight)/1.5)) + `px;
				line-height: ` + strconv.Itoa(int(float64(btnHeight)/1.5)) + `px;
				letter-spacing: ` + strconv.Itoa(btnWidth/60) + `px;
			}
		}
		@media screen and (max-width: 768px) {
			#jsext-menu ul li a {
				font-size: ` + strconv.Itoa(btnWidth/22) + `px;
				width: ` + strconv.Itoa(int(float64(btnWidth)/3)) + `px;
				height: ` + strconv.Itoa(int(float64(btnHeight)/3)) + `px;
				line-height: ` + strconv.Itoa(int(float64(btnHeight)/3)) + `px;
				letter-spacing: ` + strconv.Itoa(btnWidth/180) + `px;
			}
		}
		@media screen and (max-width: 380px) {
			#jsext-menu ul li a {
				font-size: ` + strconv.Itoa(btnWidth/60) + `px;
				width: ` + strconv.Itoa(btnWidth/12) + `px;
				height: ` + strconv.Itoa(btnHeight/12) + `px;
				line-height: ` + strconv.Itoa(btnHeight/12) + `px;
				letter-spacing: ` + strconv.Itoa(btnWidth/240) + `px;
			}
		}
	`)
	// Return the menu container and the urlItems
	return menu_container
}
