//go:build js && wasm
// +build js,wasm

package elements

import (
	"github.com/Nigel2392/jsext"
)

// All animations get rendered in a separate goroutine.
// This however does mean that there can be conflicts if multiple animations are added to the same element.

// Predefined element animations.
const (
	FADEIN_CLASS     = "jsext-fade-in"
	FADEOUT_CLASS    = "jsext-fade-out"
	BOUNCE_CLASS     = "jsext-bounce"
	FROMTOP_CLASS    = "jsext-from-top"
	FROMLEFT_CLASS   = "jsext-from-left"
	FROMRIGHT_CLASS  = "jsext-from-right"
	FROMBOTTOM_CLASS = "jsext-from-bottom"
)

// Predefined element animations.
type Animation struct {
	Type     int
	Duration int
}

// Predefined element animations.
const (
	FADEIN     = 0
	FADEOUT    = 1
	BOUNCE     = 2
	FROMTOP    = 3
	FROMLEFT   = 4
	FROMRIGHT  = 5
	FROMBOTTOM = 6
)

// Predefined element animations.
var AnimationMap = map[int]func(e *Element, timeMS int){
	FADEIN:     fadeIn,
	FADEOUT:    fadeOut,
	BOUNCE:     bounce,
	FROMTOP:    fromTop,
	FROMLEFT:   fromLeft,
	FROMRIGHT:  fromRight,
	FROMBOTTOM: fromBottom,
}

// Fade the element in once it is visible on screen.
func (e *Element) FadeIn(timeMS int) *Element {
	e.addAnim(FADEIN, timeMS)
	return e
}

// Fade the element out once it is visible on screen.
func (e *Element) FadeOut(timeMS int) *Element {
	e.addAnim(FADEOUT, timeMS)
	return e
}

// Bounce the element once it is visible on screen.
func (e *Element) Bounce(timeMS int) *Element {
	e.addAnim(BOUNCE, timeMS)
	return e
}

// Slide the element in from the top once it is visible on screen.
func (e *Element) FromTop(timeMS int) *Element {
	e.addAnim(FROMTOP, timeMS)
	return e
}

// Slide the element in from the left once it is visible on screen.
func (e *Element) FromLeft(timeMS int) *Element {
	e.addAnim(FROMLEFT, timeMS)
	return e
}

// Slide the element in from the right once it is visible on screen.
func (e *Element) FromRight(timeMS int) *Element {
	e.addAnim(FROMRIGHT, timeMS)
	return e
}

// Slide the element in from the bottom once it is visible on screen.
func (e *Element) FromBottom(timeMS int) *Element {
	e.addAnim(FROMBOTTOM, timeMS)
	return e
}

func (e *Element) animate() {
	for _, anim := range e.animations {
		var f, ok = AnimationMap[anim.Type]
		if ok {
			go f(e, anim.Duration)
		}
	}
}

func (e *Element) addAnim(typ, duration int) {
	e.animations = append(e.animations, Animation{Type: typ, Duration: duration})
	if !e.value.IsUndefined() {
		e.animate()
	}
}

func (e *Element) Animate(animations []any, opts map[string]interface{}) *Element {
	var jsArr = jsext.SliceToArray(animations)
	var jsOpts = jsext.MapToObject(opts)
	e.value.Call("animate", jsArr.Value(), jsOpts.Value())
	return e
}

func fadeIn(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"opacity": "0", "offset": "0"},
		map[string]interface{}{"opacity": "1", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
	}

	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func fadeOut(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"opacity": "1", "offset": "0"},
		map[string]interface{}{"opacity": "0", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
	}
	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func fromTop(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"transform": "translateY(-100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateY(0)", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}
	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func fromLeft(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"transform": "translateX(-100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateX(0)", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}
	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func fromRight(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"transform": "translateX(100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateX(0)", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}
	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func fromBottom(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"transform": "translateY(100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateY(0)", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}
	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func bounce(e *Element, timeMS int) {
	var arr = []any{
		map[string]interface{}{"transform": "scale(1)", "offset": "0"},
		map[string]interface{}{"transform": "scale(1.1)", "offset": "0.2"},
		map[string]interface{}{"transform": "scale(0.9)", "offset": "0.4"},
		map[string]interface{}{"transform": "scale(1.05)", "offset": "0.6"},
		map[string]interface{}{"transform": "scale(0.95)", "offset": "0.8"},
		map[string]interface{}{"transform": "scale(1)", "offset": "1"},
	}
	var opts = map[string]interface{}{
		"duration":   timeMS,
		"iterations": 1,
		"fill":       "forwards",
	}
	InViewListenerSingleExecution(e, func(this jsext.Value, event jsext.Event) {
		e.Animate(arr, opts)
	})
}

func InViewListener(e *Element, cb func(this jsext.Value, event jsext.Event)) {
	if isInViewport(e) {
		cb(jsext.Value{}, jsext.Event{})
	}
	jsext.Element(jsext.Window).AddEventListener("scroll", func(this jsext.Value, event jsext.Event) {
		if isInViewport(e) {
			cb(this, event)
		}
	})
}

func InViewListenerSingleExecution(e *Element, cb func(this jsext.Value, event jsext.Event)) {
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
