package jse

import (
	"github.com/Nigel2392/jsext/v2"
)

func Select(name, id string, opts *InputOptions) *SelectElement {
	var e = jsext.CreateElement("select")
	e.SetAttribute("name", name)
	e.SetAttribute("id", id)
	opts.Apply(e)
	return (*SelectElement)(&e)
}

type SelectElement jsext.Element

func (s *SelectElement) Value() jsext.Value {
	return jsext.Value(*s)
}

func (s *SelectElement) Element() *Element {
	return (*Element)(s)
}

func (s *SelectElement) ClassList(cls ...string) jsext.Value {
	return s.Element().ClassList(cls...)
}

func (s *SelectElement) OnChange(listener func(this *Element, event jsext.Event)) *SelectElement {
	s.Element().AddEventListener("change", listener)
	return s
}

func (s *SelectElement) Append(child any) *SelectElement {
	switch v := child.(type) {
	case jsext.Element:
		s.Value().AppendChild(v)
	case jsext.Value:
		s.Value().AppendChild(jsext.Element(v))
	case *OptionElement:
		s.Value().AppendChild(v.Value())
	case []*OptionElement:
		for _, o := range v {
			s.Value().AppendChild(o.Value())
		}
	case OptionElement:
		s.Value().AppendChild(v.Value())
	case []OptionElement:
		for _, o := range v {
			s.Value().AppendChild(o.Value())
		}
	default:
		panic("invalid child type")
	}
	return s
}

func (s *SelectElement) FormValue() string {
	var v = s.Value().Get("value")
	if v.IsNull() || v.IsUndefined() {
		return ""
	}
	return v.String()
}

func (s *SelectElement) InlineClasses(classes ...string) *SelectElement {
	s.ClassList().Add(classes...)
	return s
}

func (s *SelectElement) Set(name string, value any) *SelectElement {
	s.Value().Set(name, value)
	return s
}

func (s *SelectElement) AddEventListener(name string, listener func(this *Element, event jsext.Event)) *SelectElement {
	s.Element().AddEventListener(name, listener)
	return s
}

type OptionElement jsext.Element

func Option(text, value string, selected ...bool) *OptionElement {
	var e = jsext.CreateElement("option")
	e.SetAttribute("value", value)
	e.InnerText(text)
	if len(selected) > 0 && selected[0] {
		e.Set("selected", true)
	}
	return (*OptionElement)(&e)
}

func (s *SelectElement) Option(text, value string, selected ...bool) *OptionElement {
	var e = Option(text, value, selected...)
	s.Value().AppendChild(e.Value())
	return e
}

func (o *OptionElement) Value() jsext.Element {
	return jsext.Element(*o)
}

func (o *OptionElement) Selected(b bool) {
	if b {
		o.Value().Set("selected", true)
	} else {
		o.Value().Delete("selected")
	}
}
