//go:build js && wasm
// +build js,wasm

package menus

import (
	"strconv"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
)

type OverlayDirection string

const (
	Left        OverlayDirection = "left"
	Right       OverlayDirection = "right"
	Top         OverlayDirection = "top"
	TopLeft     OverlayDirection = "topleft"
	TopRight    OverlayDirection = "topright"
	Bottom      OverlayDirection = "bottom"
	BottomLeft  OverlayDirection = "bottomleft"
	BottomRight OverlayDirection = "bottomright"
)

type CSSField int

const (
	OpenBtnColor           CSSField = 1
	OpenBtnBg              CSSField = 2
	CloseBtnColor          CSSField = 3
	CloseBtnBg             CSSField = 4
	ControlBtnSize         CSSField = 5
	OverlayBackgroundColor CSSField = 6
	TransitionDuration     CSSField = 7
	MenuContainerCSSBlock  CSSField = 8
	MenuCssBlock           CSSField = 9
	MenuItemCSSBlock       CSSField = 10
	TextColor              CSSField = 11
	TextColorActive        CSSField = 12
	BackgroundColor        CSSField = 13
	BackgroundActive       CSSField = 14
	ButtonWidth            CSSField = 15
)

type MenuOptions struct {
	URLs             *elements.URLs
	cssMap           map[CSSField]string
	ClassPrefix      string
	ElementTag       string
	URLFunc          func(e *elements.Element) *elements.Element
	CssFunc          func(containerClass string, manuContainerClass string, menuClass string, menuItemClass string) string
	OverlayDirection OverlayDirection
}

func NewMenuOptions(overlay OverlayDirection) *MenuOptions {
	return &MenuOptions{
		URLs:             elements.NewURLs(),
		cssMap:           make(map[CSSField]string),
		ClassPrefix:      "jsext",
		OverlayDirection: overlay,
	}
}

func (m *MenuOptions) setDefaults() {
	if m.ClassPrefix == "" {
		m.ClassPrefix = "jsext"
	}
	if m.ElementTag == "" {
		m.ElementTag = "ul"
	}
	if m.URLFunc == nil {
		m.URLFunc = func(e *elements.Element) *elements.Element {
			var li = elements.Li()
			li.Append(e)
			return li
		}
	}
	if m.cssMap[OpenBtnColor] == "" {
		m.cssMap[OpenBtnColor] = "white"
	}
	if m.cssMap[OverlayBackgroundColor] == "" {
		m.cssMap[OverlayBackgroundColor] = "rgba(0,0,0,0.5)"
	}
	if m.cssMap[TransitionDuration] == "" {
		m.cssMap[TransitionDuration] = "0.5s"
	}
	if m.cssMap[OpenBtnBg] == "" {
		m.cssMap[OpenBtnBg] = "rgba(0,0,0,0.5)"
	}
	if m.cssMap[CloseBtnBg] == "" {
		m.cssMap[CloseBtnBg] = "red"
	}
	if m.cssMap[CloseBtnColor] == "" {
		m.cssMap[CloseBtnColor] = "white"
	}
	if m.cssMap[ControlBtnSize] == "" {
		m.cssMap[ControlBtnSize] = "50px"
	}
	if m.cssMap[MenuContainerCSSBlock] == "" {
		m.cssMap[MenuContainerCSSBlock] = `flex-direction: row; align-items: center; justify-content: center;`
	}
	if m.cssMap[MenuCssBlock] == "" {
		m.cssMap[MenuCssBlock] = `display:flex;flex-direction: row;`
	}
	if m.cssMap[MenuItemCSSBlock] == "" {
		m.cssMap[MenuItemCSSBlock] = ``
	}
	if m.cssMap[TextColor] == "" {
		m.cssMap[TextColor] = "white"
	}
	if m.cssMap[TextColorActive] == "" {
		m.cssMap[TextColorActive] = "white"
	}
	if m.cssMap[BackgroundColor] == "" {
		m.cssMap[BackgroundColor] = "rgba(0,0,0,0.5)"
	}
	if m.cssMap[BackgroundActive] == "" {
		m.cssMap[BackgroundActive] = "#9200ff"
	}
	if m.cssMap[ButtonWidth] == "" {
		m.cssMap[ButtonWidth] = "200px"
	}

}

