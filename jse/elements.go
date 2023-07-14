//go:build js && wasm
// +build js,wasm

package jse

import (
	"regexp"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/Nigel2392/jsext/v2"
)

func A(href string, text ...string) *Element {
	return NewElement("a", text...).SetAttr("href", href)
}

func Abbr(title string, text ...string) *Element {
	return NewElement("abbr", text...).SetAttr("title", title)
}

func Address(text ...string) *Element {
	return NewElement("address", text...)
}

func Area(alt string, coords ...string) *Element {
	return NewElement("area", alt).SetAttr("coords", coords...)
}

func Article(text ...string) *Element {
	return NewElement("article", text...)
}

func Aside(text ...string) *Element {
	return NewElement("aside", text...)
}

func Audio(src string) *Element {
	return NewElement("audio").SetAttr("src", src)
}

func B(text ...string) *Element {
	return NewElement("b", text...)
}

func Bdi(text ...string) *Element {
	return NewElement("bdi", text...)
}

func Bdo(text ...string) *Element {
	return NewElement("bdo", text...)
}

func Blockquote(text ...string) *Element {
	return NewElement("blockquote", text...)
}

func Body(text ...string) *Element {
	return NewElement("body", text...)
}

func Br() *Element {
	return NewElement("br")
}

func Button(innerText string, onClick func(*Element, jsext.Event)) *Element {
	var button = NewElement("button", innerText)
	button.AddEventListener("click", onClick)
	return button
}

func Canvas(text ...string) *Element {
	return NewElement("canvas", text...)
}

func Caption(text ...string) *Element {
	return NewElement("caption", text...)
}

func Cite(text ...string) *Element {
	return NewElement("cite", text...)
}

func Code(text ...string) *Element {
	return NewElement("code", text...)
}

func Col() *Element {
	return NewElement("col")
}

func Colgroup(text ...string) *Element {
	return NewElement("colgroup", text...)
}

func Command(text ...string) *Element {
	return NewElement("command", text...)
}

func Datalist(text ...string) *Element {
	return NewElement("datalist", text...)
}

func Dd(text ...string) *Element {
	return NewElement("dd", text...)
}

func Del(text ...string) *Element {
	return NewElement("del", text...)
}

func Details(text ...string) *Element {
	return NewElement("details", text...)
}

func Dfn(text ...string) *Element {
	return NewElement("dfn", text...)
}

func Div(classList ...string) *Element {
	var e = NewElement("div")
	if len(classList) > 0 {
		e.ClassList(classList...)
	}
	return e
}

func Dl(text ...string) *Element {
	return NewElement("dl", text...)
}

func Dt(text ...string) *Element {
	return NewElement("dt", text...)
}

func Em(text ...string) *Element {
	return NewElement("em", text...)
}

func Fieldset(text ...string) *Element {
	return NewElement("fieldset", text...)
}

func Figcaption(text ...string) *Element {
	return NewElement("figcaption", text...)
}

func Figure(text ...string) *Element {
	return NewElement("figure", text...)
}

func Footer(text ...string) *Element {
	return NewElement("footer", text...)
}

func Heading(size int, text ...string) *Element {
	var s string = strconv.Itoa(size)
	if size < 1 || size > 6 {
		s = "1"
	}
	var b = make([]byte, 2)
	b[0] = 'h'
	b[1] = s[0]
	return NewElement(string(b), text...)
}

func Header(text ...string) *Element {
	return NewElement("header", text...)
}

func Hr() *Element {
	return NewElement("hr")
}

func Html(text ...string) *Element {
	return NewElement("html", text...)
}

func I(text ...string) *Element {
	return NewElement("i", text...)
}

func Iframe(src string) *Element {
	return NewElement("iframe").SetAttr("src", src)
}

func Img(src string) *Element {
	return NewElement("img").SetAttr("src", src)
}

func Label(forElement, text string, classes ...string) *Element {
	var e = jsext.CreateElement("label")
	if len(classes) > 0 {
		e.ClassList(classes...)
	}
	e.SetAttribute("for", forElement)
	e.InnerText(text)
	return (*Element)(&e)
}

