package shortcuts

import (
	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
)

const (
	F_BIND_HTML_HTML uint32 = 1 << iota
	F_BIND_VALUE_VALUE
	F_BIND_HTML_VALUE
	F_BIND_VALUE_HTML
	F_BIND_CLICK
	F_PREVENT_DEFAULT
	F_STOP_PROPAGATION
	F_STOP_IMMEDIATE_PROPAGATION
)

func BindElements(src *jse.Element, event string, flags uint32, dest ...*jse.Element) {
	src.AddEventListener(event, func(this *jse.Element, event jsext.Event) {
		if flags&F_PREVENT_DEFAULT != 0 {
			event.PreventDefault()
		}
		if flags&F_STOP_PROPAGATION != 0 {
			event.StopPropagation()
		}
		if flags&F_STOP_IMMEDIATE_PROPAGATION != 0 {
			event.StopImmediatePropagation()
		}
		for _, elem := range dest {
			if elem == nil || elem.IsNull() || elem.IsUndefined() {
				continue
			}
			if flags&F_BIND_HTML_HTML != 0 {
				elem.Set("innerHTML", src.Get("innerHTML").Value())
			}
			if flags&F_BIND_VALUE_VALUE != 0 {
				elem.Set("value", src.Get("value").Value())
			}
			if flags&F_BIND_HTML_VALUE != 0 {
				elem.Set("value", src.Get("innerHTML").Value())
			}
			if flags&F_BIND_VALUE_HTML != 0 {
				elem.Set("innerHTML", src.Get("value").Value())
			}
			if flags&F_BIND_CLICK != 0 {
				elem.Call("click")
			}
		}
	})
}
