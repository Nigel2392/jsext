//go:build js && wasm
// +build js,wasm

package elements

import (
	"fmt"

	"github.com/Nigel2392/jsext"
)

// All animations get rendered in a separate goroutine.

// Predefined element animations.
type Animation struct {
	Animations     []any
	Options        map[string]interface{}
	WhenInViewport bool
}

// Predefined element animations.
var (
	FADEIN = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"opacity": "0", "offset": "0"},
		map[string]interface{}{"opacity": "1", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
	}}
	FADEOUT = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"opacity": "1", "offset": "0"},
		map[string]interface{}{"opacity": "0", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
	}}
	BOUNCE = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"transform": "scale(1)", "offset": "0"},
		map[string]interface{}{"transform": "scale(1.1)", "offset": "0.2"},
		map[string]interface{}{"transform": "scale(0.9)", "offset": "0.4"},
		map[string]interface{}{"transform": "scale(1.05)", "offset": "0.6"},
		map[string]interface{}{"transform": "scale(0.95)", "offset": "0.8"},
		map[string]interface{}{"transform": "scale(1)", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
	}}
	FROMTOP = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"transform": "translateY(-100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateY(0)", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}}
	FROMLEFT = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"transform": "translateX(-100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateX(0)", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}}
	FROMRIGHT = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"transform": "translateX(100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateX(0)", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}}
	FROMBOTTOM = Animation{WhenInViewport: true, Animations: []any{
		map[string]interface{}{"transform": "translateY(100%)", "offset": "0"},
		map[string]interface{}{"transform": "translateY(0)", "offset": "1"},
	}, Options: map[string]interface{}{
		"duration":   500,
		"iterations": 1,
		"fill":       "forwards",
		"easing":     "ease-in",
	}}
)

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
		go e.runAnimation(anim)
	}
}

func (e *Element) Animate(a Animation) {
	if e.value.IsUndefined() {
		e.animations = append(e.animations, a)
		return
	}
	e.runAnimation(a)
}

func (e *Element) runAnimation(a Animation) {
	var jsArr = jsext.SliceToArray(a.Animations)
	var jsOpts = jsext.MapToObject(a.Options)
	if a.WhenInViewport {
		InViewListener(e, func(this jsext.Value, event jsext.Event) {
			fmt.Println(a, "is in viewport")
			e.value.Call("animate", jsArr.Value(), jsOpts.Value())
		})
	} else {
		fmt.Println(a, "is not in viewport")
		e.value.Call("animate", jsArr.Value(), jsOpts.Value())
	}
}

func (e *Element) addAnim(a Animation, timeMS int) {
	if timeMS > 0 {
		a.Options["duration"] = timeMS
	}
	e.Animate(a)
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
