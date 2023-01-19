//go:build js && wasm
// +build js,wasm

package misc

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
	"github.com/Nigel2392/jsext/framework/helpers"
	"github.com/Nigel2392/jsext/framework/helpers/convert"
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

func NewCounter(elem *elements.Element) *Counter {
	return &Counter{
		count:   0,
		element: elem,
	}
}

type Counter struct {
	count   int
	element *elements.Element
}

func (c *Counter) Increment() {
	c.count++
	c.element.InnerHTML(strconv.Itoa(c.count))
}

func (c *Counter) Decrement() {
	c.count--
	c.element.InnerHTML(strconv.Itoa(c.count))
}

func (c *Counter) Reset() {
	c.count = 0
	c.element.InnerHTML(strconv.Itoa(c.count))
}

func (c *Counter) Set(i int) {
	c.count = i
	c.element.InnerHTML(strconv.Itoa(c.count))
}

func (c *Counter) Get() int { return c.count }

func (c *Counter) Add(i int) {
	c.count += i
	c.element.InnerHTML(strconv.Itoa(c.count))
}

func (c *Counter) Sub(i int) {
	c.count -= i
	c.element.InnerHTML(strconv.Itoa(c.count))
}

type TimeCounter struct {
	time       time.Time
	element    *elements.Element
	ticker     *time.Ticker
	formatFunc func(*convert.TimeTracker) string
	updateFunc func(*convert.TimeTracker, *elements.Element)
}

func NewTimeCounter(elem *elements.Element) *TimeCounter {
	return &TimeCounter{
		time:    time.Now(),
		element: elem,
		formatFunc: func(t *convert.TimeTracker) string {
			return t.Format("%YR years, %MO months, %DD days, %HH hours, %MM minutes, %SS seconds")
		},
	}
}

func (c *TimeCounter) Increment() { c.Add(time.Second) }

func (c *TimeCounter) FormatFunc(f func(*convert.TimeTracker) string)             { c.formatFunc = f }
func (c *TimeCounter) UpdateFunc(f func(*convert.TimeTracker, *elements.Element)) { c.updateFunc = f }

func (c *TimeCounter) Display(Time time.Time) {
	if c.formatFunc == nil {
		panic("Format function is nil")
	}
	if c.updateFunc != nil {
		c.updateFunc(c.Tracker(), c.element)
		return
	}
	// Display time until
	timeTracker := convert.NewTimeTracker(Time)
	c.element.InnerHTML(c.formatFunc(timeTracker))
}

func (c *TimeCounter) Reset() {
	c.time = time.Now()
	c.Display(c.time)
}

func (c *TimeCounter) Set(t time.Time) {
	c.time = t
	c.Display(c.time)
}

func (c *TimeCounter) Get() time.Time { return c.time }

func (c *TimeCounter) Add(t time.Duration) {
	c.time = c.time.Add(t)
	c.Display(c.time)
}

func (c *TimeCounter) Sub(t time.Duration) {
	c.time = c.time.Add(t)
	c.Display(c.time)
}

func (c *TimeCounter) SetFormat(format string) {
	c.Display(c.time)
}

func (c *TimeCounter) Live() {
	c.Display(c.time)
	c.ticker = time.NewTicker(time.Second)
	go func() {
		for range c.ticker.C {
			c.Display(c.time)

		}
	}()
}

func (c *TimeCounter) StopLive() { c.ticker.Stop() }

func (c *TimeCounter) Tracker() *convert.TimeTracker { return convert.NewTimeTracker(c.time) }

func (c *TimeCounter) Date(year int, month time.Month, day int, hour int, min int, sec int, nsec int) {
	c.Set(time.Now().Add(time.Since(time.Date(year, month, day, hour, min, sec, nsec, time.UTC)) * -1))
}

type RoadMapStyle int

const (
	RoadMapStyleOne RoadMapStyle = iota
	RoadMapStyleTwo
)

type RoadMapItem struct {
	Name         string
	Title        string
	TitleElement *elements.Element
	Description  string
	Tags         []string
	StartDate    string
	EndDate      string
}

type Translations struct {
	To      string
	Present string
}

type RoadMapOptions struct {
	Background          string
	ItemBackground      string
	TagBackgroundColors []string
	Color               string
	TitleColor          string
	TagColor            string
	DivisorColor        string
	DivisorWidth        string
	CardMargin          string
	CardBorderWidth     string
	CardBorderColor     string
	Width               string
	Items               []RoadMapItem
	classPrefix         string
	Style               RoadMapStyle
	Translations        Translations
	FontScale           float64
}

