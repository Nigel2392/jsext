//go:build js && wasm
// +build js,wasm

package elements

import (
	"strings"
)

func (e *Element) AttrStyle(styles ...string) *Element {
	if e.value.IsUndefined() {
		e.SetSemicolon("style", styles...)
		return e
	}

	// Get the current style
	var currentStyle = e.value.Get("style")

	var length = currentStyle.Length()
	var styleMap = make(map[string]string)
	if length > 0 {
		// If there is a current style, append the new style
		for i := 0; i < length; i++ {
			var key = currentStyle.Index(i).String()
			var value = currentStyle.Get(key).String()
			styleMap[key] = value
		}
	}

	// Add the new styles
	for _, style := range styles {
		// Split the style into key and value
		var split = strings.SplitN(style, ":", 2)
		if len(split) != 2 {
			continue
		}
		var key = strings.TrimSpace(split[0])
		var value = strings.TrimSpace(split[1])
		styleMap[key] = value
	}

	// Create the new style string
	var newStyle = make([]string, len(styleMap))
	var i = 0
	var b strings.Builder
	for k, v := range styleMap {
		b.Reset()
		b.WriteString(k)
		b.WriteString(":")
		b.WriteString(v)
		newStyle[i] = b.String()
		i++
	}

	e.SetSemicolon("style", newStyle...)
	return e
}

func (e *Element) AttrSrc(src string) *Element {
	e.Add("src", src)
	return e
}

func (e *Element) AttrAlt(alt string) *Element {
	e.Add("alt", alt)
	return e
}

func (e *Element) AttrTitle(title string) *Element {
	e.Add("title", title)
	return e
}

func (e *Element) AttrHref(href string) *Element {
	e.Add("href", href)
	return e
}

func (e *Element) AttrClass(classes ...string) *Element {
	e.Add("class", classes...)
	return e
}

func (e *Element) AttrID(id string) *Element {
	e.Add("id", id)
	return e
}

func (e *Element) AttrName(name string) *Element {
	e.Add("name", name)
	return e
}

func (e *Element) AttrType(t string) *Element {
	e.Add("type", t)
	return e
}

func (e *Element) AttrValue(value string) *Element {
	e.Add("value", value)
	return e
}

func (e *Element) AttrPlaceholder(placeholder string) *Element {
	e.Add("placeholder", placeholder)
	return e
}

func (e *Element) AttrDisabled(f bool) *Element {
	e.SetBool("disabled", f)
	return e
}

func (e *Element) AttrChecked(f bool) *Element {
	e.SetBool("checked", f)
	return e
}

func (e *Element) AttrRequired(f bool) *Element {
	e.SetBool("required", f)
	return e
}

func (e *Element) AttrReadOnly(f bool) *Element {
	e.SetBool("readonly", f)
	return e
}

func (e *Element) AttrMultiple(f bool) *Element {
	e.SetBool("multiple", f)
	return e
}

func (e *Element) AttrSelected(f bool) *Element {
	e.SetBool("selected", f)
	return e
}

func (e *Element) AttrAutofocus(f bool) *Element {
	e.SetBool("autofocus", f)
	return e
}

func (e *Element) AttrAutocomplete(f bool) *Element {
	e.SetBool("autocomplete", f)
	return e
}

func (e *Element) AttrAutocapitalize(f bool) *Element {
	e.SetBool("autocapitalize", f)
	return e
}

func (e *Element) AttrSpellcheck(f bool) *Element {
	e.SetBool("spellcheck", f)
	return e
}

func (e *Element) AttrNovalidate(f bool) *Element {
	e.SetBool("novalidate", f)
	return e
}

func (e *Element) AttrHidden(f bool) *Element {
	e.SetBool("hidden", f)
	return e
}

func (e *Element) AttrAsync(f bool) *Element {
	e.SetBool("async", f)
	return e
}

func (e *Element) AttrDefer(f bool) *Element {
	e.SetBool("defer", f)
	return e
}

func (e *Element) AttrAutoplay(f bool) *Element {
	e.SetBool("autoplay", f)
	return e
}

func (e *Element) AttrControls(f bool) *Element {
	e.SetBool("controls", f)
	return e
}

func (e *Element) AttrLoop(f bool) *Element {
	e.SetBool("loop", f)
	return e
}

func (e *Element) AttrMuted(f bool) *Element {
	e.SetBool("muted", f)
	return e
}

func (e *Element) AttrDefault(f bool) *Element {
	e.SetBool("default", f)
	return e
}

func (e *Element) AttrOpen(f bool) *Element {
	e.SetBool("open", f)
	return e
}

func (e *Element) AttrScoped(f bool) *Element {
	e.SetBool("scoped", f)
	return e
}
