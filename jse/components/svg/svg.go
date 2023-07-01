package svg

import (
	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
)

func BarChartSteps(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-bar-chart-steps")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M.5 0a.5.5 0 0 1 .5.5v15a.5.5 0 0 1-1 0V.5A.5.5 0 0 1 .5 0zM2 1.5a.5.5 0 0 1 .5-.5h4a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-4a.5.5 0 0 1-.5-.5v-1zm2 4a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-7a.5.5 0 0 1-.5-.5v-1zm2 4a.5.5 0 0 1 .5-.5h6a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-6a.5.5 0 0 1-.5-.5v-1zm2 4a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-7a.5.5 0 0 1-.5-.5v-1z")
	return svg
}

func Trashcan(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-trash")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z")
	path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z")
	return svg
}

func User(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-person")
	svg.ViewBox(0, 0, 24, 24)
	var g = svg.NewElement("g")

	var path = jsext.Document.Call("createElementNS", "http://www.w3.org/2000/svg", "path")
	path.Call("setAttribute", "fill", "none")
	path.Call("setAttribute", "d", "M0 0h16v16H0z")
	g.AppendChild((*jse.Element)(&path))

	path = jsext.Document.Call("createElementNS", "http://www.w3.org/2000/svg", "path")
	path.Call("setAttribute", "d", "M20 22h-2v-2a3 3 0 0 0-3-3H9a3 3 0 0 0-3 3v2H4v-2a5 5 0 0 1 5-5h6a5 5 0 0 1 5 5v2zm-8-9a6 6 0 1 1 0-12 6 6 0 0 1 0 12zm0-2a4 4 0 1 0 0-8 4 4 0 0 0 0 8z")
	g.AppendChild((*jse.Element)(&path))

	return svg
}

func ArrowDownUp(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-arrow-down-up")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M11.5 15a.5.5 0 0 0 .5-.5V2.707l3.146 3.147a.5.5 0 0 0 .708-.708l-4-4a.5.5 0 0 0-.708 0l-4 4a.5.5 0 1 0 .708.708L11 2.707V14.5a.5.5 0 0 0 .5.5zm-7-14a.5.5 0 0 1 .5.5v11.793l3.146-3.147a.5.5 0 0 1 .708.708l-4 4a.5.5 0 0 1-.708 0l-4-4a.5.5 0 0 1 .708-.708L4 13.293V1.5a.5.5 0 0 1 .5-.5z")
	return svg
}

