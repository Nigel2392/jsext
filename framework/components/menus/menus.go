//go:build js && wasm
// +build js,wasm

package menus

import (
	"strconv"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
)

type CSSField int

const (
	OpenBtnColor           CSSField = 1
	OpenBtnBg              CSSField = 2
	CloseBtnColor          CSSField = 3
	CloseBtnBg             CSSField = 4
	BtnSize                CSSField = 5
	OverlayBackgroundColor CSSField = 6
	TransitionDuration     CSSField = 7
	MenuContainerCSSBlock  CSSField = 8
	MenuCssBlock           CSSField = 9
	MenuItemCSSBlock       CSSField = 10
)

type MenuOptions struct {
	URLs        *elements.URLs
	cssMap      map[CSSField]string
	ClassPrefix string
	ElementTag  string
	URLFunc     func(e *elements.Element) *elements.Element
	CssFunc     func(containerClass string, manuContainerClass string, menuClass string, menuItemClass string) string
}

func NewMenuOptions() *MenuOptions {
	return &MenuOptions{
		URLs:        elements.NewURLs(),
		cssMap:      make(map[CSSField]string),
		ClassPrefix: "jsext",
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
	if m.cssMap[BtnSize] == "" {
		m.cssMap[BtnSize] = "50px"
	}
	if m.cssMap[MenuContainerCSSBlock] == "" {
		m.cssMap[MenuContainerCSSBlock] = `flex-direction: row; align-items: center; justify-content: center;`
	}
	if m.cssMap[MenuCssBlock] == "" {
		m.cssMap[MenuCssBlock] = `display:flex;flex-direction: row;`
	}
	if m.cssMap[MenuItemCSSBlock] == "" {
		m.cssMap[MenuItemCSSBlock] = `padding: 5px 10px; background-color: rgba(0,0,0,0.5); color: white; margin: 5px;`
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

	var openBtn = container.Div("☰").AttrClass(buttonClassOpen)
	var closeBtn = menu_container.Div().AttrClass(buttonClassClose)
	openBtn.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		var style = menu_container.JSExtElement().Style()
		style.Set("transform", "translateX(0)")
		openBtn.JSExtElement().Style().Opacity("0")
		openBtn.JSExtElement().Style().Set("pointerEvents", "none")
	})
	closeBtn.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		var style = menu_container.JSExtElement().Style()
		style.Set("transform", "translateX(-100%)")
		openBtn.JSExtElement().Style().Opacity("1")
		openBtn.JSExtElement().Style().Set("pointerEvents", "all")
	})

	menu.Children = make([]*elements.Element, 0, options.URLs.Len())
	options.URLs.ForEach(func(k string, elem *elements.Element) {
		elem.AttrClass(menuItemClass)
		menu.Append(options.URLFunc(elem))
	})
	menu_container.Append(menu)

	var css = `
	.` + containerClass + ` .` + manuContainerClass + ` {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background-color: ` + options.cssMap[OverlayBackgroundColor] + `;
		transition: transform ` + options.cssMap[TransitionDuration] + `;
		transform: translateX(-100%);
		display: flex;
		` + options.cssMap[MenuContainerCSSBlock] + `
	}
	.` + containerClass + ` .` + manuContainerClass + ` .` + menuClass + ` {
		list-style: none;
		margin: 0;
		padding: 0;
		` + options.cssMap[MenuCssBlock] + `
	}
	.` + containerClass + ` .` + manuContainerClass + ` .` + menuClass + ` .` + menuItemClass + ` {
		` + options.cssMap[MenuItemCSSBlock] + `
	}
	.` + buttonClassOpen + ` {
		position: fixed;
		top: 10px;
		left: 10px;
		height: ` + options.cssMap[BtnSize] + `;
		width: ` + options.cssMap[BtnSize] + `;
		font-size: calc(` + options.cssMap[BtnSize] + ` / 1.5);
		line-height: ` + options.cssMap[BtnSize] + `;
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
		top: 10px;
		left: 10px;
		height: ` + options.cssMap[BtnSize] + `;
		width: ` + options.cssMap[BtnSize] + `;
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
		height: calc(` + options.cssMap[BtnSize] + ` / 1.5);
		width: calc(` + options.cssMap[BtnSize] + ` / 10);
		background-color: ` + options.cssMap[CloseBtnColor] + `
	}
	.` + buttonClassClose + `:before {
		content: "";
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%) rotate(-45deg);
		height: calc(` + options.cssMap[BtnSize] + ` / 1.5);
		width: calc(` + options.cssMap[BtnSize] + ` / 10);
		background-color: ` + options.cssMap[CloseBtnColor] + `
	}
	`

	if options.CssFunc != nil {
		css += options.CssFunc(containerClass, manuContainerClass, menuClass, menuItemClass)
	}

	container.StyleBlock(css)

	return container
}

func Curtains(menu *MenuOptions) *elements.Element {
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
		var btnWidth = "200px"
		var btnHeight = "80px"
		var curtainColor = "#000"
		var complementaryColor = "#fff"
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
			color: ` + complementaryColor + `;
			text-decoration: none;
			text-transform: uppercase;
			position: relative;
			transition: all 0.4s;
			border-top: 1px solid ` + complementaryColor + `;
			border-bottom: 1px solid ` + complementaryColor + `;
			letter-spacing: calc(` + btnWidth + ` * 2 / 120);
    		font-weight: 800;
			margin: 0 !important;
			padding: calc(` + btnWidth + ` / ` + paddingDivisor + `) calc(` + btnWidth + ` / ` + paddingDivisor + `) calc(` + btnWidth + ` / ` + paddingDivisor + `) calc(` + btnWidth + ` / ` + paddingDivisor + `) !important;
		}
		.` + menuClass + ` li:first-child a {
			border-left: 1px solid ` + complementaryColor + `;
		}
		.` + menuClass + ` li:last-child a {
			border-right: 1px solid ` + complementaryColor + `;
		}
		.` + menuClass + ` li a:hover {
			color: ` + curtainColor + `;
		}
		.` + menuClass + ` li a:hover span {
			transform: scaleY(1);
		}
		.` + menuClass + ` li span {
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
