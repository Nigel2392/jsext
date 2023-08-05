package shortcuts

import (
	"strconv"
	"strings"
	"syscall/js"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
)

var (
	DEFAULT_MODAL_APPEND_TO js.Value = js.Global().Get("document").Get("body")

	DEFAULT_MODAL_CLASS_PREFIX = "jsext-"

	DEFAULT_MODAL_FILL_COLOR      string  = "rgba(0,0,0,0.5)"
	DEFAULT_MODAL_COLOR           string  = "#232323"
	DEFAULT_MODAL_BG_COLOR        string  = "#fff"
	DEFAULT_MODAL_BORDER_RADIUS   string  = "5px"
	DEFAULT_MODAL_BLOCK_BORDER    string  = "1px solid #ccc;"
	DEFAULT_MODAL_BORDER          string  = "1px solid #ccc"
	DEFAULT_MODAL_WIDTH           string  = "50%"
	DEFAULT_MODAL_HEIGHT          string  = "auto"
	DEFAULT_MODAL_OVERFLOW        string  = "auto"
	DEFAULT_MODAL_CLOSE_BTN_SCALE float64 = 1
	DEFAULT_MODAL_Z_INDEX         int     = 5
)

type Modal jse.Element

func (m *Modal) Show() {
	if !m.Element().Value().Truthy() {
		panic("modal not created or is underfined")
	}
	m.Element().Style().Display("flex")
}

func (m *Modal) Hide() {
	if !m.Element().Value().Truthy() {
		return
	}
	m.Element().Style().Display("none")
}

func (m *Modal) Element() jsext.Element {
	return (jsext.Element)(*m)
}

func (m *Modal) Create(appendToQuerySelector ...string) {
	if len(appendToQuerySelector) > 0 && appendToQuerySelector[0] != "" {
		var e, err = jsext.QuerySelector(appendToQuerySelector[0])
		if err != nil {
			panic(err)
		}
		e.Append(m.Element())
	} else {
		DEFAULT_MODAL_APPEND_TO.Call("appendChild", m.Element().JSValue())
	}
}

func (m *Modal) OpenOnClickOf(e *jse.Element) {
	e.AddEventListener("click", func(this *jse.Element, event jsext.Event) {
		var preventDefault = event.Get("preventDefault")
		if preventDefault.Truthy() {
			event.PreventDefault()
		}
		m.Show()
	})
}

func (m *Modal) CloseOnClickOf(e *jse.Element) {
	e.AddEventListener("click", func(this *jse.Element, event jsext.Event) {
		var preventDefault = event.Get("preventDefault")
		if preventDefault.Truthy() {
			event.PreventDefault()
		}
		m.Delete()
	})
}

type ButtonType int

const (
	BtnTypeAnchor ButtonType = 0
	BtnTypeButton ButtonType = 1
)

func (m *Modal) Button(tag ButtonType, innerText string) *jse.Element {
	var btn *jse.Element
	switch tag {
	case BtnTypeAnchor:
		btn = jse.A("javascript:void(0)", innerText)
	case BtnTypeButton:
		btn = jse.Button(innerText, nil)
	}
	m.OpenOnClickOf(btn)
	return btn
}

func (m *Modal) Delete() {
	if m.Element().Value().Truthy() {
		m.Element().Value().Remove()
	}
}

type ModalOptions struct {
	Header           *jse.Element
	Body             *jse.Element
	Footer           *jse.Element
	Background       string
	ModalBackground  string
	BorderRadius     string
	BlockBorder      string
	Border           string
	Width            string
	Height           string
	Color            string
	ClassPrefix      string
	CloseButton      bool
	DeleteOnClose    bool
	CloseButtonScale float64
	ZIndex           int
	OverflowX        string
	OverflowY        string
	Overflow         string
}

func ModalConfirm[T jsext.Element | *jsext.Element | jse.Element | *jse.Element | string](title string, inner T, onConfirm func(*Modal)) *Modal {
	var innerElem *jse.Element
	switch inner := any(inner).(type) {
	case string:
		innerElem = jse.Div().P().InnerText(inner)
	case jse.Element:
		innerElem = (*jse.Element)(&inner)
	case *jse.Element:
		innerElem = inner
	case jsext.Element:
		innerElem = (*jse.Element)(&inner)
	case *jsext.Element:
		innerElem = (*jse.Element)(inner)
	}
	var confirm_btn = jse.Button("Confirm", nil)
	confirm_btn.ClassList("btn", "btn-primary")
	var opts = ModalOptions{
		Header: jse.Div().Heading(4, title),
		Body:   innerElem,
		Footer: jse.Div().AppendChild(
			confirm_btn,
		),
		CloseButton: true,
	}

	var modal = CreateModal(opts)

	confirm_btn.AddEventListener("click", func(this *jse.Element, e jsext.Event) {
		e.PreventDefault()
		onConfirm(modal)
		modal.Delete()
	})

	return modal
}