func (r *RoadMapOptions) defaultOverrides() {
	if r.Background == "" {
		r.Background = "#ffffff"
	}
	if r.ItemBackground == "" {
		r.ItemBackground = `rgb(85,34,195);
		background: linear-gradient(0deg, rgba(85,34,195,0.40940126050420167) 0%, rgba(113,49,172,1) 100%)`
	}
	if r.Color == "" {
		r.Color = "#ffffff"
	}
	if r.TitleColor == "" {
		r.TitleColor = "#ffffff"
	}
	if r.TagColor == "" {
		r.TagColor = "#ffffff"
	}
	if r.TagBackgroundColors == nil {
		r.TagBackgroundColors = []string{"#ffffff"}
	}
	if r.DivisorColor == "" {
		r.DivisorColor = "#333333"
	}
	if r.DivisorWidth == "" {
		r.DivisorWidth = "1px"
	}
	if r.Width == "" {
		r.Width = "100%"
	}
	if r.classPrefix == "" {
		r.classPrefix = "jsext-jobtree-"
	}
	if r.CardMargin == "" {
		r.CardMargin = "20px"
	}
	if r.CardBorderWidth == "" {
		r.CardBorderWidth = "1px"
	}
	if r.CardBorderColor == "" {
		r.CardBorderColor = "#333333"
	}
	if r.Translations.To == "" {
		r.Translations.To = "to"
	}
	if r.Translations.Present == "" {
		r.Translations.Present = "Present"
	}
	if r.FontScale == 0 {
		r.FontScale = 1
	}
}

func RoadMap(roadMap *RoadMapOptions) *elements.Element {
	roadMap.defaultOverrides()

	switch roadMap.Style {
	case RoadMapStyleOne:
		return roadMapStyleOne(roadMap)
	case RoadMapStyleTwo:
		return roadMapStyleTwo(roadMap)
	default:
		return roadMapStyleOne(roadMap)
	}
}

func delimitRoadMapCSS(roadMap *RoadMapOptions) string {
	var cardMarginCalc string
	if roadMap.CardMargin != "" && roadMap.CardMargin != "0" {
		cardMarginCalc = `- calc(` + roadMap.CardMargin + ` / 2)`
	}
	return `.` + roadMap.classPrefix + `timeline {
		position: relative;
		width: ` + roadMap.Width + `;
	  }
	  
	  .` + roadMap.classPrefix + `timeline::after {
		content: '';
		position: absolute;
		width: ` + roadMap.DivisorWidth + `;
		background: ` + roadMap.DivisorColor + `;
		top: calc(100% / ` + strconv.Itoa(len(roadMap.Items)) + ` / 2 ` + cardMarginCalc + `);
		bottom: calc(100% / ` + strconv.Itoa(len(roadMap.Items)) + ` / 2 ` + cardMarginCalc + `);
		left: 50%;
		margin-left: calc(` + roadMap.DivisorWidth + ` / -2);
	  }`
}

