//go:build js && wasm
// +build js,wasm

package jse

import "github.com/Nigel2392/jsext/v2"

// A Element in method format: (Autogenerated)
func (e *Element) A(href string, text ...string) *Element {
	var ne = A(href, text...)
	e.AppendChild(ne)
	return ne
}

// Abbr Element in method format: (Autogenerated)
func (e *Element) Abbr(title string, innerText ...string) *Element {
	var ne = Abbr(title, innerText...)
	e.AppendChild(ne)
	return ne
}

// Address Element in method format: (Autogenerated)
func (e *Element) Address(innerText ...string) *Element {
	var ne = Address(innerText...)
	e.AppendChild(ne)
	return ne
}

// Area Element in method format: (Autogenerated)
func (e *Element) Area(alt string, coords ...string) *Element {
	var ne = Area(alt, coords...)
	e.AppendChild(ne)
	return ne
}

// Article Element in method format: (Autogenerated)
func (e *Element) Article(innerText ...string) *Element {
	var ne = Article(innerText...)
	e.AppendChild(ne)
	return ne
}

// Aside Element in method format: (Autogenerated)
func (e *Element) Aside(innerText ...string) *Element {
	var ne = Aside(innerText...)
	e.AppendChild(ne)
	return ne
}

// Audio Element in method format: (Autogenerated)
func (e *Element) Audio(src string) *Element {
	var ne = Audio(src)
	e.AppendChild(ne)
	return ne
}

// B Element in method format: (Autogenerated)
func (e *Element) B(innerText ...string) *Element {
	var ne = B(innerText...)
	e.AppendChild(ne)
	return ne
}

// Bdi Element in method format: (Autogenerated)
func (e *Element) Bdi(innerText ...string) *Element {
	var ne = Bdi(innerText...)
	e.AppendChild(ne)
	return ne
}

// Bdo Element in method format: (Autogenerated)
func (e *Element) Bdo(innerText ...string) *Element {
	var ne = Bdo(innerText...)
	e.AppendChild(ne)
	return ne
}

// Blockquote Element in method format: (Autogenerated)
func (e *Element) Blockquote(innerText ...string) *Element {
	var ne = Blockquote(innerText...)
	e.AppendChild(ne)
	return ne
}

// Body Element in method format: (Autogenerated)
func (e *Element) Body(innerText ...string) *Element {
	var ne = Body(innerText...)
	e.AppendChild(ne)
	return ne
}

// Br Element in method format: (Autogenerated)
func (e *Element) Br() *Element {
	var ne = Br()
	e.AppendChild(ne)
	return ne
}

// Button Element in method format: (Autogenerated)
func (e *Element) Button(innerText string, onClick func(this *Element, event jsext.Event)) *Element {
	var ne = Button(innerText, onClick)
	e.AppendChild(ne)
	return ne
}

// Canvas Element in method format: (Autogenerated)
func (e *Element) Canvas(innerText ...string) *Element {
	var ne = Canvas(innerText...)
	e.AppendChild(ne)
	return ne
}

// Caption Element in method format: (Autogenerated)
func (e *Element) Caption(innerText ...string) *Element {
	var ne = Caption(innerText...)
	e.AppendChild(ne)
	return ne
}

// Cite Element in method format: (Autogenerated)
func (e *Element) Cite(innerText ...string) *Element {
	var ne = Cite(innerText...)
	e.AppendChild(ne)
	return ne
}

// Code Element in method format: (Autogenerated)
func (e *Element) Code(innerText ...string) *Element {
	var ne = Code(innerText...)
	e.AppendChild(ne)
	return ne
}

// Col Element in method format: (Autogenerated)
func (e *Element) Col() *Element {
	var ne = Col()
	e.AppendChild(ne)
	return ne
}

// Colgroup Element in method format: (Autogenerated)
func (e *Element) Colgroup(innerText ...string) *Element {
	var ne = Colgroup(innerText...)
	e.AppendChild(ne)
	return ne
}

// Command Element in method format: (Autogenerated)
func (e *Element) Command(innerText ...string) *Element {
	var ne = Command(innerText...)
	e.AppendChild(ne)
	return ne
}

// Datalist Element in method format: (Autogenerated)
func (e *Element) Datalist(innerText ...string) *Element {
	var ne = Datalist(innerText...)
	e.AppendChild(ne)
	return ne
}

// Dd Element in method format: (Autogenerated)
func (e *Element) Dd(innerText ...string) *Element {
	var ne = Dd(innerText...)
	e.AppendChild(ne)
	return ne
}