func (opts *ModalOptions) SetDefaults() {
	if opts.ClassPrefix == "" {
		opts.ClassPrefix = DEFAULT_MODAL_CLASS_PREFIX
	}
	if opts.Background == "" {
		opts.Background = DEFAULT_MODAL_FILL_COLOR
	}
	if opts.ModalBackground == "" {
		opts.ModalBackground = DEFAULT_MODAL_BG_COLOR
	}
	if opts.BorderRadius == "" {
		opts.BorderRadius = DEFAULT_MODAL_BORDER_RADIUS
	}
	if opts.Border == "" {
		opts.Border = DEFAULT_MODAL_BORDER
	}
	if opts.Width == "" {
		opts.Width = DEFAULT_MODAL_WIDTH
	}
	if opts.Height == "" {
		opts.Height = DEFAULT_MODAL_HEIGHT
	}
	if opts.CloseButtonScale == 0 {
		opts.CloseButtonScale = DEFAULT_MODAL_CLOSE_BTN_SCALE
	}
	if opts.ZIndex == 0 {
		opts.ZIndex = DEFAULT_MODAL_Z_INDEX
	}
	if opts.Color == "" {
		opts.Color = DEFAULT_MODAL_COLOR
	}
	if opts.BlockBorder == "" {
		opts.BlockBorder = DEFAULT_MODAL_BLOCK_BORDER
	}
	if opts.Overflow == "" {
		if opts.OverflowX == "" {
			opts.OverflowX = DEFAULT_MODAL_OVERFLOW
		}
		if opts.OverflowY == "" {
			opts.OverflowY = DEFAULT_MODAL_OVERFLOW
		}
	} else {
		opts.OverflowX = opts.Overflow
		opts.OverflowY = opts.Overflow
	}
}

func CreateModal(opts ModalOptions) *Modal {
	opts.SetDefaults()
	var modal_container *jse.Element = jse.Div()
	modal_container.ClassList(opts.ClassPrefix + "modal-container")
	var modal = modal_container.Div()
	if opts.CloseButton {
		var close_btn = modal.Div()
		close_btn.ClassList(opts.ClassPrefix + "close-btn")
		close_btn.AddEventListener("click", func(this *jse.Element, e jsext.Event) {
			if opts.DeleteOnClose {
				modal_container.Remove()
			} else {
				(*Modal)(modal_container).Hide()
			}
		})
	}
	modal.ClassList(opts.ClassPrefix + "modal")
	if opts.Header != nil {
		modal.AppendChild(opts.Header)
	}
	if opts.Body != nil {
		modal.AppendChild(opts.Body)
	}
	if opts.Footer != nil {
		modal.AppendChild(opts.Footer)
	}

	css := strings.Join([]string{`.`, opts.ClassPrefix, `modal-container {
			position: fixed;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background: `, opts.Background, `;
			display: none;
			justify-content: center;
			align-items: center;
			z-index: `, strconv.Itoa(opts.ZIndex), `;
		}
		.`, opts.ClassPrefix, `close-btn {
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
			transform: scale(`, strconv.FormatFloat(opts.CloseButtonScale, 'f', 2, 64), `);
			transition: background 0.2s ease-in-out;
		}
		.`, opts.ClassPrefix, `close-btn:hover {
			background: #910c0c;
		}
		.`, opts.ClassPrefix, `close-btn::before {
			position: absolute;
			content: "";
			width: 20px;
			height: 2px;
			background: #fff;
			transform: rotate(45deg);
		}
		.`, opts.ClassPrefix, `close-btn::after {
			position: absolute;
			content: "";
			width: 20px;
			height: 2px;
			background: #fff;
			transform: rotate(-45deg);
		}
		
		.`, opts.ClassPrefix, `modal {
			background: ` + opts.ModalBackground + `;
			border-radius: ` + opts.BorderRadius + `;
			border: ` + opts.Border + `;
			color: ` + opts.Color + `;
			width: ` + opts.Width + `;
			height: ` + opts.Height + `;
			overflow-x: ` + opts.OverflowX + `;
			overflow-y: ` + opts.OverflowY + `;
			max-width: 95%;
			max-height: 95%;
			display: flex;
			flex-direction: column;
			position: relative;
		}
		.`, opts.ClassPrefix, `modal > * {
			padding: 10px;
		}
		.`, opts.ClassPrefix, `modal > *:first-child {
			border-bottom: ` + opts.BlockBorder + `
		}
		.`, opts.ClassPrefix, `modal > *:last-child {
			border-top: ` + opts.BlockBorder + `
		}
		.`, opts.ClassPrefix, `modal > *:only-child {
			border: none;
		}`}, "")

	modal.StyleBlock(css)
	return (*Modal)(modal_container)
}