func Input(t, name string, opts *InputOptions) *Element {
	var e = jsext.CreateElement("input")
	e.SetAttribute("type", t)
	e.SetAttribute("name", name)
	e.SetAttribute("id", name)
	opts.Apply(e)
	return (*Element)(&e)
}

type sub_add int8

const (
	sub sub_add = iota
	add
)

func timeAdder(input *Element, typ sub_add) func(this *Element, event jsext.Event) {
	return func(this *Element, event jsext.Event) {
		event.PreventDefault()
		var value = input.Get("value").String()
		var intVal, err = strconv.Atoi(value)
		if err != nil {
			return
		}
		switch typ {
		case sub:
			intVal--
		case add:
			intVal++
		}
		if intVal < 0 {
			intVal = 0
		}
		input.Set("value", strconv.Itoa(intVal))
	}
}

type Timer struct {
	*Element
	hoursInput   *Element
	minutesInput *Element
	secondsInput *Element
}

func (t *Timer) Time() time.Time {
	var (
		hours   = t.hoursInput.Get("value").String()
		minutes = t.minutesInput.Get("value").String()
		seconds = t.secondsInput.Get("value").String()
	)
	var (
		h, _ = strconv.Atoi(hours)
		m, _ = strconv.Atoi(minutes)
		s, _ = strconv.Atoi(seconds)
	)
	return time.Date(0, 0, 0, h, m, s, 0, time.UTC)
}

func (t *Timer) Duration() time.Duration {
	var (
		hours   = t.hoursInput.Get("value").String()
		minutes = t.minutesInput.Get("value").String()
		seconds = t.secondsInput.Get("value").String()
	)
	var (
		h, _ = strconv.Atoi(hours)
		m, _ = strconv.Atoi(minutes)
		s, _ = strconv.Atoi(seconds)
	)

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
}

func (t *Timer) SetHours(hours int) {
	t.hoursInput.Set("value", strconv.Itoa(hours))
}

func (t *Timer) SetMinutes(minutes int) {
	t.minutesInput.Set("value", strconv.Itoa(minutes))
}

func (t *Timer) SetSeconds(seconds int) {
	t.secondsInput.Set("value", strconv.Itoa(seconds))
}

func (t *Timer) SetTime(time time.Time) {
	t.hoursInput.Set("value", strconv.Itoa(time.Hour()))
	t.minutesInput.Set("value", strconv.Itoa(time.Minute()))
	t.secondsInput.Set("value", strconv.Itoa(time.Second()))
}

func (t *Timer) SetDuration(duration time.Duration) {
	t.hoursInput.Set("value", strconv.Itoa(int(duration.Hours())))
	t.minutesInput.Set("value", strconv.Itoa(int(duration.Minutes())))
	t.secondsInput.Set("value", strconv.Itoa(int(duration.Seconds())))
}

func (t *Timer) Value() *Element {
	return t.Element
}

func checkTimeCss() {
	var head = js.Global().Get("document").Get("head")
	var style = head.Call("querySelector", "style#time-input-css")
	if style.IsUndefined() {
		style = StyleBlock(`
			.time-input-container{
				display: flex;
				flex-direction: row;
				justify-content: space-between;
				align-items: center;
				margin: 0.5rem;
				width: 100%;
			}
			.time-input-group {
				display: flex;
				flex-direction: column;
				justify-content: space-between;
				align-items: center;
			}
			.time-input-container > .time-input-group {
				flex: 1 1 auto;
				margin: 0.5rem;
			}
			.time-input-container > .time-input-group:first-child {
				margin-left: 0;
			}
			.time-input-container > .time-input-group:last-child {
				margin-right: 0;
			}
			.time-input-group > .time-input-items {
				display: flex;
				flex-direction: row;
				justify-content: space-between;
				align-items: stretch;
				margin: 0.5rem;
				width: 100%;
			}
			.time-input-items > .time-input{
				width: 50%;
				border-radius: 0.5rem 0 0 0.5rem !important;
				padding: 0.5rem;
			}
			.time-input-items > .time-input-label {
				width: 50%;
				border-radius: 0 0.5rem 0.5rem 0 !important;
				padding: 0.5rem;
				font-weight: bold;
			}
		`).JSValue()
		style.Set("id", "time-input-css")
		head.Call("appendChild", style)
	}
}