func Unstyled(options *MenuOptions) *elements.Element {
	options.setDefaults()

	if options.URLFunc == nil {
		panic("URLFunc must be set")
	}

	var (
		containerClass     = options.ClassPrefix + "-overlay-container"
		menuClass          = options.ClassPrefix + "-overlay-menu"
		manuContainerClass = options.ClassPrefix + "-overlay-menu-container"
		buttonClassOpen    = options.ClassPrefix + "-open-btn"
		buttonClassClose   = options.ClassPrefix + "-close-btn"
		menuItemClass      = options.ClassPrefix + "-menu-item"
	)

	var container = elements.Div().AttrClass(containerClass)
	var menu_container = container.Div().AttrClass(manuContainerClass)
	var menu = elements.NewElement(options.ElementTag).AttrClass(menuClass)

	var translateStart = "translateX(-100%)"
	var translateEnd = "translateX(0)"
	var BorderRadiusStart = "0 50% 50% 0"
	var BorderRadiusEnd = "0 0 0 0"
	var ButtonCss = "top: 10px; left: 10px;"
	switch options.OverlayDirection {
	case Left:
		translateStart = "translateX(-100%)"
		translateEnd = "translateX(0)"
		ButtonCss = "top: 10px; left: 10px;"
		BorderRadiusStart = "0 50% 50% 0"
	case Right:
		translateStart = "translateX(100%)"
		translateEnd = "translateX(0)"
		ButtonCss = "top: 10px; right: 10px;"
		BorderRadiusStart = "50% 0 0 50%"
	case Top:
		translateStart = "translateY(-100%)"
		translateEnd = "translateY(0)"
		ButtonCss = "top: 10px; left: 10px;"
		BorderRadiusStart = "0 0 50% 50%"
	case Bottom:
		translateStart = "translateY(100%)"
		translateEnd = "translateY(0)"
		ButtonCss = "bottom: 10px; left: 10px;"
		BorderRadiusStart = "50% 50% 0 0"
	case TopLeft:
		translateStart = "translate(-100%, -100%)"
		translateEnd = "translate(0, 0)"
		ButtonCss = "top: 10px; left: 10px;"
		BorderRadiusStart = "0 0 50% 0"
	case TopRight:
		translateStart = "translate(100%, -100%)"
		translateEnd = "translate(0, 0)"
		ButtonCss = "top: 10px; right: 10px;"
		BorderRadiusStart = "0 0 0 50%"
	case BottomLeft:
		translateStart = "translate(-100%, 100%)"
		translateEnd = "translate(0, 0)"
		ButtonCss = "bottom: 10px; left: 10px;"
		BorderRadiusStart = "0 50% 0 0"
	case BottomRight:
		translateStart = "translate(100%, 100%)"
		translateEnd = "translate(0, 0)"
		ButtonCss = "bottom: 10px; right: 10px;"
		BorderRadiusStart = "50% 0 0 0"
	}

	var openBtn = container.Div("☰").AttrClass(buttonClassOpen)
	var closeBtn = menu_container.Div().AttrClass(buttonClassClose)
	openBtn.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		var style = menu_container.JSExtElement().Style()
		style.Set("transform", translateEnd)
		openBtn.JSExtElement().Style().Opacity("0")
		menu_container.JSExtElement().Style().Set("borderRadius", BorderRadiusEnd)
		openBtn.JSExtElement().Style().Set("pointerEvents", "none")
	})
	closeBtn.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		var style = menu_container.JSExtElement().Style()
		style.Set("transform", translateStart)
		openBtn.JSExtElement().Style().Opacity("1")
		menu_container.JSExtElement().Style().Set("borderRadius", BorderRadiusStart)
		openBtn.JSExtElement().Style().Set("pointerEvents", "all")
	})

	menu.Children = make([]*elements.Element, 0, options.URLs.Len())
	options.URLs.ForEach(func(k string, elem *elements.Element) {
		elem.AttrClass(menuItemClass)
		menu.Append(options.URLFunc(elem))
	})
	menu_container.Append(menu)

	var css = `
	.` + manuContainerClass + ` {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: ` + options.cssMap[OverlayBackgroundColor] + `;
		transition: transform ` + options.cssMap[TransitionDuration] + `, border-radius ` + options.cssMap[TransitionDuration] + `;
		border-radius: ` + BorderRadiusStart + `;
		transform: ` + translateStart + `;
		display: flex;
		` + options.cssMap[MenuContainerCSSBlock] + `
	}
	.` + menuClass + ` {
		list-style: none;
		margin: 0;
		padding: 0;
		` + options.cssMap[MenuCssBlock] + `
	}
	.` + menuItemClass + ` {
		` + options.cssMap[MenuItemCSSBlock] + `
	}
	.` + buttonClassOpen + ` {
		position: fixed;
		` + ButtonCss + `
		height: ` + options.cssMap[ControlBtnSize] + `;
		width: ` + options.cssMap[ControlBtnSize] + `;
		font-size: calc(` + options.cssMap[ControlBtnSize] + ` / 1.5);
		line-height: ` + options.cssMap[ControlBtnSize] + `;
		text-align: center;
		background-color: ` + options.cssMap[OpenBtnBg] + `;
		color: ` + options.cssMap[OpenBtnColor] + `;
		border-radius: 5px;
		cursor: pointer;
		transition: all ` + options.cssMap[TransitionDuration] + `;
		z-index: 999;
	}
	.` + buttonClassClose + ` {
		position: absolute;
		` + ButtonCss + `
		height: ` + options.cssMap[ControlBtnSize] + `;
		width: ` + options.cssMap[ControlBtnSize] + `;
		background-color: ` + options.cssMap[CloseBtnBg] + `;
		border-radius: 5px;
		cursor: pointer;
	}
	.` + buttonClassClose + `:after {
		content: "";
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%) rotate(45deg);
		height: calc(` + options.cssMap[ControlBtnSize] + ` / 1.5);
		width: calc(` + options.cssMap[ControlBtnSize] + ` / 10);
		background-color: ` + options.cssMap[CloseBtnColor] + `
	}
	.` + buttonClassClose + `:before {
		content: "";
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%) rotate(-45deg);
		height: calc(` + options.cssMap[ControlBtnSize] + ` / 1.5);
		width: calc(` + options.cssMap[ControlBtnSize] + ` / 10);
		background-color: ` + options.cssMap[CloseBtnColor] + `
	}
	`
	if options.CssFunc != nil {
		css += options.CssFunc(containerClass, manuContainerClass, menuClass, menuItemClass)
	}

	container.StyleBlock(css)

	return container
}

