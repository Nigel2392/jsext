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

type JobtreeItem struct {
	Company   string
	Title     string
	Tags      []string
	StartDate string
	EndDate   string
}

type Jobtree struct {
	Background          string
	ItemBackground      string
	TagBackgroundColors []string
	Color               string
	TitleColor          string
	TagColor            string
	DivisorColor        string
	DivisorWidth        string
	Width               string
	Items               []JobtreeItem
	classPrefix         string
}

func (j *Jobtree) defaultOverrides() {
	if j.Background == "" {
		j.Background = "#ffffff"
	}
	if j.ItemBackground == "" {
		j.ItemBackground = `rgb(85,34,195);
		background: linear-gradient(0deg, rgba(85,34,195,0.40940126050420167) 0%, rgba(113,49,172,1) 100%)`
	}
	if j.Color == "" {
		j.Color = "#ffffff"
	}
	if j.TitleColor == "" {
		j.TitleColor = "#ffffff"
	}
	if j.TagColor == "" {
		j.TagColor = "#ffffff"
	}
	if j.TagBackgroundColors == nil {
		j.TagBackgroundColors = []string{"#ffffff"}
	}
	if j.DivisorColor == "" {
		j.DivisorColor = "#333333"
	}
	if j.DivisorWidth == "" {
		j.DivisorWidth = "1px"
	}
	if j.Width == "" {
		j.Width = "100%"
	}
	if j.classPrefix == "" {
		j.classPrefix = "jsext-jobtree-"
	}
}

func JobTree(jobTree *Jobtree) *elements.Element {
	jobTree.defaultOverrides()

	var timeline = elements.Div().AttrClass(jobTree.classPrefix + "timeline")

	for i, item := range jobTree.Items {
		var container = timeline.Div().AttrClass(jobTree.classPrefix + "container")
		if i%2 == 0 {
			container.AttrClass(jobTree.classPrefix + "right")
		} else {
			container.AttrClass(jobTree.classPrefix + "left")
		}
		container.Div(item.StartDate + " - " + item.EndDate).AttrClass(jobTree.classPrefix + "date")
		var content = container.Div().AttrClass(jobTree.classPrefix + "content")
		content.H1(item.Company)
		content.H2(item.Title)
		var paragraph = content.P()
		var ct = 0
		var color string
		for _, tag := range item.Tags {
			color, ct = helpers.GetColor(jobTree.TagBackgroundColors, ct, "#5555ff")
			paragraph.Span(tag).AttrClass(jobTree.classPrefix + "content-tag-item").AttrStyle("background-color:" + color)
		}
	}

	var css = `
	.` + jobTree.classPrefix + `timeline {
	  position: relative;
	  width: ` + jobTree.Width + `;
	}
	
	.` + jobTree.classPrefix + `timeline::after {
	  content: '';
	  position: absolute;
	  width: ` + jobTree.DivisorWidth + `;
	  background: ` + jobTree.DivisorColor + `;
	  top: 0;
	  bottom: 0;
	  left: 50%;
	  margin-left: calc(` + jobTree.DivisorWidth + ` / -2);
	}
	
	.` + jobTree.classPrefix + `container {
	  position: relative;
	  background: inherit;
	  width: 50%;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `left {
	  left: 0;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right {
	  left: 50%;
	}
	
	.` + jobTree.classPrefix + `container::after {
	  content: '';
	  position: absolute;
	  width: 16px;
	  height: 16px;
	  top: calc(50% - 8px);
	  right: -8px;
	  background: #ffffff;
	  border: 2px solid ` + jobTree.DivisorColor + `;
	  border-radius: 16px;
	  z-index: 1;
	  transform: translateX(200%);
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right::after {
	  left: -8px;
	  transform: translateX(-200%);
	}
	
	.` + jobTree.classPrefix + `container::before {
	  content: '';
	  position: absolute;
	  width: 50px;
	  height: 2px;
	  top: calc(50% - 1px);
	  right: 8px;
	  background: ` + jobTree.DivisorColor + `;
	  z-index: 1;
	  transform: translateX(100%);
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right::before {
	  left: 8px;
	  transform: translateX(-100%);
	}
	
	.` + jobTree.classPrefix + `container .` + jobTree.classPrefix + `date {
	  position: absolute;
	  display: inline-block;
	  top: calc(50% - 8px);
	  text-align: center;
	  font-size: 14px;
	  font-weight: bold;
	  color: ` + jobTree.DivisorColor + `;
	  text-transform: uppercase;
	  letter-spacing: 1px;
	  z-index: 1;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `left .` + jobTree.classPrefix + `date {
		right:1em;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right .` + jobTree.classPrefix + `date {
	  	left: 1em;
	}
	
	.` + jobTree.classPrefix + `container .` + jobTree.classPrefix + `icon {
	  position: absolute;
	  display: inline-block;
	  width: 40px;
	  height: 40px;
	  padding: 9px 0;
	  top: calc(50% - 20px);
	  background: #F6D155;
	  border: 2px solid ` + jobTree.DivisorColor + `;
	  border-radius: 40px;
	  text-align: center;
	  font-size: 18px;
	  color: ` + jobTree.DivisorColor + `;
	  z-index: 1;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `left .` + jobTree.classPrefix + `icon {
	  right: 56px;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right .` + jobTree.classPrefix + `icon {
	  left: 56px;
	}
	
	.` + jobTree.classPrefix + `container .` + jobTree.classPrefix + `content {
	  color: ` + jobTree.Color + `;
	  padding: 30px 90px 30px 30px;
	  background: ` + jobTree.ItemBackground + `;
	  position: relative;
	  border-radius: 0 500px 500px 0;
	}
	
	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right .` + jobTree.classPrefix + `content {
	  padding: 30px 30px 30px 90px;
	  border-radius: 500px 0 0 500px;
	  text-align: right;
	}
	
	.` + jobTree.classPrefix + `container .` + jobTree.classPrefix + `content h1 {
	  	margin: 0 0 10px 0;
		font-size: 20px;
		font-weight: bold;
		color: ` + jobTree.TitleColor + `;
	}

	.` + jobTree.classPrefix + `container .` + jobTree.classPrefix + `content h2 {
	  margin: 0 0 10px 0;
	  font-size: 18px;
	  font-weight: normal;
	  color: ` + jobTree.TitleColor + `;
	}
	
	.` + jobTree.classPrefix + `container .` + jobTree.classPrefix + `content p {
	  margin: 0;
	  font-size: 16px;
	  line-height: 22px;
	  color: #000000;
	}

	.` + jobTree.classPrefix + "content-tag-item" + `{
		display: inline-block;
		padding: 2px 5px;
		margin: 0 5px 5px 0;
		border-radius: 5px;
		color: ` + jobTree.TagColor + `;
	}
	
	@media (max-width: 767.98px) {
	  	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `left .` + jobTree.classPrefix + `date {
			right: -100%;
	  	}
	  	.` + jobTree.classPrefix + `container.` + jobTree.classPrefix + `right .` + jobTree.classPrefix + `date {
			left: -100%;
		}
	}
	`

	timeline.StyleBlock(css)
	return timeline
}