func TimeInput(name string, h, m, s int, opts *InputOptions) *Timer {

	checkTimeCss()

	var inputDiv = Div("time-input-container")

	var (
		inputGroupHour = inputDiv.Div("time-input-group")
		inputHour      = Input("text", name, opts).InlineClasses("time-input")
		inputHourLabel = P("H").InlineClasses("time-input-label")
		hourAdd        = Button("+", nil).OnHoldClick(time.Millisecond*500, time.Millisecond*50, timeAdder(inputHour, add))
		hourSub        = Button("-", nil).OnHoldClick(time.Millisecond*500, time.Millisecond*50, timeAdder(inputHour, sub))
	)
	inputGroupHour.AppendChild(
		hourAdd,
		Div("time-input-items").AppendChild(
			inputHour,
			inputHourLabel,
		),
		hourSub,
	)

	var (
		inputGroupMinute = inputDiv.Div("time-input-group")
		inputMinute      = Input("text", name, opts).InlineClasses("time-input")
		inputMinuteLabel = P("M").InlineClasses("time-input-label")
		minuteAdd        = Button("+", nil).OnHoldClick(time.Millisecond*500, time.Millisecond*50, timeAdder(inputMinute, add))
		minuteSub        = Button("-", nil).OnHoldClick(time.Millisecond*500, time.Millisecond*50, timeAdder(inputMinute, sub))
	)
	inputGroupMinute.AppendChild(
		minuteAdd,
		Div("time-input-items").AppendChild(
			inputMinute,
			inputMinuteLabel,
		),
		minuteSub,
	)

	var (
		inputGroupSecond = inputDiv.Div("time-input-group")
		inputSecond      = Input("text", name, opts).InlineClasses("time-input")
		inputSecondLabel = P("S").InlineClasses("time-input-label")
		secondAdd        = Button("+", nil).OnHoldClick(time.Millisecond*500, time.Millisecond*50, timeAdder(inputSecond, add))
		secondSub        = Button("-", nil).OnHoldClick(time.Millisecond*500, time.Millisecond*50, timeAdder(inputSecond, sub))
	)
	inputGroupSecond.AppendChild(
		secondAdd,
		Div("time-input-items").AppendChild(
			inputSecond,
			inputSecondLabel,
		),
		secondSub,
	)

	hourAdd.InlineClasses("time-input-add", "time-input-button").SetAttr("type", "button")
	minuteAdd.InlineClasses("time-input-add", "time-input-button").SetAttr("type", "button")
	secondAdd.InlineClasses("time-input-add", "time-input-button").SetAttr("type", "button")

	hourSub.InlineClasses("time-input-sub", "time-input-button").SetAttr("type", "button")
	minuteSub.InlineClasses("time-input-sub", "time-input-button").SetAttr("type", "button")
	secondSub.InlineClasses("time-input-sub", "time-input-button").SetAttr("type", "button")

	if len(opts.ButtonClasses) > 0 {
		hourAdd.ClassList(opts.ButtonClasses...)
		hourSub.ClassList(opts.ButtonClasses...)
		minuteAdd.ClassList(opts.ButtonClasses...)
		minuteSub.ClassList(opts.ButtonClasses...)
		secondAdd.ClassList(opts.ButtonClasses...)
		secondSub.ClassList(opts.ButtonClasses...)
	}

	var t = &Timer{
		Element:      inputDiv,
		hoursInput:   inputHour,
		minutesInput: inputMinute,
		secondsInput: inputSecond,
	}

	t.SetHours(h)
	t.SetMinutes(m)
	t.SetSeconds(s)

	return t
}

func TextArea(name string, opts *InputOptions) *Element {
	var e = jsext.CreateElement("textarea")
	e.SetAttribute("name", name)
	e.SetAttribute("id", name)
	opts.Apply(e)
	return (*Element)(&e)
}

func Ins(text ...string) *Element {
	return NewElement("ins", text...)
}

func Kbd(text ...string) *Element {
	return NewElement("kbd", text...)
}

func Legend(text ...string) *Element {
	return NewElement("legend", text...)
}

func Li(text ...string) *Element {
	return NewElement("li", text...)
}

func Link(href string) *Element {
	return NewElement("link").SetAttr("href", href)
}

