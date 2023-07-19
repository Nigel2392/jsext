package dom

import (
	"strings"
	"syscall/js"
)

type Node struct {
	js.Value
	Depth int
}

func (d Node) Parent() Node {
	return Node{d.Get("parentNode"), d.Depth - 1}
}

func (d Node) Children() []Node {
	var children = d.Get("childNodes")
	var nodes = make([]Node, children.Length())
	for i := 0; i < children.Length(); i++ {
		nodes[i] = Node{children.Index(i), d.Depth + 1}
	}
	return nodes
}

func (d Node) NextSibling() Node {
	return Node{d.Get("nextSibling"), d.Depth}
}

func (d Node) PreviousSibling() Node {
	return Node{d.Get("previousSibling"), d.Depth}
}

func (d Node) FirstChild() Node {
	return Node{d.Get("firstChild"), d.Depth + 1}
}

func (d Node) LastChild() Node {
	return Node{d.Get("lastChild"), d.Depth + 1}
}

func (d Node) NodeType() NodeType {
	return NodeType(d.Get("nodeType").Int())
}

func (d Node) NodeName() string {
	return d.Get("nodeName").String()
}

func (d Node) Get(key string) js.Value {
	var keys = strings.Split(key, ".")
	var v = d.Value
	for _, k := range keys {
		v = v.Get(k)
	}
	return v
}

func (d Node) Set(key string, value any) {
	var keys = strings.Split(key, ".")
	var v = d.Value
	for i := 0; i < len(keys)-1; i++ {
		v = v.Get(keys[i])
	}
	v.Set(keys[len(keys)-1], value)
}

func (d Node) Call(key string, args ...any) js.Value {
	var keys = strings.Split(key, ".")
	var v = d.Value
	for _, k := range keys {
		v = v.Get(k)
	}
	return v.Invoke(args...)
}

func (d Node) Walk(nodeTypes []NodeType, fn func(Node)) {
	Walk(nodeTypes, d.Value, fn)
}
