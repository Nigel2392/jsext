//go:build js && wasm
// +build js,wasm

package forms

import (
	"strings"
	"time"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/elements"
)

// Form delimiter for nested structs
const Delimiter = "___"

// Form struct
type Form struct {
	Inner      *elements.Element
	Validators map[string]func(string) error
}

// Get a new form
func NewForm(action string, method string) *Form {
	return &Form{
		Inner:      elements.Form(action, method),
		Validators: make(map[string]func(string) error),
	}
}

// Get the form value from the form element
func (f *Form) Value() jsext.Element {
	return f.Inner.JSExtElement()
}

// Get the form value from the form element
func (f *Form) Element() *elements.Element {
	return f.Inner
}

// Render the form
func (f *Form) Render() jsext.Element {
	return f.Inner.Render()
}

// Set the form ID
func (f *Form) AttrID(id string) *Form {
	f.Inner.AttrID(id)
	return f
}

// Add a validator to the form.
// The validator will be called when the form is submitted.
// WATCH OUT!
//
//   - The form names could be transformed when you parse a struct into a form;
//     Please use the same name as the struct field name, for embedded structs
//     the delimiter is used:
//
//     "Field" for a normal field
//     "Other___Field" for an embedded field
func (f *Form) AddValidator(name string, fn func(string) error) {
	f.Validators[name] = fn
}

// Eventlistener for when the form is submitted.
func (f *Form) OnSubmit(cb func(data map[string]string, elements []jsext.Element)) {
	f.Inner.AddEventListener("submit", func(this jsext.Value, event jsext.Event) {
		event.PreventDefault()
		var data = make(map[string]string)
		var elements = this.Get("elements")
		var length = elements.Get("length").Int()
		var elemList = make([]jsext.Element, length)
		for i := 0; i < length; i++ {
			var element = elements.Index(i)
			var name = element.Get("name").String()
			var value = element.Get("value").String()
			if name == "" {
				continue
			}

			if fn, ok := f.Validators[name]; ok {
				var err = fn(value)
				if err != nil {
					return
				}
			}

			data[name] = value
			elemList[i] = element.ToElement()
		}
		cb(data, elemList)
	})
}

// Valid formtypes
type FORMTYPES string

// Check if two formtypes are equal
func (ft FORMTYPES) Equals(other FORMTYPES) bool {
	return strings.EqualFold(string(ft), string(other))
}

// Formtypes to use in forms
const (
	FORMTYP_TEXT     FORMTYPES = "text"
	FORMTYP_CHECKBOX FORMTYPES = "checkbox"
	FORMTYP_NUMBER   FORMTYPES = "number"
	FORMTYP_DATETIME FORMTYPES = "datetime-local"
	FORMTYP_FILE     FORMTYPES = "file"
	FORMTYP_STRUCT   FORMTYPES = "struct"
	// FORMTYP_LIST     FORMTYPES = "select"
	FORMTYP_INVALID FORMTYPES = "text"
)

func isValidTyp(typ FORMTYPES) bool {
	switch typ {
	case FORMTYP_TEXT, FORMTYP_CHECKBOX, FORMTYP_NUMBER, FORMTYP_DATETIME, FORMTYP_FILE: //, FORMTYP_LIST:
		return true
	default:
		return false
	}
}

// Format a value for display in a form
func FormatIfDateTime(val any) any {
	switch val := val.(type) {
	case time.Time:
		return val.Format("2006-01-02T15:04")
	default:
		return val
	}
}
