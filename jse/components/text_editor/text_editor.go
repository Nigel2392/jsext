package texteditor

import (
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
	"github.com/Nigel2392/jsext/v2/jse/components/svg"
)

type Editor struct {
	*jse.Element              // The root element
	Lint         *jse.Element // Lint element
	Editor       *jse.Element // Editor element
	Footer       *jse.Element // Footer element
	Conf         *Config
}

type LintElement struct {
	*jse.Element
	OnClick    func(*Editor)
	ConfigData map[string]interface{}
}

type Config struct {
	Placeholder string
	Value       string

	ButtonSize int
	Width      string
	Height     string
	MinHeight  string
	MaxHeight  string
	MinWidth   string
	MaxWidth   string

	AllowResizeX bool
	AllowResizeY bool

	LintElements     map[string]*LintElement // [name]LintElement
	LintElementOrder []string
	Options          []TextEditorOption
}

func (c *Config) SetOrder(order ...string) {
	if c.LintElementOrder == nil {
		c.LintElementOrder = make([]string, 0)
	}
	c.LintElementOrder = order
}

func (c *Config) AddLintElement(name string, svgFunc func(w, h int) *jse.SVG, onclick func(*Editor), config map[string]interface{}) {
	if c.ButtonSize == 0 {
		c.ButtonSize = 16
	}
	if c.LintElements == nil {
		c.LintElements = make(map[string]*LintElement)
	}
	c.LintElements[name] = makeLintElement(c.ButtonSize, svgFunc, onclick, config)
}

func (c *Config) AddLintElementCommand(name string, svgFunc func(w, h int) *jse.SVG, command string) {
	if c.ButtonSize == 0 {
		c.ButtonSize = 16
	}
	if c.LintElements == nil {
		c.LintElements = make(map[string]*LintElement)
	}
	c.LintElements[name] = makeLintElementCommand(c.ButtonSize, svgFunc, command)
}

func makeLintElement(size int, svgFunc func(w, h int) *jse.SVG, onclick func(*Editor), config map[string]interface{}) *LintElement {
	var elem = jse.NewElement("button")
	elem.AppendChild(svgFunc(size, size).Element())
	return &LintElement{
		Element:    elem,
		OnClick:    onclick,
		ConfigData: config,
	}
}

func makeLintElementCommand(size int, svgFunc func(w, h int) *jse.SVG, command string) *LintElement {
	return makeLintElement(size, svgFunc, func(e *Editor) {
		// replace the selection with the command
		jsext.Document.Call("execCommand", command, false, nil)
	}, nil)
}

func (c *Config) Defaults() {
	if c.LintElements == nil {
		c.LintElements = make(map[string]*LintElement)
	}
	if c.LintElementOrder == nil {
		c.LintElementOrder = make([]string, 0)
	}
	if c.Options == nil {
		c.Options = make([]TextEditorOption, 0)
	}
	if c.ButtonSize == 0 {
		c.ButtonSize = 16
	}
	if c.Width == "" {
		c.Width = "100%"
	}
	if c.Height == "" {
		c.Height = "350px"
	}
	if c.MaxHeight == "" {
		c.MaxHeight = "1000px"
	}
	if c.MaxWidth == "" {
		c.MaxWidth = "100%"
	}

	c.LintElements["bold"] = makeLintElementCommand(c.ButtonSize, svg.Bold, "bold")
	c.LintElements["italic"] = makeLintElementCommand(c.ButtonSize, svg.Italic, "italic")
	c.LintElements["underline"] = makeLintElementCommand(c.ButtonSize, svg.Underline, "underline")
	c.LintElements["strikethrough"] = makeLintElementCommand(c.ButtonSize, svg.StrikeThrough, "strikeThrough")
	c.LintElements["ol"] = makeLintElementCommand(c.ButtonSize, svg.ListOL, "insertOrderedList")
	c.LintElements["ul"] = makeLintElementCommand(c.ButtonSize, svg.ListUL, "insertUnorderedList")
	c.LintElements["text-left"] = makeLintElementCommand(c.ButtonSize, svg.TextLeft, "justifyLeft")
	c.LintElements["text-center"] = makeLintElementCommand(c.ButtonSize, svg.TextCenter, "justifyCenter")
	c.LintElements["text-right"] = makeLintElementCommand(c.ButtonSize, svg.TextRight, "justifyRight")

	c.LintElements["link"] = makeLintElement(c.ButtonSize, svg.Link, func(e *Editor) {
		// check if there is a selection
		var selected = e.CurrentlySelectedText()
		var _, err = url.Parse(selected)
		if err == nil {
			jsext.Document.Call("execCommand", "createLink", false, selected)
		}
	}, nil)

	if len(c.LintElementOrder) == 0 {
		c.LintElementOrder = []string{
			"bold", "italic", "underline",
			"strikethrough", "ol", "ul", "link",
			"text-left", "text-center",
			"text-right", //"image",
		}
	}
}

