package jse

import (
	"fmt"

	"github.com/Nigel2392/jsext/v2"
)

type SVG jsext.Element

func NewSVG(xmlns string) *SVG {
	var e = jsext.Document.Call("createElementNS", xmlns, "svg")
	e.Call("setAttribute", "xmlns", xmlns)
	return (*SVG)(&e)
}

func (s *SVG) Element() *Element {
	return (*Element)(s)
}

func (s *SVG) XMLNS(u string) *SVG {
	s.Element().SetAttr("xmlns", u)
	return s
}

func (s *SVG) ViewBox(x, y, width, height int) *SVG {
	s.Element().Call("setAttributeNS", nil, "viewBox", fmt.Sprintf("%d %d %d %d", x, y, width, height))
	return s
}

func (s *SVG) Fill(color string) *SVG {
	s.Element().SetAttr("fill", color)
	return s
}

func (s *SVG) Width(width int) *SVG {
	s.Element().SetAttr("width", fmt.Sprintf("%d", width))
	return s
}

func (s *SVG) Height(height int) *SVG {
	s.Element().SetAttr("height", fmt.Sprintf("%d", height))
	return s
}

func (s *SVG) CreateElement(name string) *Element {
	var ns = s.Element().Get("namespaceURI")
	var e = jsext.Document.Call("createElementNS", ns.Value(), name)
	s.Element().Call("appendChild", e)
	return (*Element)(&e)
}