func roadMapStyleOne(roadMap *RoadMapOptions) *elements.Element {
	var timeline = elements.Div().AttrClass(roadMap.classPrefix + "timeline")

	for i, item := range roadMap.Items {
		var container = timeline.Div().AttrClass(roadMap.classPrefix + "container")
		if i%2 == 0 {
			container.AttrClass(roadMap.classPrefix + "right")
		} else {
			container.AttrClass(roadMap.classPrefix + "left")
		}
		container.Div(item.StartDate + " - " + item.EndDate).AttrClass(roadMap.classPrefix + "date")
		var content = container.Div().AttrClass(roadMap.classPrefix + "content")
		content.H1(item.Name)
		content.H2(item.Title)
		if item.Description != "" {
			content.P(item.Description)
		}
		var paragraph = content.P()
		var ct = 0
		var color string
		for _, tag := range item.Tags {
			color, ct = helpers.GetColor(roadMap.TagBackgroundColors, ct, "#5555ff")
			paragraph.Span(tag).AttrClass(roadMap.classPrefix + "content-tag-item").AttrStyle("background-color:" + color)
		}
	}
	roadMap.CardMargin = "0"
	var css = delimitRoadMapCSS(roadMap) + `	
	.` + roadMap.classPrefix + `container {
	  position: relative;
	  background: inherit;
	  width: 50%;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `left {
	  left: 0;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right {
	  left: 50%;
	}
	.` + roadMap.classPrefix + `container::after {
	  content: '';
	  position: absolute;
	  width: 16px;
	  height: 16px;
	  top: calc(50% - 8px);
	  right: -8px;
	  background: #ffffff;
	  border: 2px solid ` + roadMap.DivisorColor + `;
	  border-radius: 16px;
	  z-index: 1;
	  transform: translateX(200%);
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right::after {
	  left: -8px;
	  transform: translateX(-200%);
	}
	.` + roadMap.classPrefix + `container::before {
	  content: '';
	  position: absolute;
	  width: 50px;
	  height: 2px;
	  top: calc(50% - 1px);
	  right: 8px;
	  background: ` + roadMap.DivisorColor + `;
	  z-index: 1;
	  transform: translateX(100%);
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right::before {
	  left: 8px;
	  transform: translateX(-100%);
	}
	.` + roadMap.classPrefix + `container .` + roadMap.classPrefix + `date {
	  position: absolute;
	  display: inline-block;
	  top: calc(50% - 8px);
	  text-align: center;
	  font-size: calc(14px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  font-weight: bold;
	  color: ` + roadMap.DivisorColor + `;
	  text-transform: uppercase;
	  letter-spacing: 1px;
	  z-index: 1;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `left .` + roadMap.classPrefix + `date {
		right:1em;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `date {
	  	left: 1em;
	}
	.` + roadMap.classPrefix + `container .` + roadMap.classPrefix + `icon {
	  position: absolute;
	  display: inline-block;
	  width: calc(40px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  height: calc(40px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  padding: 9px 0;
	  top: calc(50% - calc(calc(40px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) / 2));
	  background: #F6D155;
	  border: 2px solid ` + roadMap.DivisorColor + `;
	  border-radius: calc(40px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  text-align: center;
	  font-size: calc(18px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);;
	  color: ` + roadMap.DivisorColor + `;
	  z-index: 1;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `left .` + roadMap.classPrefix + `icon {
	  right: 56px;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `icon {
	  left: 56px;
	}
	.` + roadMap.classPrefix + `container .` + roadMap.classPrefix + `content {
	  color: ` + roadMap.Color + `;
	  padding: calc(30px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) calc(90px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) calc(30px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) calc(30px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  background: ` + roadMap.ItemBackground + `;
	  position: relative;
	  border-radius: 0 500px 500px 0;
	}
	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `content {
	  padding: calc(30px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) calc(30px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) calc(30px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `) calc(90px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  border-radius: 500px 0 0 500px;
	  text-align: right;
	}
	.` + roadMap.classPrefix + `container .` + roadMap.classPrefix + `content h1 {
	  	margin: 0 0 10px 0;
		font-size: calc(20px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-weight: bold;
		color: ` + roadMap.TitleColor + `;
	}
	.` + roadMap.classPrefix + `container .` + roadMap.classPrefix + `content h2 {
	  margin: 0 0 10px 0;
	  font-size: calc(18px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  font-weight: normal;
	  color: ` + roadMap.TitleColor + `;
	}
	.` + roadMap.classPrefix + `container .` + roadMap.classPrefix + `content p {
	  margin: 0;
	  font-size: calc(16px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	  line-height: 22px;
	  color: #000000;
	}
	.` + roadMap.classPrefix + "content-tag-item" + `{
		display: inline-block;
		padding: 2px 5px;
		margin: 0 5px 5px 0;
		border-radius: 5px;
		color: ` + roadMap.TagColor + `;
	}
	@media (max-width: 767.98px) {
	  	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `left .` + roadMap.classPrefix + `date {
			right: -100%;
	  	}
	  	.` + roadMap.classPrefix + `container.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `date {
			left: -100%;
		}
	}
	`
	timeline.StyleBlock(css)
	return timeline
}

func roadMapStyleTwo(roadMap *RoadMapOptions) *elements.Element {
	var timeline = elements.Div().AttrClass(roadMap.classPrefix + "timeline")

	for i, item := range roadMap.Items {
		var container = timeline.Div().AttrClass(roadMap.classPrefix + "card-container")
		var card = container.Div().AttrClass(roadMap.classPrefix + "card")
		if i%2 == 0 {
			container.AttrClass(roadMap.classPrefix + "left")
		} else {
			container.AttrClass(roadMap.classPrefix + "right")
		}

		var card_header = card.Div().AttrClass(roadMap.classPrefix + "card-header")

		if item.Description != "" && len(item.Tags) != 0 {
			var card_body = card.Div().AttrClass(roadMap.classPrefix + "card-body")
			if item.Description != "" {
				card_body.P(item.Description)
			}
			var tagsParagraph = card_body.P()
			for i, tag := range item.Tags {
				var spacing = ""
				if i < len(item.Tags)-1 {
					spacing = ", "
				}
				tagsParagraph.Span(tag + spacing).AttrClass(roadMap.classPrefix + "content-tag-item")
			}
		}

		var (
			card_footer  = card.Div().AttrClass(roadMap.classPrefix + "card-footer")
			card_company = card.Div().AttrClass(roadMap.classPrefix + "card-name")
		)

		if item.TitleElement != nil {
			card_header.Append(item.TitleElement)
		} else {
			card_header.Div(item.Title)
		}

		if item.StartDate != "" || item.EndDate != "" {
			if item.StartDate != "" && item.EndDate != "" {
				card_footer.Div(item.StartDate + " " + roadMap.Translations.To + " " + item.EndDate)
			} else if item.StartDate != "" {
				card_footer.Div(item.StartDate + " - " + roadMap.Translations.Present)
			} else {
				card_footer.Div(roadMap.Translations.Present)
			}
		}

		card_company.Div(item.Name)
		card.Animations.FadeIn(500, elements.UseIntersectionObserver, true)
	}

	var css = delimitRoadMapCSS(roadMap) + `
	.` + roadMap.classPrefix + `card-container {
		position: relative;
		width: 100%;
	}
	.` + roadMap.classPrefix + `card-container::before {
		content: "";
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translateY(-50%) translateX(-50%);
		width: calc(20px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		height: calc(20px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		border-radius: 50%;
		background: ` + roadMap.DivisorColor + `;
		z-index: 1;
	}
	.` + roadMap.classPrefix + `card {
		position: relative;
		width: calc(50% + calc(` + roadMap.DivisorWidth + ` / 2) - calc(` + roadMap.CardBorderWidth + ` * 2));
		background: ` + roadMap.ItemBackground + `;
		border-radius: 5px;
		padding: 10px 0;
		margin-bottom:` + roadMap.CardMargin + `;
		color: ` + roadMap.Color + `;
		border: ` + roadMap.CardBorderWidth + ` solid ` + roadMap.CardBorderColor + `;
	}
	.` + roadMap.classPrefix + `left .` + roadMap.classPrefix + `card {
		left: 0;
		text-align: left;
		margin-right: auto;
		box-shadow: 5px 5px 5px 0 rgb(0 0 0 / 20%)
	}
	.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `card {
		right: 0;
		text-align: right;
		margin-left: auto;
		box-shadow: 5px 5px 5px 0 rgb(0 0 0 / 20%)
	}
	.` + roadMap.classPrefix + `card .` + roadMap.classPrefix + `card-header {
		line-height: calc(22px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-size: calc(22px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-weight: bold;
		color: ` + roadMap.TitleColor + `;
		border-bottom: ` + roadMap.CardBorderWidth + ` solid ` + roadMap.CardBorderColor + `;
		padding: 0 10px;
		padding-bottom: 5px;
	}
	.` + roadMap.classPrefix + `card .` + roadMap.classPrefix + `card-body {
		padding: 0 10px;
	}
	.` + roadMap.classPrefix + `card .` + roadMap.classPrefix + `card-body * {
		line-height: calc(16px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-size: calc(16px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
	}
	.` + roadMap.classPrefix + `card .` + roadMap.classPrefix + `card-footer {
		line-height: calc(14px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-size: calc(14px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-weight: bold;
		color: ` + roadMap.TitleColor + `;
		margin-bottom: 0;
		margin-top: auto;
		padding: 0 10px;
		padding-top: 5px;
		border-top: ` + roadMap.CardBorderWidth + ` solid ` + roadMap.CardBorderColor + `;
	}
	.` + roadMap.classPrefix + `left .` + roadMap.classPrefix + `card-footer {
		text-align: right;
	}
	.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `card-footer {
		text-align: left;
	}
	.` + roadMap.classPrefix + `card .` + roadMap.classPrefix + `card-name {
		position: absolute;
		font-size: calc(24px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `);
		font-weight: bold;
		top: 50%;
		transform: translateY(-50%);
		color: ` + roadMap.DivisorColor + `;
		width: 100%;
	}
	.` + roadMap.classPrefix + `left .` + roadMap.classPrefix + `card-name {
		left: calc(100% + ` + roadMap.DivisorWidth + ` + calc(10px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `));
	}
	.` + roadMap.classPrefix + `right .` + roadMap.classPrefix + `card-name {
		right: calc(100% + ` + roadMap.DivisorWidth + ` + calc(10px *` + strconv.FormatFloat(roadMap.FontScale, 'f', 2, 64) + `));
	}
	`
	timeline.StyleBlock(css)

	return timeline
}

type Modal elements.Element

func (m *Modal) e() *elements.Element {
	return (*elements.Element)(m)
}

func (m *Modal) Show() {
	if !m.e().Value().Truthy() {
		m.Create()
	}
	m.e().AttrStyle("display:flex")
}

func (m *Modal) Hide() {
	m.e().AttrStyle("display:none")
}

func (m *Modal) Render() jsext.Element {
	return m.e().Render()
}

func (m *Modal) Create(appendToQuerySelector ...string) {
	if len(appendToQuerySelector) > 0 && appendToQuerySelector[0] != "" {
		var e, err = jsext.QuerySelector(appendToQuerySelector[0])
		if err != nil {
			panic(err)
		}
		e.Append(m.Render())
	} else {
		jsext.Body.Append(m.Render())
	}
}

func (m *Modal) OpenOnClickOf(e *elements.Element) {
	e.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		var preventDefault = event.Get("preventDefault")
		if preventDefault.Truthy() {
			event.PreventDefault()
		}
		m.Show()
	})
}

func (m *Modal) CloseOnClickOf(e *elements.Element) {
	e.AddEventListener("click", func(this jsext.Value, event jsext.Event) {
		var preventDefault = event.Get("preventDefault")
		if preventDefault.Truthy() {
			event.PreventDefault()
		}
		m.Delete()
	})
}

type ButtonType int

const (
	Anchor ButtonType = 0
	Button ButtonType = 1
)

func (m *Modal) Button(tag ButtonType, innerText string) *elements.Element {
	var btn *elements.Element
	switch tag {
	case Anchor:
		btn = elements.A("javascript:void(0)", innerText)
	case Button:
		btn = elements.Button(innerText)
	}
	m.OpenOnClickOf(btn)
	return btn
}

func (m *Modal) Delete() {
	if m.e().Value().Truthy() {
		m.e().Value().Remove()
	}
}

type ModalOptions struct {
	Header           *elements.Element
	Body             *elements.Element
	Footer           *elements.Element
	Background       string
	ModalBackground  string
	BorderRadius     string
	Border           string
	Width            string
	Height           string
	ClassPrefix      string
	CloseButton      bool
	CloseButtonScale float64
	ZIndex           int
}

func (opts *ModalOptions) SetDefaults() {
	if opts.ClassPrefix == "" {
		opts.ClassPrefix = "jsext-modal-"
	}
	if opts.Background == "" {
		opts.Background = "rgba(0,0,0,0.5)"
	}
	if opts.ModalBackground == "" {
		opts.ModalBackground = "#fff"
	}
	if opts.BorderRadius == "" {
		opts.BorderRadius = "5px"
	}
	if opts.Border == "" {
		opts.Border = "1px solid #ccc"
	}
	if opts.Width == "" {
		opts.Width = "50%"
	}
	if opts.Height == "" {
		opts.Height = "auto"
	}
	if opts.CloseButtonScale == 0 {
		opts.CloseButtonScale = 1
	}
	if opts.ZIndex == 0 {
		opts.ZIndex = 9999
	}
}

func CreateModal(opts ModalOptions) *Modal {
	opts.SetDefaults()
	var modal_container *elements.Element = elements.Div().AttrClass(opts.ClassPrefix + "modal-container")
	if opts.CloseButton {
		var close_btn = modal_container.Div().AttrClass(opts.ClassPrefix + "close-btn")
		close_btn.AddEventListener("click", func(this jsext.Value, e jsext.Event) {
			(*Modal)(modal_container).Hide()
		})
	}
	var modal = modal_container.Div().AttrClass(opts.ClassPrefix + "modal")
	if opts.Header != nil {
		modal.Append(opts.Header)
	}
	if opts.Body != nil {
		modal.Append(opts.Body)
	}
	if opts.Footer != nil {
		modal.Append(opts.Footer)
	}

	css := `
		.` + opts.ClassPrefix + `modal-container {
			position: fixed;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background: ` + opts.Background + `;
			display: none;
			justify-content: center;
			align-items: center;
			z-index: ` + strconv.Itoa(opts.ZIndex) + `;
		}
		.` + opts.ClassPrefix + `close-btn {
			position: absolute;
			top: 10px;
			right: 10px;
			width: 30px;
			height: 30px;
			background: #ff0000;
			border-radius: 50%;
			display: flex;
			justify-content: center;
			align-items: center;
			cursor: pointer;
			transform: scale(` + strconv.FormatFloat(opts.CloseButtonScale, 'f', 2, 64) + `);
			transition: background 0.2s ease-in-out;
		}
		.` + opts.ClassPrefix + `close-btn:hover {
			background: #910c0c;
		}
		.` + opts.ClassPrefix + `close-btn::before {
			position: absolute;
			content: "";
			width: 20px;
			height: 2px;
			background: #fff;
			transform: rotate(45deg);
		}
		.` + opts.ClassPrefix + `close-btn::after {
			position: absolute;
			content: "";
			width: 20px;
			height: 2px;
			background: #fff;
			transform: rotate(-45deg);
		}
		
		.` + opts.ClassPrefix + `modal {
			background: ` + opts.ModalBackground + `;
			border-radius: ` + opts.BorderRadius + `;
			border: ` + opts.Border + `;
			max-width: 95%;
			overflow-x: auto;
			width: ` + opts.Width + `;
			max-height: 95%;
			overflow-y: auto;
			height: ` + opts.Height + `;
			display: flex;
			flex-direction: column;
		}
		.` + opts.ClassPrefix + `modal > * {
			padding: 10px;
		}
		.` + opts.ClassPrefix + `modal > *:first-child {
			border-bottom: 1px solid #ccc;
		}
		.` + opts.ClassPrefix + `modal > *:last-child {
			border-top: 1px solid #ccc;
		}
		.` + opts.ClassPrefix + `modal > *:only-child {
			border: none;
		}
		`

	modal.StyleBlock(css)
	return (*Modal)(modal_container)
}

type JiggleOptions struct {
	ChangeColor string
	ClassPrefix string
	Words       bool
}

func (opts *JiggleOptions) SetDefaults() {
	if opts.ChangeColor == "" {
		opts.ChangeColor = "red"
	}
	if opts.ClassPrefix == "" {
		opts.ClassPrefix = "jsext"
	}
}

func (opts *JiggleOptions) MainElementClass() string {
	opts.SetDefaults()
	var hash = helpers.FNVHashString(opts.ChangeColor)
	return opts.ClassPrefix + "-jiggle-" + hash
}

func JiggleText(tag, text string, opts *JiggleOptions) *elements.Element {
	var options = *opts
	options.SetDefaults()
	var jiggle = elements.Animation{
		Options: map[string]interface{}{
			"duration":   200,
			"iterations": 1,
			"fill":       "forwards",
			"easing":     "ease-out",
		},
	}
	var textLen int
	var split []string
	var extra = ""
	if options.Words {
		// Extra css
		extra = "padding: 0px 0.15em;"
		// Split the words
		split = helpers.SplitWords(text)
		textLen = len(split)
		jiggle.Animations = []any{
			map[string]interface{}{"transform": "scale(1) rotate(0deg)", "offset": "0.0"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.1)", "offset": "0.125"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.15) rotate(-3deg)", "offset": "0.25"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.25)", "offset": "0.375"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.15) rotate(3deg)", "offset": "0.75"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.1)", "offset": "0.875"},
			map[string]interface{}{"transform": "scale(1) rotate(0deg)", "offset": "1.0"},
		}
	} else {
		// Extra css
		extra = "white-space: pre;"
		// Get the length of the text
		textLen = len(text)
		jiggle.Animations = []any{
			map[string]interface{}{"transform": "scale(1) rotate(0deg)", "offset": "0.0"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.1)", "offset": "0.125"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.15) rotate(-8deg)", "offset": "0.25"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.25)", "offset": "0.375"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.15) rotate(8deg)", "offset": "0.75"},
			map[string]interface{}{"color": options.ChangeColor, "transform": "scale(1.1)", "offset": "0.875"},
			map[string]interface{}{"transform": "scale(1) rotate(0deg)", "offset": "1.0"},
		}
	}
	// Create the main element
	var hash = helpers.FNVHashString(options.ChangeColor)
	var main = elements.NewElement(tag).AttrClass(options.ClassPrefix + "-jiggle-" + hash)
	// Initialize a slice large enough to hold all the characters/words
	main.Children = make([]*elements.Element, textLen)
	// Add the spans of characters/words
	for i := 0; i < textLen; i++ {
		switch options.Words {
		case true:
			main.Children[i] = elements.Span(split[i] + " ")
		default:
			main.Children[i] = elements.Span(string(text[i]))
		}
		main.Children[i].AttrClass(options.ClassPrefix + "-jiggle-text-" + hash)
	}

	// Set up the eventlisteners
	go func(j *elements.Animation) {
		for i := 0; i < len(main.Children); i++ {
			main.Children[i].AddEventListener("mouseover", func(this jsext.Value, event jsext.Event) {
				this.Call("animate", jsext.SliceToArray(jiggle.Animations).Value(), jsext.MapToObject(jiggle.Options).Value())
			})
		}
	}(&jiggle)

	jsext.StyleBlock(hash, `
	.`+options.ClassPrefix+`-jiggle-`+hash+` {
		pointer-events: none;
		display: flex;
		margin: 1em 0px;
		flex-wrap: wrap;
	}
	.`+options.ClassPrefix+`-jiggle-text-`+hash+` {
		display: inline;
		pointer-events: all;
		transform-origin: center;
		transition: 0.2s;
		`+extra+`
	}
	.`+options.ClassPrefix+`-jiggle-text-`+hash+`:hover {
		color: `+options.ChangeColor+`;
	}
	`)

	return main
}

type Grid interface {
	Columns() []*elements.Element
	Rows() []*elements.Element
	Element() *elements.Element
	Render() jsext.Element
}

var spacing = " "
var jsextPrefix = "jsext-"

type gridFromPattern elements.Element

// Returns the rows of the grid.
func (g *gridFromPattern) Rows() []*elements.Element {
	var grid = elements.Element(*g)
	return grid.Children
}

// Returns the columns of the grid.
func (g *gridFromPattern) Columns() []*elements.Element {
	var grid = elements.Element(*g)
	var columns = make([]*elements.Element, 0)
	for _, row := range grid.Children {
		columns = append(columns, row.Children...)
	}
	return columns
}

// Returns the inner element of the grid.
func (g *gridFromPattern) Element() *elements.Element {
	var grid = elements.Element(*g)
	return &grid
}

// Renders the grid.
// Adheres to the component interface.
func (g *gridFromPattern) Render() jsext.Element {
	return g.Element().Render()
}

// Provide a grid based on a pattern.
// Returns the rows and columns of the grid.
func GridFromPattern(gridItemCallback func(row int, column int, index int) []*elements.Element, gridPattern ...string) Grid {
	var gridContainer = elements.Div().AttrClass(jsextPrefix + "grid-container")
	gridContainer.Children = make([]*elements.Element, len(gridPattern))
	var index = 0
	for i, pattern := range gridPattern {
		var className = jsextPrefix + helpers.FNVHashString(pattern)
		gridContainer.Children[i] = elements.Div().AttrClass(className + "grid-row")

		// If the pattern is the same as the previous pattern, then we can just copy the previous row.
		if i != 0 && gridPattern[i-1] == pattern {
			var gridColumnPtr = gridContainer.Children[i-1].Children[0]
			var gridColumn = *gridColumnPtr
			var newColumn = &gridColumn
			gridContainer.Children[i].Children = append(gridContainer.Children[i].Children, newColumn)
			continue
		}

		var splitGrid = strings.Split(pattern, spacing)
		var fractionSlice = make([]string, len(splitGrid))
		if len(splitGrid) == 0 {
			return nil
		}

		for j, v := range splitGrid {
			fractionSlice[j] = strconv.Itoa(len(v)) + "fr"
			gridContainer.Children[i].Div().AttrClass(className + "-grid-item")
			if gridItemCallback != nil {
				var elems = gridItemCallback(i, j, index)
				if len(elems) > 0 {
					gridContainer.Children[i].Children[j].Children = elems
				}
				index++
			}
		}
		go func() {
			var fractionString = strings.Join(fractionSlice, " ")
			var css = `
			.` + className + `grid-row {
				display: grid;
				grid-template-columns: ` + fractionString + `;
				grid-template-rows: auto;
				grid-gap: 1em;
				width: 100%;
			}
			.` + className + `-grid-item {
				width: 100%;
			}`
			jsext.StyleBlock(className, css)
		}()
	}
	var gridCSS = `
		.` + jsextPrefix + `grid-container {
			display: grid;
			grid-template-rows: auto;
			grid-gap: 1em;
			width: 100%;
		}
		.` + jsextPrefix + `grid-row {
			width: 100%;
		}
	`
	jsext.StyleBlock(jsextPrefix+"grid", gridCSS)
	var gridElement = gridFromPattern(*gridContainer)
	return &gridElement
}

type GridItem struct {
	Header *elements.Element
	Body   *elements.Element
	Footer *elements.Element
}

func (g *GridItem) Container(className string) *elements.Element {
	var container = elements.Div().AttrClass(className)
	if g.Header != nil {
		container.Append(g.Header)
	}
	if g.Body != nil {
		container.Append(g.Body)
	}
	if g.Footer != nil {
		container.Append(g.Footer)
	}
	return container
}

type GridOptions struct {
	// The grid items to display inside of the grid.
	GridItems []*GridItem
	// Grid width, 100%, etc.
	Width string
	// GridItemWidth grid width, 1fr, etc.
	GridItemWidth string
	// Amount of rows and columns in the grid.
	Columns int
	Rows    int
	// Column fractions
	GridColumnWidth []string
	// If there are too many grid items supplied, and ElementOnOverflow is not nil,
	// the last grid item will be replaced with the ElementOnOverflow when used in a supportive function.
	ElementOnOverflow *GridItem
	// ExtraCSS is a map with the following keys, to edit their respective elements css styles.
	// 	- grid
	// 	- grid-item
	ExtraCSS map[string]string
	// Class prefix for the grid, and its elements.
	ClassPrefix string
	// Function to be called for adding extra CSS.
	CSSFunc func(gridClass, itemClass, itemHeaderClass, itemBodyClass, itemFooterClass string) string
}

var GridDelimiter = " "

func (c *GridOptions) Fractions(fr ...string) {
	if len(fr) != c.Columns {
		panic(errors.New("the amount of fractions must be equal to the amount of columns"))
	}
	c.GridColumnWidth = fr
}

func (c *GridOptions) setDefaults() {
	if c.Width == "" {
		c.Width = "100%"
	}
	if c.GridItemWidth == "" {
		c.GridItemWidth = "1fr"
	}
	if c.ClassPrefix == "" {
		c.ClassPrefix = "jsext"
	}
	if c.Columns == 1 && c.Rows == 1 {
		panic("Cannot have 1 column and 1 row. You do not need a grid.")
	}
	if c.Columns == 0 {
		c.Columns = 2
	}
	if c.Rows == 0 {
		c.Rows = 1
	}
	if c.ExtraCSS == nil {
		c.ExtraCSS = make(map[string]string)
	}
}

func GridNoOverflow(opts *GridOptions) *elements.Element {
	var options = *opts
	options.setDefaults()
	var (
		GridItemClass       = options.ClassPrefix + `-grid-item`
		GridItemHeaderClass = options.ClassPrefix + `-grid-item-header`
		GridItemBodyClass   = options.ClassPrefix + `-grid-item-body`
		GridItemFooterClass = options.ClassPrefix + `-grid-item-footer`
	)
	var isGreater = len(options.GridItems) > options.Columns*options.Rows
	if isGreater {
		if options.ElementOnOverflow != nil {
			// Remove the last grid item and replace it with the overflow element.
			options.GridItems = options.GridItems[:options.Columns*options.Rows-1]
		} else {
			options.GridItems = options.GridItems[:options.Columns*options.Rows]
		}
	}

	var mainGrid = NewGrid(&options)

	if options.ElementOnOverflow != nil && isGreater {
		var gridItem = elements.Div().AttrClass(GridItemClass)
		gridItem.Append(
			options.ElementOnOverflow.Header.AttrClass(GridItemHeaderClass),
			options.ElementOnOverflow.Body.AttrClass(GridItemBodyClass),
			options.ElementOnOverflow.Footer.AttrClass(GridItemFooterClass),
		)
		mainGrid.Append(gridItem)
	}

	return mainGrid
}

func NewGrid(opts *GridOptions) *elements.Element {
	var options = *opts
	options.setDefaults()
	var (
		GridClass           = options.ClassPrefix + `-grid`
		gridItemClass       = options.ClassPrefix + `-grid-item`
		gridItemHeaderClass = options.ClassPrefix + `-grid-item-header`
		gridItemBodyClass   = options.ClassPrefix + `-grid-item-body`
		gridItemFooterClass = options.ClassPrefix + `-grid-item-footer`
	)

	var mainGrid = elements.Div().AttrClass(GridClass)
	mainGrid.Children = make([]*elements.Element, len(options.GridItems))
	for i, gridItem := range options.GridItems {
		var gridItemContainer = elements.Div().AttrClass(gridItemClass)
		if gridItem.Header != nil {
			gridItemContainer.Append(gridItem.Header.AttrClass(gridItemHeaderClass))
		}
		if gridItem.Body != nil {
			gridItemContainer.Append(gridItem.Body.AttrClass(gridItemBodyClass))
		}
		if gridItem.Footer != nil {
			gridItemContainer.Append(gridItem.Footer.AttrClass(gridItemFooterClass))
		}
		mainGrid.Children[i] = gridItemContainer
	}

	var colFracs = `repeat(` + strconv.Itoa(options.Columns) + `, ` + options.GridItemWidth + `)`
	if len(options.GridColumnWidth) > 0 {
		colFracs = strings.Join(options.GridColumnWidth, " ")
	}

	var css = `
		.` + GridClass + `{
			display: grid;
			grid-template-columns: ` + colFracs + `;
			grid-template-rows: repeat(` + strconv.Itoa(options.Rows) + `, auto);
			grid-gap: 1rem;
			width: ` + options.Width + `;
			` + options.ExtraCSS["grid"] + `
		}
		.` + gridItemClass + `{
			width: 100%;
			` + options.ExtraCSS["grid-item"] + `
		}
	`

	if options.CSSFunc != nil {
		css += options.CSSFunc(
			GridClass,
			gridItemClass,
			gridItemHeaderClass,
			gridItemBodyClass,
			gridItemFooterClass,
		)
	}

	mainGrid.StyleBlock(css)

	return mainGrid
}
