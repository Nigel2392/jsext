//go:build js && wasm
// +build js,wasm

package elements

import (
	"github.com/Nigel2392/jsext"
)

// All animations get rendered in a separate goroutine.
const Infinity = "Infinity"

var fadeIn = Animation{Animations: []any{
	map[string]interface{}{"opacity": "0", "offset": "0"},
	map[string]interface{}{"opacity": "1", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var fadeOut = Animation{Animations: []any{
	map[string]interface{}{"opacity": "1", "offset": "0"},
	map[string]interface{}{"opacity": "0", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var bounce = Animation{Animations: []any{
	map[string]interface{}{"transform": "scale(1)", "offset": "0"},
	map[string]interface{}{"transform": "scale(1.1)", "offset": "0.2"},
	map[string]interface{}{"transform": "scale(0.9)", "offset": "0.4"},
	map[string]interface{}{"transform": "scale(1.05)", "offset": "0.6"},
	map[string]interface{}{"transform": "scale(0.95)", "offset": "0.8"},
	map[string]interface{}{"transform": "scale(1)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var fromTop = Animation{Animations: []any{
	map[string]interface{}{"transform": "translateY(-100%)", "offset": "0"},
	map[string]interface{}{"transform": "translateY(0)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
	"easing":     "ease-in",
}}
var fromLeft = Animation{Animations: []any{
	map[string]interface{}{"transform": "translateX(-100%)", "offset": "0"},
	map[string]interface{}{"transform": "translateX(0)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
	"easing":     "ease-in",
}}
var fromRight = Animation{Animations: []any{
	map[string]interface{}{"transform": "translateX(100%)", "offset": "0"},
	map[string]interface{}{"transform": "translateX(0)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
	"easing":     "ease-in",
}}
var fromBottom = Animation{Animations: []any{
	map[string]interface{}{"transform": "translateY(100%)", "offset": "0"},
	map[string]interface{}{"transform": "translateY(0)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
	"easing":     "ease-in",
}}

// Predefined element animations.
type Animation struct {
	Animations     []any
	Options        map[string]interface{}
	WhenInViewport bool
}

func (e *Element) Rainbow(colorsPerSecond float64, colors ...string) *Element {
	if len(colors) == 0 {
		colors = []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"}
	}
	e.AttrStyle("color:" + colors[0])
	var anim = Animation{WhenInViewport: true, Animations: make([]any, len(colors)+1), Options: map[string]interface{}{
		"duration":   1000 / colorsPerSecond * float64(len(colors)),
		"iterations": Infinity,
	}}
	for i, color := range colors {
		anim.Animations[i] = map[string]interface{}{"color": color, "offset": float64(i) / float64(len(colors))}
	}
	e.Animate(anim)
	return e
}

// Fade the element in
func (e *Element) FadeIn(timeMS int, wheninViewport ...bool) *Element {
	var anim = fadeIn
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

// Fade the element out
func (e *Element) FadeOut(timeMS int, wheninViewport ...bool) *Element {
	var anim = fadeOut
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

// Bounce the element
func (e *Element) Bounce(timeMS int, wheninViewport ...bool) *Element {
	var anim = bounce
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

// Slide the element in from the top
func (e *Element) FromTop(timeMS int, wheninViewport ...bool) *Element {
	var anim = fromTop
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

// Slide the element in from the left
func (e *Element) FromLeft(timeMS int, wheninViewport ...bool) *Element {
	var anim = fromLeft
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

// Slide the element in from the right
func (e *Element) FromRight(timeMS int, wheninViewport ...bool) *Element {
	var anim = fromRight
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

// Slide the element in from the bottom
func (e *Element) FromBottom(timeMS int, wheninViewport ...bool) *Element {
	var anim = fromBottom
	anim.Options["duration"] = timeMS
	anim.WhenInViewport = len(wheninViewport) > 0 && wheninViewport[0]
	e.Animate(anim)
	return e
}

func (e *Element) animate() {
	for _, anim := range e.animations {
		go e.Animate(anim)
	}
}

func (e *Element) Animate(a Animation) {
	if e.value.IsUndefined() || e.value.IsNull() {
		e.animations = append(e.animations, a)
		return
	}
	var jsArr = jsext.SliceToArray(a.Animations)
	var jsOpts = jsext.MapToObject(a.Options)
	if a.WhenInViewport {
		InViewListener(e, func(this jsext.Value, event jsext.Event) {
			e.value.Call("animate", jsArr.Value(), jsOpts.Value())
		})
	} else {
		e.value.Call("animate", jsArr.Value(), jsOpts.Value())
	}
}

func InViewListener(e *Element, cb func(this jsext.Value, event jsext.Event)) {
	var ran = false
	if isInViewport(e) {
		if !ran {
			cb(jsext.Value{}, jsext.Event{})
			ran = true
		}
	}
	jsext.Element(jsext.Window).AddEventListener("scroll", func(this jsext.Value, event jsext.Event) {
		if !ran {
			if isInViewport(e) {
				cb(this, event)
				ran = true
			}
		}
	})
}

// isInViewport checks if the element is in the viewport
func isInViewport(e *Element) bool {
	var (
		bounding   = e.value.Call("getBoundingClientRect")
		elemHeight = e.value.Get("offsetHeight").Int()
		elemWidth  = e.value.Get("offsetWidth").Int()
	)

	if bounding.Get("top").Int() >= -elemHeight &&
		bounding.Get("left").Int() >= -elemWidth &&
		bounding.Get("bottom").Int() <= (jsext.Window.Get("innerHeight").Int()+elemHeight) &&
		bounding.Get("right").Int() <= (jsext.Window.Get("innerWidth").Int()+elemWidth) {
		return true
	}

	return false
}