type TextEditorOption func(*Editor)

func (e *Editor) CurrentlySelectedText() string {
	var selected = jsext.Window.Call("getSelection")
	if selected.Get("isCollapsed").Bool() {
		return ""
	}
	return selected.Call("toString").String()
}

func New(conf *Config) *Editor {
	if conf == nil {
		conf = &Config{}
	}

	conf.Defaults()

	var e = jse.NewElement("text-editor-box")
	var editor = &Editor{
		Element: e,
		Lint:    e.NewElement("text-editor-lint"),
		Editor:  e.NewElement("text-editor"),
		Footer:  e.NewElement("text-editor-footer"),

		Conf: conf,
	}

	editor.Element.Style().Width(conf.Width)
	editor.Element.Style().Height(conf.Height)
	editor.Element.Style().MaxHeight(conf.MaxHeight)
	editor.Element.Style().MaxWidth(conf.MaxWidth)
	if conf.MinHeight != "" {
		editor.Element.Style().MinHeight(conf.MinHeight)
	}
	if conf.MinWidth != "" {
		editor.Element.Style().MinWidth(conf.MinWidth)
	}
	editor.Editor.SetAttr("contenteditable", "true")

	if conf.AllowResizeX || conf.AllowResizeY {
		var resize = svg.TextareaResize(conf.ButtonSize, conf.ButtonSize)
		resize.Element().Style().Cursor("se-resize")
		resize.Element().AddEventListener("mousedown", func(_ *jse.Element, e jsext.Event) {
			e.PreventDefault()
			var mousemove = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
				var event = args[0]
				event.Call("preventDefault")
				if conf.AllowResizeX {
					var x = event.Get("clientX").Int()
					var w = x - editor.Element.Get("offsetLeft").Int()
					editor.Element.Style().Width(fmt.Sprintf("%dpx", w))
				}
				if conf.AllowResizeY {
					var y = event.Get("clientY").Int()
					var h = y - editor.Element.Get("offsetTop").Int()
					editor.Element.Style().Height(fmt.Sprintf("%dpx", h))
				}
				return nil
			})
			var mouseup *js.Func
			var mUp = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
				var event = args[0]
				event.Call("preventDefault")
				js.Global().Call("removeEventListener", "mousemove", mousemove)
				js.Global().Call("removeEventListener", "mouseup", *mouseup)
				return nil
			})
			mouseup = &mUp
			js.Global().Call("addEventListener", "mousemove", mousemove)
			js.Global().Call("addEventListener", "mouseup", *mouseup)
		})
		editor.Footer.AppendChild(resize.Element())
	}

	var lintElements = make([]*LintElement, 0, len(conf.LintElementOrder))
	for _, name := range conf.LintElementOrder {
		var elem = conf.LintElements[name]
		lintElements = append(lintElements, elem)
		elem.Element.SetAttr("data-command", name)
		elem.Element.OnClick(func(this *jse.Element, event jsext.Event) {
			event.PreventDefault()
			editor.Editor.Call("focus")
			elem.OnClick(editor)
		})
		editor.Lint.AppendChild(elem.Element)
	}

	for _, opt := range conf.Options {
		opt(editor)
	}

	editor.Element.StyleBlock(`
		text-editor-box {
			display: flex;
			flex-direction: column;
			padding: 5px 15px;
		}
		text-editor-lint {
			border-radius: 0.5em 0.5em 0 0;
			display: flex;
			flex-direction: row;
			align-items: center;
			flex-wrap: wrap;
			padding: 3px 10px;  
		}
		text-editor-lint button {
			border-radius: 3px;
			padding: 5px;
			margin: 5px;
			cursor: pointer;
			font-weight: bold;
			text-transform: uppercase;
			font-size: 0.8rem;
			outline: none;
		}
		text-editor-lint button:nth-child(1) {
			margin-left: 0;
		}
		text-editor {
			flex: 1;
			padding: 5px;
			overflow: auto;
		}
		text-editor:focus {
			outline: none;
		}
		text-editor-footer {
			padding: 5px 15px;
			border-radius: 0 0 0.5em 0.5em;
			display: flex;
			flex-direction: row;
			align-items: center;
			justify-content: flex-end;  
		}
		`)
	return editor
}

func WithEditorHeight(h string) TextEditorOption {
	return func(e *Editor) {
		e.Editor.Style().Height(fmt.Sprint(h))
	}
}

func WithEditorWidth(w string) TextEditorOption {
	return func(e *Editor) {
		e.Editor.Style().Width(fmt.Sprint(w))
	}
}

func (e *Editor) Value() *jse.Element {
	return e.Editor
}

func (e *Editor) JSValue() js.Value {
	return e.Editor.JSValue()
}
