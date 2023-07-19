package dom

import (
	"syscall/js"
)

func Walk(nodetypes []NodeType, e js.Value, fn func(Node)) {
	walk(nodetypes, e, fn, 0)
}

func walk(nodetypes []NodeType, e js.Value, fn func(Node), depth int) {
	fn(Node{e, depth})
	var childNodes = e.Get("childNodes")
	for i := 0; i < childNodes.Length(); i++ {
		var child = childNodes.Index(i)
	inner:
		for _, nodetype := range nodetypes {
			if child.Get("nodeType").Int() == int(nodetype) {
				walk(nodetypes, child, fn, depth+1)
				break inner
			}
		}
	}
}

var domParser = js.Global().Get("DOMParser").New()

func Parse(html string) Document {
	var doc = domParser.Call("parseFromString", html, "text/html")
	return Document(doc)
}

type Document js.Value

func (d Document) Head() js.Value {
	return (js.Value)(d).Get("head")
}

func (d Document) Body() js.Value {
	return (js.Value)(d).Get("body")
}

func (d Document) QuerySelector(query string) js.Value {
	return (js.Value)(d).Call("querySelector", query)
}

func (d Document) QuerySelectorAll(query string) []js.Value {
	var e = (js.Value)(d).Call("querySelectorAll", query)
	var elements = make([]js.Value, e.Length())
	for i := 0; i < e.Length(); i++ {
		elements[i] = e.Index(i)
	}
	return elements
}

func (d Document) GetElementById(id string) js.Value {
	return (js.Value)(d).Call("getElementById", id)
}

func (d Document) Walk(nodeTypes []NodeType, fn func(Node)) {
	Walk(nodeTypes, (js.Value)(d), fn)
}
