package jse

func (e *Element) Attrs(m map[string]interface{}) *Element {
	for k, v := range m {
		e.Set(k, v)
	}
	return e
}

func (e *Element) AttrAccept(v string) *Element {
	e.SetAttr("accept", v)
	return e
}

func (e *Element) AttrAcceptCharset(v string) *Element {
	e.SetAttr("accept-charset", v)
	return e
}

func (e *Element) AttrAccessKey(v string) *Element {
	e.SetAttr("accesskey", v)
	return e
}

func (e *Element) AttrAction(v string) *Element {
	e.SetAttr("action", v)
	return e
}

func (e *Element) AttrAlt(v string) *Element {
	e.SetAttr("alt", v)
	return e
}

func (e *Element) AttrAsync(v string) *Element {
	e.SetAttr("async", v)
	return e
}

func (e *Element) AttrAutocomplete(v string) *Element {
	e.SetAttr("autocomplete", v)
	return e
}

func (e *Element) AttrAutofocus(v string) *Element {
	e.SetAttr("autofocus", v)
	return e
}

func (e *Element) AttrAutoplay(v string) *Element {
	e.SetAttr("autoplay", v)
	return e
}

func (e *Element) AttrCharset(v string) *Element {
	e.SetAttr("charset", v)
	return e
}

func (e *Element) AttrChecked(v bool) *Element {
	e.Set("checked", v)
	return e
}

func (e *Element) AttrCite(v string) *Element {
	e.SetAttr("cite", v)
	return e
}

func (e *Element) AttrID(v string) *Element {
	e.SetAttr("id", v)
	return e
}

func (e *Element) AttrType(v string) *Element {
	e.SetAttr("type", v)
	return e
}

func (e *Element) AttrName(v string) *Element {
	e.SetAttr("name", v)
	return e
}

func (e *Element) AttrValue(v string) *Element {
	e.SetAttr("value", v)
	return e
}

func (e *Element) AttrPlaceholder(v string) *Element {
	e.SetAttr("placeholder", v)
	return e
}

func (e *Element) AttrSrc(v string) *Element {
	e.SetAttr("src", v)
	return e
}

func (e *Element) AttrHref(v string) *Element {
	e.SetAttr("href", v)
	return e
}

func (e *Element) AttrDisabled(v bool) *Element {
	e.Set("disabled", v)
	return e
}

func (e *Element) AttrDownload(v string) *Element {
	e.SetAttr("download", v)
	return e
}

func (e *Element) AttrEncType(v string) *Element {
	e.SetAttr("enctype", v)
	return e
}

func (e *Element) AttrFor(v string) *Element {
	e.SetAttr("for", v)
	return e
}

func (e *Element) AttrReadOnly(v bool) *Element {
	e.Set("readonly", v)
	return e
}

func (e *Element) AttrRequired(v bool) *Element {
	e.Set("required", v)
	return e
}

func (e *Element) AttrSelected(v bool) *Element {
	e.Set("selected", v)
	return e
}

func (e *Element) AttrWidth(v string) *Element {
	e.SetAttr("width", v)
	return e
}

func (e *Element) AttrHeight(v string) *Element {
	e.SetAttr("height", v)
	return e
}

func (e *FormElement) AttrID(v string) *FormElement {
	e.Element().SetAttr("id", v)
	return e
}

func (e *FormElement) AttrType(v string) *FormElement {
	e.Element().SetAttr("type", v)
	return e
}

func (e *FormElement) AttrName(v string) *FormElement {
	e.Element().SetAttr("name", v)
	return e
}

func (e *FormElement) AttrValue(v string) *FormElement {
	e.Element().SetAttr("value", v)
	return e
}

func (e *FormElement) AttrPlaceholder(v string) *FormElement {
	e.Element().SetAttr("placeholder", v)
	return e
}

func (e *FormElement) AttrDisabled(v bool) *FormElement {
	e.Element().Set("disabled", v)
	return e
}

func (e *FormElement) AttrFor(v string) *FormElement {
	e.Element().SetAttr("for", v)
	return e
}

func (e *FormElement) AttrReadOnly(v bool) *FormElement {
	e.Element().Set("readonly", v)
	return e
}

func (e *FormElement) AttrRequired(v bool) *FormElement {
	e.Element().Set("required", v)
	return e
}