func MergeHorizontal(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-arrow-down-up")
	svg.ViewBox(0, 0, 24, 24)
	var path = svg.NewElement("path")
	path.SetAttr("fill", "currentColor")
	path.SetAttr("d", "M12 8.9757L16.2426 4.73306L14.8284 3.31885L12 6.14727L9.17157 3.31885L7.75736 4.73306L12 8.9757Z")
	path = svg.NewElement("path")
	path.SetAttr("fill", "currentColor")
	path.SetAttr("d", "M5 11.9998C5 11.4475 5.44772 10.9997 6 10.9997H18C18.5523 10.9997 19 11.4475 19 11.9998C19 12.552 18.5523 12.9998 18 12.9998H6C5.44772 12.9998 5 12.552 5 11.9998Z")
	path = svg.NewElement("path")
	path.SetAttr("fill", "currentColor")
	path.SetAttr("d", "M12 15.0242L7.75736 19.2668L9.17157 20.681L12 17.8526L14.8284 20.681L16.2426 19.2668L12 15.0242Z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-type-bold" viewBox="0 0 16 16">
//
//	<path d="M8.21 13c2.106 0 3.412-1.087 3.412-2.823 0-1.306-.984-2.283-2.324-2.386v-.055a2.176 2.176 0 0 0 1.852-2.14c0-1.51-1.162-2.46-3.014-2.46H3.843V13H8.21zM5.908 4.674h1.696c.963 0 1.517.451 1.517 1.244 0 .834-.629 1.32-1.73 1.32H5.908V4.673zm0 6.788V8.598h1.73c1.217 0 1.88.492 1.88 1.415 0 .943-.643 1.449-1.832 1.449H5.907z"/>
//
// </svg>
func Bold(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-type-bold")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M8.21 13c2.106 0 3.412-1.087 3.412-2.823 0-1.306-.984-2.283-2.324-2.386v-.055a2.176 2.176 0 0 0 1.852-2.14c0-1.51-1.162-2.46-3.014-2.46H3.843V13H8.21zM5.908 4.674h1.696c.963 0 1.517.451 1.517 1.244 0 .834-.629 1.32-1.73 1.32H5.908V4.673zm0 6.788V8.598h1.73c1.217 0 1.88.492 1.88 1.415 0 .943-.643 1.449-1.832 1.449H5.907z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-type-italic" viewBox="0 0 16 16">
//
//	<path d="M7.991 11.674 9.53 4.455c.123-.595.246-.71 1.347-.807l.11-.52H7.211l-.11.52c1.06.096 1.128.212 1.005.807L6.57 11.674c-.123.595-.246.71-1.346.806l-.11.52h3.774l.11-.52c-1.06-.095-1.129-.211-1.006-.806z"/>
//
// </svg>
func Italic(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-type-italic")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M7.991 11.674 9.53 4.455c.123-.595.246-.71 1.347-.807l.11-.52H7.211l-.11.52c1.06.096 1.128.212 1.005.807L6.57 11.674c-.123.595-.246.71-1.346.806l-.11.52h3.774l.11-.52c-1.06-.095-1.129-.211-1.006-.806z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-type-underline" viewBox="0 0 16 16">
//
//	<path d="M5.313 3.136h-1.23V9.54c0 2.105 1.47 3.623 3.917 3.623s3.917-1.518 3.917-3.623V3.136h-1.23v6.323c0 1.49-.978 2.57-2.687 2.57-1.709 0-2.687-1.08-2.687-2.57V3.136zM12.5 15h-9v-1h9v1z"/>
//
// </svg>
func Underline(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-type-underline")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M5.313 3.136h-1.23V9.54c0 2.105 1.47 3.623 3.917 3.623s3.917-1.518 3.917-3.623V3.136h-1.23v6.323c0 1.49-.978 2.57-2.687 2.57-1.709 0-2.687-1.08-2.687-2.57V3.136zM12.5 15h-9v-1h9v1z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-type-strikethrough" viewBox="0 0 16 16">
//
//	<path d="M6.333 5.686c0 .31.083.581.27.814H5.166a2.776 2.776 0 0 1-.099-.76c0-1.627 1.436-2.768 3.48-2.768 1.969 0 3.39 1.175 3.445 2.85h-1.23c-.11-1.08-.964-1.743-2.25-1.743-1.23 0-2.18.602-2.18 1.607zm2.194 7.478c-2.153 0-3.589-1.107-3.705-2.81h1.23c.144 1.06 1.129 1.703 2.544 1.703 1.34 0 2.31-.705 2.31-1.675 0-.827-.547-1.374-1.914-1.675L8.046 8.5H1v-1h14v1h-3.504c.468.437.675.994.675 1.697 0 1.826-1.436 2.967-3.644 2.967z"/>
//
// </svg>
func StrikeThrough(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-type-strikethrough")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M6.333 5.686c0 .31.083.581.27.814H5.166a2.776 2.776 0 0 1-.099-.76c0-1.627 1.436-2.768 3.48-2.768 1.969 0 3.39 1.175 3.445 2.85h-1.23c-.11-1.08-.964-1.743-2.25-1.743-1.23 0-2.18.602-2.18 1.607zm2.194 7.478c-2.153 0-3.589-1.107-3.705-2.81h1.23c.144 1.06 1.129 1.703 2.544 1.703 1.34 0 2.31-.705 2.31-1.675 0-.827-.547-1.374-1.914-1.675L8.046 8.5H1v-1h14v1h-3.504c.468.437.675.994.675 1.697 0 1.826-1.436 2.967-3.644 2.967z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-list-ul" viewBox="0 0 16 16">
//
//	<path fill-rule="evenodd" d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm-3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"/>
//
// </svg>
func ListUL(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-list-ul")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm-3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm0 4a1 1 0 1 0 0-2 1 1 0 0 0 0 2z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-list-ol" viewBox="0 0 16 16">
//
//	<path fill-rule="evenodd" d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5z"/>
//	<path d="M1.713 11.865v-.474H2c.217 0 .363-.137.363-.317 0-.185-.158-.31-.361-.31-.223 0-.367.152-.373.31h-.59c.016-.467.373-.787.986-.787.588-.002.954.291.957.703a.595.595 0 0 1-.492.594v.033a.615.615 0 0 1 .569.631c.003.533-.502.8-1.051.8-.656 0-1-.37-1.008-.794h.582c.008.178.186.306.422.309.254 0 .424-.145.422-.35-.002-.195-.155-.348-.414-.348h-.3zm-.004-4.699h-.604v-.035c0-.408.295-.844.958-.844.583 0 .96.326.96.756 0 .389-.257.617-.476.848l-.537.572v.03h1.054V9H1.143v-.395l.957-.99c.138-.142.293-.304.293-.508 0-.18-.147-.32-.342-.32a.33.33 0 0 0-.342.338v.041zM2.564 5h-.635V2.924h-.031l-.598.42v-.567l.629-.443h.635V5z"/>
//
// </svg>
func ListOL(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-list-ol")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zm0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5z")
	path = svg.NewElement("path")
	path.SetAttr("d", "M1.713 11.865v-.474H2c.217 0 .363-.137.363-.317 0-.185-.158-.31-.361-.31-.223 0-.367.152-.373.31h-.59c.016-.467.373-.787.986-.787.588-.002.954.291.957.703a.595.595 0 0 1-.492.594v.033a.615.615 0 0 1 .569.631c.003.533-.502.8-1.051.8-.656 0-1-.37-1.008-.794h.582c.008.178.186.306.422.309.254 0 .424-.145.422-.35-.002-.195-.155-.348-.414-.348h-.3zm-.004-4.699h-.604v-.035c0-.408.295-.844.958-.844.583 0 .96.326.96.756 0 .389-.257.617-.476.848l-.537.572v.03h1.054V9H1.143v-.395l.957-.99c.138-.142.293-.304.293-.508 0-.18-.147-.32-.342-.32a.33.33 0 0 0-.342.338v.041zM2.564 5h-.635V2.924h-.031l-.598.42v-.567l.629-.443h.635V5z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-link" viewBox="0 0 16 16">
//
//	<path d="M6.354 5.5H4a3 3 0 0 0 0 6h3a3 3 0 0 0 2.83-4H9c-.086 0-.17.01-.25.031A2 2 0 0 1 7 10.5H4a2 2 0 1 1 0-4h1.535c.218-.376.495-.714.82-1z"/>
//	<path d="M9 5.5a3 3 0 0 0-2.83 4h1.098A2 2 0 0 1 9 6.5h3a2 2 0 1 1 0 4h-1.535a4.02 4.02 0 0 1-.82 1H12a3 3 0 1 0 0-6H9z"/>
//
// </svg>
func Link(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-link")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M6.354 5.5H4a3 3 0 0 0 0 6h3a3 3 0 0 0 2.83-4H9c-.086 0-.17.01-.25.031A2 2 0 0 1 7 10.5H4a2 2 0 1 1 0-4h1.535c.218-.376.495-.714.82-1z")
	path = svg.NewElement("path")
	path.SetAttr("d", "M9 5.5a3 3 0 0 0-2.83 4h1.098A2 2 0 0 1 9 6.5h3a2 2 0 1 1 0 4h-1.535a4.02 4.02 0 0 1-.82 1H12a3 3 0 1 0 0-6H9z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-text-left" viewBox="0 0 16 16">
//
//	<path fill-rule="evenodd" d="M2 12.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm0-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm0-3a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm0-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z"/>
//
// </svg>
func TextLeft(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-text-left")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M2 12.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm0-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm0-3a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm0-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-text-center" viewBox="0 0 16 16">
//
//	<path fill-rule="evenodd" d="M4 12.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-2-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm2-3a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-2-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z"/>
//
// </svg>
func TextCenter(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-text-center")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M4 12.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-2-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm2-3a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-2-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-text-right" viewBox="0 0 16 16">
//
//	<path fill-rule="evenodd" d="M6 12.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-4-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm4-3a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-4-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z"/>
//
// </svg>
func TextRight(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-text-right")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("fill-rule", "evenodd")
	path.SetAttr("d", "M6 12.5a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-4-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5zm4-3a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1-.5-.5zm-4-3a.5.5 0 0 1 .5-.5h11a.5.5 0 0 1 0 1h-11a.5.5 0 0 1-.5-.5z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-file-image" viewBox="0 0 16 16">
//
//	<path d="M8.002 5.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z"/>
//	<path d="M12 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2zM3 2a1 1 0 0 1 1-1h8a1 1 0 0 1 1 1v8l-2.083-2.083a.5.5 0 0 0-.76.063L8 11 5.835 9.7a.5.5 0 0 0-.611.076L3 12V2z"/>
//
// </svg>
func FileImage(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-file-image")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M8.002 5.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0z")
	path = svg.NewElement("path")
	path.SetAttr("d", "M12 0H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2zM3 2a1 1 0 0 1 1-1h8a1 1 0 0 1 1 1v8l-2.083-2.083a.5.5 0 0 0-.76.063L8 11 5.835 9.7a.5.5 0 0 0-.611.076L3 12V2z")
	return svg
}

// <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-textarea-resize" viewBox="0 0 16 16">
//
//	<path d="M0 4.5A2.5 2.5 0 0 1 2.5 2h11A2.5 2.5 0 0 1 16 4.5v7a2.5 2.5 0 0 1-2.5 2.5h-11A2.5 2.5 0 0 1 0 11.5v-7zM2.5 3A1.5 1.5 0 0 0 1 4.5v7A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5v-7A1.5 1.5 0 0 0 13.5 3h-11zm10.854 4.646a.5.5 0 0 1 0 .708l-3 3a.5.5 0 0 1-.708-.708l3-3a.5.5 0 0 1 .708 0zm0 2.5a.5.5 0 0 1 0 .708l-.5.5a.5.5 0 0 1-.708-.708l.5-.5a.5.5 0 0 1 .708 0z"/>
//
// </svg>
func TextareaResize(w, h int) *jse.SVG {
	var svg = jse.NewSVG("http://www.w3.org/2000/svg")
	svg.Width(w)
	svg.Height(h)
	svg.Fill("currentColor")
	svg.Element().ClassList("bi", "bi-textarea-resize")
	svg.ViewBox(0, 0, 16, 16)
	var path = svg.NewElement("path")
	path.SetAttr("d", "M0 4.5A2.5 2.5 0 0 1 2.5 2h11A2.5 2.5 0 0 1 16 4.5v7a2.5 2.5 0 0 1-2.5 2.5h-11A2.5 2.5 0 0 1 0 11.5v-7zM2.5 3A1.5 1.5 0 0 0 1 4.5v7A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5v-7A1.5 1.5 0 0 0 13.5 3h-11zm10.854 4.646a.5.5 0 0 1 0 .708l-3 3a.5.5 0 0 1-.708-.708l3-3a.5.5 0 0 1 .708 0zm0 2.5a.5.5 0 0 1 0 .708l-.5.5a.5.5 0 0 1-.708-.708l.5-.5a.5.5 0 0 1 .708 0z")
	return svg
}
