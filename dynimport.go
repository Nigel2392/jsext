package jsext

import (
	"syscall/js"
)

var Imported = make(map[string]Import)

func RemoveImport(name string) {
	var imp = Imported[name]
	if imp.name == "" {
		return
	}
	imp.jsVal.Remove()
	delete(Imported, name)
}

type Import struct {
	name  string
	jsVal Element
}

func (i Import) Value() Value {
	return Value(i.jsVal)
}

func (i Import) JSValue() js.Value {
	return js.Value(i.jsVal)
}

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

func StyleBlock(name, code string) Import {
	var style = CreateElement("style")
	style.Set("type", "text/css")
	style.Set("text", code)
	var i = Import{
		name,
		style,
	}
	i.run()
	return i
}

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

func (i Import) Remove() {
	i.jsVal.Remove()
	delete(Imported, i.name)
}