func Blurry(menu *MenuOptions) *elements.Element {
	menu.setDefaults()
	menu.URLFunc = func(e *elements.Element) *elements.Element {
		var li = elements.Li()
		li.Append(e)
		return li
	}
	var btnWidth = menu.cssMap[ButtonWidth]
	var btnHeight = "calc(" + btnWidth + " / 2.5)"
	menu.cssMap[MenuCssBlock] = `display: flex;`
	menu.CssFunc = func(containerClass, manuContainerClass, menuClass, menuItemClass string) string {
		return `.` + menuClass + ` li {
			list-style: none;
			margin: 0 20px;
			transition: 0.5s;
		}
		.` + menuClass + ` li a {
			display: block;
			position: relative;
			text-decoration: none;
			font-size: calc(` + btnHeight + ` / 2);
			line-height: ` + btnHeight + `;
			text-align: center;
			color: ` + menu.cssMap[TextColor] + `;
			text-transform: uppercase;
			transition: all 0.3s;
			z-index: 1;
			width: ` + btnWidth + `;
			height: ` + btnHeight + `;
		}
		.` + menuClass + ` li a:hover {
			transform: scale(1.5);
			filter: blur(0) !important;
			opacity: 1 !important;
		}
		.` + menuClass + ` li a:before {
			content: "";
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background: ` + menu.cssMap[BackgroundActive] + `;
			transition: 0.3s;
			transform: scale(0);
			z-index: -1;
		}
		.` + menuClass + ` li a:hover:before {
			transform: scale(1);
		}
		.` + menuClass + `:hover li a {
			filter: blur(calc(` + btnWidth + ` / 40));
			opacity: 0.2;
		}
		`
	}

	var m = Unstyled(menu)
	return m
}

