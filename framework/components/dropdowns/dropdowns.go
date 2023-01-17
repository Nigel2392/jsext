//go:build js && wasm
// +build js,wasm

package dropdowns

import (
	"strconv"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
)

type DropdownOptions struct {
	Background        string
	Color             string
	BorderWidth       string
	Width             string
	Height            string
	ButtonBorderWidth string
	ButtonWidth       string
	ButtonHeight      string
	ButtonText        string
	Header            *elements.Element
	Footer            *elements.Element
	MenuItems         []*elements.Element
	ItemsPerColumn    int
	Prefix            string
}

func (d *DropdownOptions) SetDefaults() {
	if d.Background == "" {
		d.Background = "#ffffff"
	}
	if d.Color == "" {
		d.Color = "#333333"
	}
	if d.BorderWidth == "" {
		d.BorderWidth = "0px"
	}
	if d.ButtonText == "" {
		d.ButtonText = "Menu"
	}
	if d.Width == "" {
		d.Width = "100%"
	}
	if d.Height == "" {
		d.Height = "400px"
	}
	if d.ButtonBorderWidth == "" {
		d.ButtonBorderWidth = "0px"
	}
	if d.ButtonWidth == "" {
		d.ButtonWidth = "100%"
	}
	if d.ButtonHeight == "" {
		d.ButtonHeight = "50px"
	}
	if d.ItemsPerColumn == 0 {
		d.ItemsPerColumn = 8
	}
	if d.Prefix == "" {
		d.Prefix = "jsext-dropdown-"
	}
	d.Width = "calc(" + d.Width + " - calc(" + d.BorderWidth + " * 2))"
}

func Dropdown(options DropdownOptions) *elements.Element {
	if len(options.MenuItems) == 0 {
		panic("DropDownMenu: No menu items provided")
	}

	options.SetDefaults()

	var dropDownContainer = elements.Div().AttrClass(options.Prefix + "container")
	var dropDownButton = dropDownContainer.Button(options.ButtonText).AttrClass(options.Prefix + "dropbtn")
	var dropDownContentContainer = dropDownContainer.Div().AttrClass(options.Prefix + "content-container").AttrStyle("display:none")
	var dropDownContent = dropDownContentContainer.Div().AttrClass(options.Prefix + "content")

	dropDownButton.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		if dropDownContentContainer.JSExtElement().Style().Display() == "block" {
			dropDownContentContainer.JSExtElement().Style().Display("none")
			dropDownButton.JSExtElement().ClassList().Remove(options.Prefix + "active")

		} else {
			dropDownContentContainer.JSExtElement().Style().Display("block")
			dropDownButton.JSExtElement().ClassList().Add(options.Prefix + "active")
		}
	})

	columns := len(options.MenuItems) / options.ItemsPerColumn
	if len(options.MenuItems)%options.ItemsPerColumn != 0 {
		columns++
	}

	for i := 0; i < columns; i++ {
		var column = dropDownContent.Div().AttrClass(options.Prefix + "column" + strconv.Itoa(i))
		for j := 0; j < options.ItemsPerColumn; j++ {
			var index = i*options.ItemsPerColumn + j
			if index >= len(options.MenuItems) {
				break
			}
			column.Append(options.MenuItems[index].AttrClass(options.Prefix + "item"))
		}
	}

	var css = `
		.` + options.Prefix + `container {
			position: relative;
			display: inline-block;
			width: ` + options.ButtonWidth + `;
		}
		.` + options.Prefix + `dropbtn {
			background-color: ` + options.Background + `;
			color: ` + options.Color + `;
			padding: calc(` + options.ButtonHeight + ` / 2), calc(` + options.ButtonWidth + ` / 2);
			font-size: calc(` + options.ButtonHeight + ` / 2);
			border: none;
			cursor: pointer;
			width: ` + options.ButtonWidth + `;
			height: ` + options.ButtonHeight + `;
			transition: 0.3s;
			border: ` + options.ButtonBorderWidth + ` solid ` + options.Color + `;
			font-weight: bold;
		}
		.` + options.Prefix + `dropbtn:hover {
			background-color: ` + options.Color + `;
			color: ` + options.Background + `;
		}
		.` + options.Prefix + `content-container {
			display: none;
			position: absolute;
			background-color: ` + options.Background + `;
			min-width: ` + options.Width + `;
			max-width: ` + options.Width + `;
			min-height: ` + options.Height + `;
			max-height: ` + options.Height + `;
			overflow: auto;
			z-index: 3;
			border: ` + options.BorderWidth + ` solid ` + options.Color + `;
		}
		.` + options.Prefix + `content {
			color: ` + options.Color + `;
			text-decoration: none;
			display: grid;
			grid-template-columns: repeat(` + strconv.Itoa(columns) + `, 1fr);
		}
		.` + options.Prefix + `item {
			background-color: ` + options.Background + `;
			color: ` + options.Color + `;
			font-size: 20px;
			text-decoration: none;
			display: block;
			text-align: center;
			height: calc(` + options.Height + ` / ` + strconv.Itoa(options.ItemsPerColumn) + `);
			line-height: calc(` + options.Height + `/ ` + strconv.Itoa(options.ItemsPerColumn) + `);
			transition: 0.3s;
		}
		.` + options.Prefix + `item:hover {
			background-color: ` + options.Color + `;
			color: ` + options.Background + `;
		}
		.` + options.Prefix + `active {
			background-color: ` + options.Color + `;
			color: ` + options.Background + `;
		}
		`

	for i := 0; i < columns; i++ {
		css += `
			.` + options.Prefix + `column` + strconv.Itoa(i) + ` {
				display:block;
			}
		`
	}

	dropDownContainer.StyleBlock(css)

	return dropDownContainer
}

