//go:build js && wasm
// +build js,wasm

package jsext

import (
	"strings"
	"syscall/js"
)

// Imports keeps track of all imported files.
var Imported = make(map[string]Import)

// Remove an import from the DOM, and remove it from the Imported map.
func RemoveImport(name string) {
	var imp = Imported[name]
	if imp.name == "" {
		return
	}
	imp.jsVal.Remove()
	delete(Imported, name)
}

// Import is a file that has been imported into the DOM.
type Import struct {
	name  string
	jsVal Element
}

// MarshalJS returns the underlying js.Value.
func (e Import) MarshalJS() js.Value {
	return e.jsVal.MarshalJS()
}

// Value returns the jsext.Value of the import.
func (i Import) Value() Value {
	return Value(i.jsVal)
}

// JSValue returns the js.Value of the import.
func (i Import) JSValue() js.Value {
	return js.Value(i.jsVal)
}

// Import a link
func ImportLink(name, src, typ, rel string) Import {
	var link = CreateElement("link")
	link.Set("type", typ)
	if rel == "" {
		rel = "stylesheet"
	}
	link.Set("rel", rel)
	link.Set("href", src)
	var i = Import{
		name,
		link,
	}
	i.run()
	return i
}

// Import a style block for use with raw sourcecode.
func StyleBlock(name, code string) Import {
	var style = CreateElement("style")
	style.Set("type", "text/css")
	style.Set("innerHTML", code)
	var i = Import{
		name,
		style,
	}
	i.run()
	return i
}

// Import a single css block into the Global jsext style element.
func ImportCSS(name, selector string, src ...string) Import {
	var css = strings.Join(src, ";")
	var style Element
	var err error
	style, err = GetElementById("jsext-main-style")
	if err == nil {
		style.Set("innerHTML", style.Get("innerHTML").String()+"\n"+selector+" { "+css+" }")
		return Import{
			name,
			style,
		}
	} else {
		style = CreateElement("style")
	}
	style.Set("type", "text/css")
	style.Set("id", "jsext-main-style")
	style.Set("innerHTML", selector+" { "+css+" }")
	var i = Import{
		name,
		style,
	}
	i.run()
	return i
}

// Import a script
func ImportScript(name, src, typ string) Import {
	var script = CreateElement("script")
	script.Set("type", typ)
	script.Set("src", src)
	var i = Import{
		name,
		script,
	}
	i.run()
	return i
}

// Import a script block for use with raw sourcecode.
func ScriptBlock(name, code string) Import {
	var script = CreateElement("script")
	script.Set("type", "text/javascript")
	script.Set("text", code)
	var i = Import{
		name,
		script,
	}
	i.run()
	return i
}

func (i Import) run() {
	imp, ok := Imported[i.name]
	if ok {
		imp.Remove()
	}
	Imported[i.name] = i
	Head.AppendChild(i.jsVal)
}

// Remove an import.
func (i Import) Remove() {
	i.jsVal.Remove()
	delete(Imported, i.name)
}
