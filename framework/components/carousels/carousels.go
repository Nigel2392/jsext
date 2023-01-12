//go:build js && wasm
// +build js,wasm

package carousels

import (
	"strconv"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
)

type Options struct {
	// The width of the carousel.
	Width string
	// The height of the carousel.
	Height string
	// The background color of the carousel.
	Background string
	// Border of the carousel.
	Border string
	// Show the carousel controls.
	Controls      bool
	ControlsColor string
	ControlsSize  string
	// Show the carousel indicators.
	Indicators bool
	// Items
	Items []*elements.Element
	// Class prefix
	Prefix string
	// Active item
	ActiveItem int
}

func (c *Options) SetDefaults() {
	if c.Width == "" {
		c.Width = "100%"
	}
	if c.Height == "" {
		c.Height = "100%"
	}
	if c.Background == "" {
		c.Background = "#ffffff"
	}
	if c.ControlsColor == "" {
		c.ControlsColor = "#000000"
	}
	if c.Prefix == "" {
		c.Prefix = "jsext-carousel-"
	}
	if c.ControlsSize == "" {
		c.ControlsSize = "calc(1.25em + 1.5vw)"
	}
	if c.Border == "" {
		c.Border = "none"
	}
}

func Plain(options *Options) *elements.Element {
	if len(options.Items) == 0 {
		panic("Carousel requires at least one item.")
	}
	options.SetDefaults()

	var CarouselContainer = elements.Div().AttrClass(options.Prefix + "container")
	var Carousel = CarouselContainer.Div().AttrClass(options.Prefix + "carousel")
	var CarouselInner = Carousel.Div().AttrClass(options.Prefix + "carousel-inner")

	var activeItem = options.ActiveItem
	if activeItem < 0 || activeItem >= len(options.Items) {
		activeItem = 0
	}

	var items = make([]*elements.Element, len(options.Items))
	for i, item := range options.Items {
		var itemClass = options.Prefix + "carousel-item"
		items[i] = item
		if i == activeItem {
			CarouselInner.Append(item.AttrClass(itemClass, "active"))
		} else {
			CarouselInner.Append(item.AttrClass(itemClass))
		}
	}

	var css = `.` + options.Prefix + `container {
		width: ` + options.Width + `;
		height: ` + options.Height + `;
		background-color: ` + options.Background + `;
		position: relative;
		border: ` + options.Border + `;
		border-radius: 0.25em;
		padding: calc(0.5em + 0.5vw);
	}

	.` + options.Prefix + `carousel {
		width: 100%;
		height: 100%;
	}
	.` + options.Prefix + `carousel-inner {
		width: 100%;
		height: 100%;
		position: relative;
	}
	.` + options.Prefix + `carousel-item {
		width: ` + options.Width + `;
		height: ` + options.Height + `;
		position: absolute;
		top: 0;
		left: 0;
		opacity: 0;
		transition: opacity 0.5s;
		object-fit: contain;
	}
	.` + options.Prefix + `carousel-item.active {
		opacity: 1;
	}`
	var indicatorList = make([]*elements.Element, len(options.Items))

	if options.Controls {
		var left = CarouselContainer.Div().AttrClass(options.Prefix+"arrow-left", options.Prefix+"arrow")
		var right = CarouselContainer.Div().AttrClass(options.Prefix+"arrow-right", options.Prefix+"arrow")
		// CarouselContainer.Div().AttrClass(options.Prefix+"arrow-left", options.Prefix+"arrow")
		// CarouselContainer.Div().AttrClass(options.Prefix+"arrow-right", options.Prefix+"arrow")
		css += `
		.` + options.Prefix + `arrow {
			opacity: 0;
			width: 0; 
			height: 0; 
			border-top: ` + options.ControlsSize + ` solid transparent;
			border-bottom: ` + options.ControlsSize + ` solid transparent; 
			position: absolute;
			top: 50%;
			transform: translateY(-50%);
			transition: opacity 0.3s ease-in;
			cursor: pointer;
		  }
		.` + options.Prefix + `container:hover .` + options.Prefix + `arrow {
			opacity: 1;
		}
		.` + options.Prefix + `arrow-right {
			border-left: ` + options.ControlsSize + ` solid ` + options.ControlsColor + `;
			right: calc(0.5em + 0.5vw);
	    }
	    .` + options.Prefix + `arrow-left {
			border-right:` + options.ControlsSize + ` solid ` + options.ControlsColor + `; 
			left: calc(0.5em + 0.5vw);
	    }`
		var setActive = func(i int) {
			items[i].JSExtElement().ClassList().Add("active")
			if options.Indicators {
				indicatorList[i].JSExtElement().ClassList().Add("active")
			}
		}

		left.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
			for i, item := range items {
				if item.JSExtElement().ClassList().Call("contains", "active").Bool() {
					item.JSExtElement().ClassList().Remove("active")
					if options.Indicators {
						indicatorList[i].JSExtElement().ClassList().Remove("active")
					}
					if i == 0 {
						setActive(len(items) - 1)
					} else {
						setActive(i - 1)
					}
					break
				}
			}
		})
		right.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
			for i, item := range items {
				if item.JSExtElement().ClassList().Call("contains", "active").Bool() {
					item.JSExtElement().ClassList().Remove("active")
					if options.Indicators {
						indicatorList[i].JSExtElement().ClassList().Remove("active")
					}
					if i == len(items)-1 {
						setActive(0)
					} else {
						setActive(i + 1)
					}
					break
				}
			}
		})
	}
	if options.Indicators {
		var indicators = CarouselContainer.Div().AttrClass(options.Prefix + "indicators")
		for i := range options.Items {
			var indicator = indicators.Div().AttrClass(options.Prefix + "indicator")
			if i == activeItem {
				indicator.AttrClass("active")
			}
			indicator.Set("data-index", strconv.Itoa(i))
			indicator.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
				indexInt, _ := strconv.Atoi(this.Get("dataset").Get("index").String())
				for i, item := range items {
					if i == indexInt {
						item.JSExtElement().ClassList().Add("active")
					} else {
						item.JSExtElement().ClassList().Remove("active")
					}
				}
				for _, indicator := range indicatorList {
					indicator.JSExtElement().ClassList().Remove("active")
				}
				this.Get("classList").Call("add", "active")
			})
			indicatorList[i] = indicator
		}
		css += `
		.` + options.Prefix + `indicators {
			position: absolute;
			bottom: 0;
			left: 0;
			width: 100%;
			display: flex;
			justify-content: center;
			align-items: center;
			padding: 10px 0;
		}
		.` + options.Prefix + `indicator {
			width: 10px;
			height: 10px;
			border-radius: 50%;
			background-color: ` + options.Background + `;
			border: 1px solid ` + options.ControlsColor + `;
			margin: 0 5px;
			cursor: pointer;
			transition: background-color 0.5s;
		}
		.` + options.Prefix + `indicator.active {
			background-color: ` + options.ControlsColor + `;
			border: 1px solid ` + options.Background + `;
		}`
	}

	CarouselContainer.StyleBlock(css)

	return CarouselContainer
}

func Image(imageUrls []string, options *Options) *elements.Element {
	var items = make([]*elements.Element, len(imageUrls))
	for i, url := range imageUrls {
		items[i] = elements.Img(url).AttrAlt(url + " Is not available")
	}
	options.Items = items
	return Plain(options)
}