func DropdownElement(opts DropdownOptions, absolutePosition bool) *elements.Element {

	if len(opts.MenuItems) == 0 {
		panic("No menu items provided")
	}
	if opts.Header == nil {
		panic("No header provided")
	}

	opts.SetDefaults()

	var container = elements.Div().AttrClass(opts.Prefix + "dropdown")
	var header = container.Div().AttrClass(opts.Prefix + "dropdown-header").Append(opts.Header)
	container.Append()
	var menu = container.Div().AttrClass(opts.Prefix + "dropdown-content")

	menu.Append(elements.Div().AttrClass(opts.Prefix + "dropdown-item").Append(opts.MenuItems[0]))

	if opts.Footer != nil {
		var footer_no_ptr = *opts.Footer
		var footer = &footer_no_ptr
		container.Append(elements.Div().AttrClass(opts.Prefix + "dropdown-footer").Append(footer))
	}
	header.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		menu.JSExtElement().ClassList().Toggle(opts.Prefix + "show")
	})

	var pos = "relative"
	if absolutePosition {
		pos = "absolute"
	}

	var css = `
		.` + opts.Prefix + `dropdown {
			position: relative;
			display: inline-block;
			width: ` + opts.Width + `;
		}
		.` + opts.Prefix + `dropdown-header {
			background-color: ` + opts.Background + `;
			color: ` + opts.Color + `;
			text-decoration: none;
			display: block;
			border: ` + opts.ButtonBorderWidth + ` solid ` + opts.Color + `;
			width: ` + opts.Width + `;
			cursor: pointer;
		}
		.` + opts.Prefix + `dropdown-content {
			width: calc(` + opts.Width + ` - ` + opts.BorderWidth + ` * 2 + ` + opts.ButtonBorderWidth + ` * 2);
			height: 0px;
			position: ` + pos + `;
			background-color: ` + opts.Background + `;
			border: none;
			overflow: hidden;
			box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
			z-index: 1;
			transition: height 0.5s, border 0.5s;
		}
		.` + opts.Prefix + `dropdown-content > * {
			opacity: 0;
			transition: opacity 0.3s;
		}
		.` + opts.Prefix + `dropdown-item {
			color: ` + opts.Color + `;
			text-decoration: none;
			display: block;
			width: 100%;
			height: 100%;
		}
		.` + opts.Prefix + `dropdown-item:not(:last-child) {
			border-bottom: ` + opts.BorderWidth + ` solid ` + opts.Color + `;
		}
		.` + opts.Prefix + `dropdown a:hover {background-color: #f1f1f1}
		.` + opts.Prefix + `show {height: ` + opts.Height + `; border:` + opts.BorderWidth + ` solid ` + opts.Color + `;}
		.` + opts.Prefix + `show > * {
			opacity: 1;
		}
	`

	container.StyleBlock(css)

	return container
}
