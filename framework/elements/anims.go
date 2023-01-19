//go:build js && wasm
// +build js,wasm

package elements

import (
	"syscall/js"

	"github.com/Nigel2392/jsext"
)

// All animations get rendered in a separate goroutine.
const Infinity = "Infinity"

var scaleUp = Animation{Animations: []any{
	map[string]interface{}{"transform": "scale(0)", "offset": "0"},
	map[string]interface{}{"transform": "scale(1)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var scaleDown = Animation{Animations: []any{
	map[string]interface{}{"transform": "scale(1)", "offset": "0"},
	map[string]interface{}{"transform": "scale(0)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var buzz = Animation{Animations: []any{
	map[string]interface{}{"transform": "rotate(0deg)", "offset": "0"},
	map[string]interface{}{"transform": "rotate(5deg)", "offset": "0.1"},
	map[string]interface{}{"transform": "rotate(-5deg)", "offset": "0.2"},
	map[string]interface{}{"transform": "rotate(5deg)", "offset": "0.3"},
	map[string]interface{}{"transform": "rotate(-5deg)", "offset": "0.4"},
	map[string]interface{}{"transform": "rotate(5deg)", "offset": "0.5"},
	map[string]interface{}{"transform": "rotate(-5deg)", "offset": "0.6"},
	map[string]interface{}{"transform": "rotate(5deg)", "offset": "0.7"},
	map[string]interface{}{"transform": "rotate(-5deg)", "offset": "0.8"},
	map[string]interface{}{"transform": "rotate(5deg)", "offset": "0.9"},
	map[string]interface{}{"transform": "rotate(0deg)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var shake = Animation{Animations: []any{
	map[string]interface{}{"transform": "translate(0, 0)", "offset": "0"},
	map[string]interface{}{"transform": "translate(-10px, 0)", "offset": "0.1"},
	map[string]interface{}{"transform": "translate(10px, 0)", "offset": "0.2"},
	map[string]interface{}{"transform": "translate(-10px, 0)", "offset": "0.3"},
	map[string]interface{}{"transform": "translate(10px, 0)", "offset": "0.4"},
	map[string]interface{}{"transform": "translate(-10px, 0)", "offset": "0.5"},
	map[string]interface{}{"transform": "translate(10px, 0)", "offset": "0.6"},
	map[string]interface{}{"transform": "translate(-10px, 0)", "offset": "0.7"},
	map[string]interface{}{"transform": "translate(10px, 0)", "offset": "0.8"},
	map[string]interface{}{"transform": "translate(-10px, 0)", "offset": "0.9"},
	map[string]interface{}{"transform": "translate(0, 0)", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
var flash = Animation{Animations: []any{
	map[string]interface{}{"opacity": "1", "offset": "0"},
	map[string]interface{}{"opacity": "0", "offset": "0.25"},
	map[string]interface{}{"opacity": "1", "offset": "0.5"},
	map[string]interface{}{"opacity": "0", "offset": "0.75"},
	map[string]interface{}{"opacity": "1", "offset": "1"},
}, Options: map[string]interface{}{
	"iterations": 1,
	"fill":       "forwards",
}}
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

type Animations struct {
	animations []Animation
	element    *Element
}

func (a *Animations) Add(anim Animation) *Animations {
	a.animations = append(a.animations, anim)
	return a
}

// Predefined element animations.
type Animation struct {
	Animations             []any
	Options                map[string]interface{}
	whenInViewportAndReset bool
	ResetWhenLeaveViewport bool
}

func (a *Animations) Element() *Element {
	return a.element
}

func (a *Animations) Rainbow(colorsPerSecond float64, colors ...string) *Animations {
	if len(colors) == 0 {
		colors = []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet"}
	}
	a.element.AttrStyle("color:" + colors[0])
	var anim = Animation{whenInViewportAndReset: true, Animations: make([]any, len(colors)+1), Options: map[string]interface{}{
		"duration":   1000 / colorsPerSecond * float64(len(colors)),
		"iterations": Infinity,
	}}
	for i, color := range colors {
		anim.Animations[i] = map[string]interface{}{"color": color, "offset": float64(i) / float64(len(colors))}
	}
	a.Animate(anim)
	return a
}

// Scale the element up
func (a *Animations) ScaleUp(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = scaleUp
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Scale the element down
func (a *Animations) ScaleDown(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = scaleDown
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Buzz the element
func (a *Animations) Buzz(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = buzz
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Shake the element
func (a *Animations) Shake(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = shake
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Flash the element
func (a *Animations) Flash(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = flash
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Fade the element in
func (a *Animations) FadeIn(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = fadeIn
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Fade the element out
func (a *Animations) FadeOut(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = fadeOut
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Bounce the element
func (a *Animations) Bounce(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = bounce
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Slide the element in from the top
func (a *Animations) FromTop(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = fromTop
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Slide the element in from the left
func (a *Animations) FromLeft(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = fromLeft
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Slide the element in from the right
func (a *Animations) FromRight(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = fromRight
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

// Slide the element in from the bottom
func (a *Animations) FromBottom(timeMS int, whenInViewportAndReset ...bool) *Animations {
	var anim = fromBottom
	anim.Options["duration"] = timeMS
	anim.whenInViewportAndReset = len(whenInViewportAndReset) > 0 && whenInViewportAndReset[0]
	anim.ResetWhenLeaveViewport = len(whenInViewportAndReset) > 1 && whenInViewportAndReset[1]
	a.Animate(anim)
	return a
}

func (a *Animations) animate() {
	for _, anim := range a.animations {
		go a.Animate(anim)
	}
}

func (a *Animations) Animate(anim Animation) {
	if a.element.value.IsUndefined() || a.element.value.IsNull() {
		a.animations = append(a.animations, anim)
		return
	}
	var jsArr = jsext.SliceToArray(anim.Animations)
	var jsOpts = jsext.MapToObject(anim.Options)
	if anim.whenInViewportAndReset {
		var observer = jsext.Get("IntersectionObserver")
		var observerInstance jsext.Value
		observerInstance = observer.New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			var entries = args[0]
			var entry = entries.Index(0)
			var isIntersecting = entry.Get("isIntersecting").Bool()
			if isIntersecting {
				a.element.value.Call("animate", jsArr.Value(), jsOpts.Value())
				if !anim.ResetWhenLeaveViewport {
					observerInstance.Call("unobserve", entry.Get("target"))
				}
			}
			return nil
		}), jsext.MapToObject(map[string]interface{}{
			"root":       nil,
			"rootMargin": "0px",
			"threshold":  0.5,
		}).Value())

		observerInstance.Call("observe", a.element.value)

	} else {
		a.element.value.Call("animate", jsArr.Value(), jsOpts.Value())
	}
}

//	func InViewListener(e *Element, cb func(this jsext.Value, event jsext.Event), resetOnLeave bool) {
//		var ran = false
//		if isInViewport(e) {
//			if !ran {
//				cb(jsext.Value{}, jsext.Event{})
//				ran = true
//			}
//		}
//		jsext.Element(jsext.Window).AddEventListener("scroll", func(this jsext.Value, event jsext.Event) {
//			var isInView = isInViewport(e)
//			if !ran && isInView {
//				cb(this, event)
//				ran = true
//			} else if ran && !isInView && resetOnLeave {
//				ran = false
//			}
//		})
//	}
//
//	// isInViewport checks if the element is in the viewport
//	func isInViewport(e *Element) bool {
//		var (
//			bounding   = e.value.Call("getBoundingClientRect")
//			elemHeight = e.value.Get("offsetHeight").Int()
//			elemWidth  = e.value.Get("offsetWidth").Int()
//		)
//
//		if bounding.Get("top").Int() >= -elemHeight &&
//			bounding.Get("left").Int() >= -elemWidth &&
//			bounding.Get("bottom").Int() <= (jsext.Window.Get("innerHeight").Int()+elemHeight) &&
//			bounding.Get("right").Int() <= (jsext.Window.Get("innerWidth").Int()+elemWidth) {
//			return true
//		}
//
//		return false
//	}
//