func Curtains(menu *MenuOptions) *elements.Element {
	menu.setDefaults()
	menu.URLFunc = func(e *elements.Element) *elements.Element {
		var li = elements.Li()
		li.Append(e)
		e.Span()
		e.Span()
		e.Span()
		e.Span()
		return li
	}
	menu.CssFunc = func(containerClass, manuContainerClass, menuClass, menuItemClass string) string {
		var btnWidth = menu.cssMap[ButtonWidth]
		var btnHeight = "calc(" + btnWidth + " / 2.5)"
		var backgroundColor = menu.cssMap[BackgroundColor]
		var backgroundColorActive = menu.cssMap[BackgroundActive]
		var textColor = menu.cssMap[TextColor]
		var textColorActive = menu.cssMap[TextColorActive]
		var padding = 50
		for i := 0; i < menu.URLs.Len(); i++ {
			if padding == 0 {
				break
			}
			padding -= 5
		}
		var paddingDivisor = strconv.Itoa(padding)

		return `.` + menuItemClass + ` {
			display: block;
			width: ` + btnWidth + `;
			height: ` + btnHeight + `;
			line-height: ` + btnHeight + `;
			font-size: calc(` + btnHeight + ` / 2);
			text-align: center;
			color: ` + textColor + `;
			text-decoration: none;
			text-transform: uppercase;
			position: relative;
			transition: all 0.4s;
			border-top: 1px solid ` + textColor + `;
			border-bottom: 1px solid ` + textColor + `;
			letter-spacing: calc(` + btnWidth + ` * 2 / 120);
    		font-weight: 800;
			margin: 0 !important;
			padding: calc(` + btnWidth + ` / ` + paddingDivisor + `) calc(` + btnWidth + ` / ` + paddingDivisor + `) calc(` + btnWidth + ` / ` + paddingDivisor + `) calc(` + btnWidth + ` / ` + paddingDivisor + `) !important;
			background: ` + backgroundColor + ` ;
			z-index: 1;
		}
		.` + menuClass + ` li:first-child a {
			border-left: 1px solid ` + textColor + `;
		}
		.` + menuClass + ` li:last-child a {
			border-right: 1px solid ` + textColor + `;
		}
		.` + menuClass + ` li a:hover {
			color: ` + textColorActive + `;
		}
		.` + menuClass + ` li a:hover span {
			transform: scaleY(1);
		}
		.` + menuClass + ` li span {
			background: ` + backgroundColorActive + `;
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
		.` + menuClass + ` li span:nth-child(1) {
			transition-delay: 0.2s;
		}
		.` + menuClass + ` li span:nth-child(2) {
			left: 25%;
			transition-delay: 0.1s;
		}
		.` + menuClass + ` li span:nth-child(3) {
			left: 50%;
		}
		.` + menuClass + ` li span:nth-child(4) {
			left: 75%;
			transition-delay: 0.3s;
		}
	`
	}
	var m = Unstyled(menu)
	return m
}