// Del Element in method format: (Autogenerated)
func (e *Element) Del(innerText ...string) *Element {
	var ne = Del(innerText...)
	e.AppendChild(ne)
	return ne
}

// Details Element in method format: (Autogenerated)
func (e *Element) Details(innerText ...string) *Element {
	var ne = Details(innerText...)
	e.AppendChild(ne)
	return ne
}

// Dfn Element in method format: (Autogenerated)
func (e *Element) Dfn(innerText ...string) *Element {
	var ne = Dfn(innerText...)
	e.AppendChild(ne)
	return ne
}

// Div Element in method format: (Autogenerated)
func (e *Element) Div(classlist ...string) *Element {
	var ne = Div(classlist...)
	e.AppendChild(ne)
	return ne
}

// Dl Element in method format: (Autogenerated)
func (e *Element) Dl(innerText ...string) *Element {
	var ne = Dl(innerText...)
	e.AppendChild(ne)
	return ne
}

// Dt Element in method format: (Autogenerated)
func (e *Element) Dt(innerText ...string) *Element {
	var ne = Dt(innerText...)
	e.AppendChild(ne)
	return ne
}

// Em Element in method format: (Autogenerated)
func (e *Element) Em(innerText ...string) *Element {
	var ne = Em(innerText...)
	e.AppendChild(ne)
	return ne
}

// Fieldset Element in method format: (Autogenerated)
func (e *Element) Fieldset(innerText ...string) *Element {
	var ne = Fieldset(innerText...)
	e.AppendChild(ne)
	return ne
}

// Figcaption Element in method format: (Autogenerated)
func (e *Element) Figcaption(innerText ...string) *Element {
	var ne = Figcaption(innerText...)
	e.AppendChild(ne)
	return ne
}

// Figure Element in method format: (Autogenerated)
func (e *Element) Figure(innerText ...string) *Element {
	var ne = Figure(innerText...)
	e.AppendChild(ne)
	return ne
}

// Footer Element in method format: (Autogenerated)
func (e *Element) Footer(innerText ...string) *Element {
	var ne = Footer(innerText...)
	e.AppendChild(ne)
	return ne
}

// Form Element in method format: (Autogenerated)
func (e *Element) Form(action, method, id string) *FormElement {
	var ne = Form(action, method, id)
	e.AppendChild((*Element)(ne))
	return ne
}

func (e *FormElement) InlineClasses(s ...string) *FormElement {
	e.ClassList(s...)
	return e
}

// H1 Element in method format: (Autogenerated)
func (e *Element) Heading(size int, innerText ...string) *Element {
	var ne = Heading(size, innerText...)
	e.AppendChild(ne)
	return ne
}

// Header Element in method format: (Autogenerated)
func (e *Element) Header(innerText ...string) *Element {
	var ne = Header(innerText...)
	e.AppendChild(ne)
	return ne
}

// Hr Element in method format: (Autogenerated)
func (e *Element) Hr() *Element {
	var ne = Hr()
	e.AppendChild(ne)
	return ne
}

// Html Element in method format: (Autogenerated)
func (e *Element) Html(innerText ...string) *Element {
	var ne = Html(innerText...)
	e.AppendChild(ne)
	return ne
}

// I Element in method format: (Autogenerated)
func (e *Element) I(innerText ...string) *Element {
	var ne = I(innerText...)
	e.AppendChild(ne)
	return ne
}

// Iframe Element in method format: (Autogenerated)
func (e *Element) Iframe(src string) *Element {
	var ne = Iframe(src)
	e.AppendChild(ne)
	return ne
}

// Img Element in method format: (Autogenerated)
func (e *Element) Img(src string) *Element {
	var ne = Img(src)
	e.AppendChild(ne)
	return ne
}

// Ins Element in method format: (Autogenerated)
func (e *Element) Ins(innerText ...string) *Element {
	var ne = Ins(innerText...)
	e.AppendChild(ne)
	return ne
}

// Kbd Element in method format: (Autogenerated)
func (e *Element) Kbd(innerText ...string) *Element {
	var ne = Kbd(innerText...)
	e.AppendChild(ne)
	return ne
}

// Legend Element in method format: (Autogenerated)
func (e *Element) Legend(innerText ...string) *Element {
	var ne = Legend(innerText...)
	e.AppendChild(ne)
	return ne
}

// Li Element in method format: (Autogenerated)
func (e *Element) Li(innerText ...string) *Element {
	var ne = Li(innerText...)
	e.AppendChild(ne)
	return ne
}

// Link Element in method format: (Autogenerated)
func (e *Element) Link(href string) *Element {
	var ne = Link(href)
	e.AppendChild(ne)
	return ne
}