func Main(text ...string) *Element {
	return NewElement("main", text...)
}

func Map(text ...string) *Element {
	return NewElement("map", text...)
}

func Mark(text ...string) *Element {
	return NewElement("mark", text...)
}

func Meta(name, content string) *Element {
	return NewElement("meta").SetAttr("name", name).SetAttr("content", content)
}

func Meter(text ...string) *Element {
	return NewElement("meter", text...)
}

func Nav(text ...string) *Element {
	return NewElement("nav", text...)
}

func Noscript(text ...string) *Element {
	return NewElement("noscript", text...)
}

func Object(text ...string) *Element {
	return NewElement("object", text...)
}

func Ol(text ...string) *Element {
	return NewElement("ol", text...)
}

func Optgroup(text ...string) *Element {
	return NewElement("optgroup", text...)
}

func Output(text ...string) *Element {
	return NewElement("output", text...)
}

func P(text ...string) *Element {
	return NewElement("p", text...)
}

func Param(text ...string) *Element {
	return NewElement("param", text...)
}

func Pre(text ...string) *Element {
	return NewElement("pre", text...)
}

func Progress(text ...string) *Element {
	return NewElement("progress", text...)
}

func Q(text ...string) *Element {
	return NewElement("q", text...)
}

func Rp(text ...string) *Element {
	return NewElement("rp", text...)
}

func Rt(text ...string) *Element {
	return NewElement("rt", text...)
}

func Ruby(text ...string) *Element {
	return NewElement("ruby", text...)
}

func Samp(text ...string) *Element {
	return NewElement("samp", text...)
}

func Script(src string) *Element {
	return NewElement("script").SetAttr("src", src)
}

func Section(text ...string) *Element {
	return NewElement("section", text...)
}

func Small(text ...string) *Element {
	return NewElement("small", text...)
}

func Source(src string) *Element {
	return NewElement("source").SetAttr("src", src)
}

func Span(text ...string) *Element {
	return NewElement("span", text...)
}

func Strong(text ...string) *Element {
	return NewElement("strong", text...)
}

func Style(text ...string) *Element {
	return NewElement("style", text...)
}

// Create a style element with the given text as its content.
func StyleBlock(t ...string) *Element {
	var v = NewElement("style")
	v.SetAttr("type", "text/css")
	if len(t) > 0 {
		var sourceCode = strings.Join(t, "\n")
		sourceCode = strings.ReplaceAll(sourceCode, "\n\n", "\n")
		var re = regexp.MustCompile(`\s+([a-z-]+ *:[^:;]+;)\n`)
		sourceCode = re.ReplaceAllString(sourceCode, "$1")
		v.InnerHTML(sourceCode)
	}
	return v
}

func Sub(text ...string) *Element {
	return NewElement("sub", text...)
}

func Summary(text ...string) *Element {
	return NewElement("summary", text...)
}

func Sup(text ...string) *Element {
	return NewElement("sup", text...)
}

func Table(classList ...string) *Element {
	var e = NewElement("table")
	if len(classList) > 0 {
		e.ClassList(classList...)
	}
	return e
}

func Tbody(text ...string) *Element {
	return NewElement("tbody", text...)
}

func Td(text ...string) *Element {
	return NewElement("td", text...)
}

func Textarea(name, placeholder string) *Element {
	return NewElement("textarea").SetAttr("name", name).SetAttr("placeholder", placeholder)
}

func Tfoot(text ...string) *Element {
	return NewElement("tfoot", text...)
}

func Th(text ...string) *Element {
	return NewElement("th", text...)
}

func Thead(text ...string) *Element {
	return NewElement("thead", text...)
}

func Time(text ...string) *Element {
	return NewElement("time", text...)
}

func Title(text ...string) *Element {
	return NewElement("title", text...)
}

func Tr(text ...string) *Element {
	return NewElement("tr", text...)
}

func Track(text ...string) *Element {
	return NewElement("track", text...)
}

func Ul(text ...string) *Element {
	return NewElement("ul", text...)
}

func Var(text ...string) *Element {
	return NewElement("var", text...)
}

func Video(text ...string) *Element {
	return NewElement("video", text...)
}

func Wbr(text ...string) *Element {
	return NewElement("wbr", text...)
}
