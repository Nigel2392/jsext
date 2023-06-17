//go:build js && wasm
// +build js,wasm

package elements

import (
	"regexp"
	"strconv"
	"strings"
)

func A(href string, text ...string) *Element {
	return NewElement("a", text...).Set("href", href)
}

func Abbr(title string, text ...string) *Element {
	return NewElement("abbr", text...).Set("title", title)
}

func Address(text ...string) *Element {
	return NewElement("address", text...)
}

func Area(alt string, coords ...string) *Element {
	return NewElement("area", alt).Set("coords", coords...)
}

func Article(text ...string) *Element {
	return NewElement("article", text...)
}

func Aside(text ...string) *Element {
	return NewElement("aside", text...)
}

func Audio(src string) *Element {
	return NewElement("audio").Set("src", src)
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

func Button(text ...string) *Element {
	return NewElement("button", text...)
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
	return NewElement("div").AttrClass(classList...)
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

func Form(action, method string) *Element {
	return NewElement("form").Set("action", action).Set("method", method)
}

func Heading(size int, text ...string) *Element {
	var s string = strconv.Itoa(size)
	if size < 1 || size > 6 {
		s = "1"
	}
	var b = make([]byte, 1+len(s))
	b[0] = 'h'
	b = append(b, s...)
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
	return NewElement("iframe").Set("src", src)
}

func Img(src string) *Element {
	return NewElement("img").Set("src", src)
}

func Input(typ, name string) *Element {
	return NewElement("input").Set("type", typ).Set("name", name)
}

func Ins(text ...string) *Element {
	return NewElement("ins", text...)
}

func Kbd(text ...string) *Element {
	return NewElement("kbd", text...)
}

func Label(text, forElem string) *Element {
	return NewElement("label", text).Set("for", forElem)
}

func Legend(text ...string) *Element {
	return NewElement("legend", text...)
}

func Li(text ...string) *Element {
	return NewElement("li", text...)
}

func Link(href string) *Element {
	return NewElement("link").Set("href", href)
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
	return NewElement("meta").Set("name", name).Set("content", content)
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

func Option(text, value string, selected ...bool) *Element {
	var e = NewElement("option", text)
	e.Set("value", value)
	if len(selected) > 0 && selected[0] {
		e.Set("selected", "selected")
	}
	return e
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
	return NewElement("script").Set("src", src)
}

func Section(text ...string) *Element {
	return NewElement("section", text...)
}

func Select(text ...string) *Element {
	return NewElement("select", text...)
}

func Small(text ...string) *Element {
	return NewElement("small", text...)
}

func Source(src string) *Element {
	return NewElement("source").Set("src", src)
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

func StyleBlock(t ...string) *Element {
	var v = NewElement("style")
	v.Set("type", "text/css")
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

func Table(text ...string) *Element {
	return NewElement("table", text...)
}

func Tbody(text ...string) *Element {
	return NewElement("tbody", text...)
}

func Td(text ...string) *Element {
	return NewElement("td", text...)
}

func Textarea(name, placeholder string) *Element {
	return NewElement("textarea").Set("name", name).Set("placeholder", placeholder)
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