// Main Element in method format: (Autogenerated)
func (e *Element) Main(innerText ...string) *Element {
	var ne = Main(innerText...)
	e.AppendChild(ne)
	return ne
}

// Map Element in method format: (Autogenerated)
func (e *Element) Map(innerText ...string) *Element {
	var ne = Map(innerText...)
	e.AppendChild(ne)
	return ne
}

// Mark Element in method format: (Autogenerated)
func (e *Element) Mark(innerText ...string) *Element {
	var ne = Mark(innerText...)
	e.AppendChild(ne)
	return ne
}

// Meta Element in method format: (Autogenerated)
func (e *Element) Meta(name string, content string) *Element {
	var ne = Meta(name, content)
	e.AppendChild(ne)
	return ne
}

// Meter Element in method format: (Autogenerated)
func (e *Element) Meter(innerText ...string) *Element {
	var ne = Meter(innerText...)
	e.AppendChild(ne)
	return ne
}

// Nav Element in method format: (Autogenerated)
func (e *Element) Nav(innerText ...string) *Element {
	var ne = Nav(innerText...)
	e.AppendChild(ne)
	return ne
}

// Noscript Element in method format: (Autogenerated)
func (e *Element) Noscript(innerText ...string) *Element {
	var ne = Noscript(innerText...)
	e.AppendChild(ne)
	return ne
}

// Object Element in method format: (Autogenerated)
func (e *Element) Object(innerText ...string) *Element {
	var ne = Object(innerText...)
	e.AppendChild(ne)
	return ne
}

// Ol Element in method format: (Autogenerated)
func (e *Element) Ol(innerText ...string) *Element {
	var ne = Ol(innerText...)
	e.AppendChild(ne)
	return ne
}

// Optgroup Element in method format: (Autogenerated)
func (e *Element) Optgroup(innerText ...string) *Element {
	var ne = Optgroup(innerText...)
	e.AppendChild(ne)
	return ne
}

// Output Element in method format: (Autogenerated)
func (e *Element) Output(innerText ...string) *Element {
	var ne = Output(innerText...)
	e.AppendChild(ne)
	return ne
}

// P Element in method format: (Autogenerated)
func (e *Element) P(innerText ...string) *Element {
	var ne = P(innerText...)
	e.AppendChild(ne)
	return ne
}

// Param Element in method format: (Autogenerated)
func (e *Element) Param(innerText ...string) *Element {
	var ne = Param(innerText...)
	e.AppendChild(ne)
	return ne
}

// Pre Element in method format: (Autogenerated)
func (e *Element) Pre(innerText ...string) *Element {
	var ne = Pre(innerText...)
	e.AppendChild(ne)
	return ne
}

// Progress Element in method format: (Autogenerated)
func (e *Element) Progress(innerText ...string) *Element {
	var ne = Progress(innerText...)
	e.AppendChild(ne)
	return ne
}

// Q Element in method format: (Autogenerated)
func (e *Element) Q(innerText ...string) *Element {
	var ne = Q(innerText...)
	e.AppendChild(ne)
	return ne
}

// Rp Element in method format: (Autogenerated)
func (e *Element) Rp(innerText ...string) *Element {
	var ne = Rp(innerText...)
	e.AppendChild(ne)
	return ne
}

// Rt Element in method format: (Autogenerated)
func (e *Element) Rt(innerText ...string) *Element {
	var ne = Rt(innerText...)
	e.AppendChild(ne)
	return ne
}

// Ruby Element in method format: (Autogenerated)
func (e *Element) Ruby(innerText ...string) *Element {
	var ne = Ruby(innerText...)
	e.AppendChild(ne)
	return ne
}

// Samp Element in method format: (Autogenerated)
func (e *Element) Samp(innerText ...string) *Element {
	var ne = Samp(innerText...)
	e.AppendChild(ne)
	return ne
}

// Script Element in method format: (Autogenerated)
func (e *Element) Script(src string) *Element {
	var ne = Script(src)
	e.AppendChild(ne)
	return ne
}

func (e *Element) StyleBlock(src ...string) *Element {
	var ne = StyleBlock(src...)
	e.AppendChild(ne)
	return ne
}

// Section Element in method format: (Autogenerated)
func (e *Element) Section(innerText ...string) *Element {
	var ne = Section(innerText...)
	e.AppendChild(ne)
	return ne
}

// Small Element in method format: (Autogenerated)
func (e *Element) Small(innerText ...string) *Element {
	var ne = Small(innerText...)
	e.AppendChild(ne)
	return ne
}

// Source Element in method format: (Autogenerated)
func (e *Element) Source(src string) *Element {
	var ne = Source(src)
	e.AppendChild(ne)
	return ne
}

