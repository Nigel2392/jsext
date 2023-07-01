package jse

import (
	"strconv"

	"github.com/Nigel2392/jsext/v2"
)

type InputOptions struct {
	Required      bool
	Disabled      bool
	ReadOnly      bool
	Autofocus     bool
	Checked       bool
	AutoComplete  string
	Placeholder   string
	Multiple      bool
	Min           int
	Max           int
	Step          int
	Value         string
	Classlist     []string
	ButtonClasses []string
}

func (o *InputOptions) Apply(e jsext.Element) {
	if o == nil {
		return
	}
	if o.Required {
		e.SetAttribute("required", "true")
	}
	if o.Disabled {
		e.SetAttribute("disabled", "true")
	}
	if o.ReadOnly {
		e.SetAttribute("readonly", "true")
	}
	if o.Autofocus {
		e.SetAttribute("autofocus", "true")
	}
	if o.Checked {
		e.SetAttribute("checked", "true")
	}
	if o.AutoComplete != "" {
		e.SetAttribute("autocomplete", o.AutoComplete)
	}
	if o.Placeholder != "" {
		e.SetAttribute("placeholder", o.Placeholder)
	}
	if o.Multiple {
		e.SetAttribute("multiple", "true")
	}
	if o.Min != 0 {
		e.SetAttribute("min", strconv.Itoa(o.Min))
	}
	if o.Max != 0 {
		e.SetAttribute("max", strconv.Itoa(o.Max))
	}
	if o.Step != 0 {
		e.SetAttribute("step", strconv.Itoa(o.Step))
	}
	if o.Value != "" {
		e.SetAttribute("value", o.Value)
	}
	if len(o.Classlist) > 0 {
		e.ClassList(o.Classlist...)
	}
}