// Span Element in method format: (Autogenerated)
func (e *Element) Span(innerText ...string) *Element {
	var ne = Span(innerText...)
	e.AppendChild(ne)
	return ne
}

// Strong Element in method format: (Autogenerated)
func (e *Element) Strong(innerText ...string) *Element {
	var ne = Strong(innerText...)
	e.AppendChild(ne)
	return ne
}

// Sub Element in method format: (Autogenerated)
func (e *Element) Sub(innerText ...string) *Element {
	var ne = Sub(innerText...)
	e.AppendChild(ne)
	return ne
}

// Summary Element in method format: (Autogenerated)
func (e *Element) Summary(innerText ...string) *Element {
	var ne = Summary(innerText...)
	e.AppendChild(ne)
	return ne
}

// Sup Element in method format: (Autogenerated)
func (e *Element) Sup(innerText ...string) *Element {
	var ne = Sup(innerText...)
	e.AppendChild(ne)
	return ne
}

// Table Element in method format: (Autogenerated)
func (e *Element) Table(classList ...string) *Element {
	var ne = Table(classList...)
	e.AppendChild(ne)
	return ne
}

// Tbody Element in method format: (Autogenerated)
func (e *Element) Tbody(innerText ...string) *Element {
	var ne = Tbody(innerText...)
	e.AppendChild(ne)
	return ne
}

// Td Element in method format: (Autogenerated)
func (e *Element) Td(innerText ...string) *Element {
	var ne = Td(innerText...)
	e.AppendChild(ne)
	return ne
}

// Textarea Element in method format: (Autogenerated)
func (e *Element) Textarea(name string, placeholder string) *Element {
	var ne = Textarea(name, placeholder)
	e.AppendChild(ne)
	return ne
}

// Tfoot Element in method format: (Autogenerated)
func (e *Element) Tfoot(innerText ...string) *Element {
	var ne = Tfoot(innerText...)
	e.AppendChild(ne)
	return ne
}

// Th Element in method format: (Autogenerated)
func (e *Element) Th(innerText ...string) *Element {
	var ne = Th(innerText...)
	e.AppendChild(ne)
	return ne
}

// Thead Element in method format: (Autogenerated)
func (e *Element) Thead(innerText ...string) *Element {
	var ne = Thead(innerText...)
	e.AppendChild(ne)
	return ne
}

// Time Element in method format: (Autogenerated)
func (e *Element) Time(innerText ...string) *Element {
	var ne = Time(innerText...)
	e.AppendChild(ne)
	return ne
}

// Title Element in method format: (Autogenerated)
func (e *Element) Title(innerText ...string) *Element {
	var ne = Title(innerText...)
	e.AppendChild(ne)
	return ne
}

// Tr Element in method format: (Autogenerated)
func (e *Element) Tr(innerText ...string) *Element {
	var ne = Tr(innerText...)
	e.AppendChild(ne)
	return ne
}

// Track Element in method format: (Autogenerated)
func (e *Element) Track(innerText ...string) *Element {
	var ne = Track(innerText...)
	e.AppendChild(ne)
	return ne
}

// Ul Element in method format: (Autogenerated)
func (e *Element) Ul(innerText ...string) *Element {
	var ne = Ul(innerText...)
	e.AppendChild(ne)
	return ne
}

// Var Element in method format: (Autogenerated)
func (e *Element) Var(innerText ...string) *Element {
	var ne = Var(innerText...)
	e.AppendChild(ne)
	return ne
}

// Video Element in method format: (Autogenerated)
func (e *Element) Video(innerText ...string) *Element {
	var ne = Video(innerText...)
	e.AppendChild(ne)
	return ne
}

// Wbr Element in method format: (Autogenerated)
func (e *Element) Wbr(innerText ...string) *Element {
	var ne = Wbr(innerText...)
	e.AppendChild(ne)
	return ne
}

// Select Element in method format: (Autogenerated)
func (e *Element) Select(name, id string, opts *InputOptions) *SelectElement {
	var ne = Select(name, id, opts)
	e.AppendChild((*Element)(ne))
	return ne
}

func (e *Element) Option(text, value string, selected ...bool) *OptionElement {
	var ne = Option(text, value, selected...)
	e.AppendChild((*Element)(ne))
	return ne
}

func (e *Element) Label(forElement, text string, classes ...string) *Element {
	var ne = Label(forElement, text, classes...)
	e.AppendChild(ne)
	return ne
}

func (e *Element) Input(typ string, name string, opts *InputOptions) *Element {
	var ne = Input(typ, name, opts)
	e.AppendChild(ne)
	return ne
}
