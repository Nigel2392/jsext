//go:build js && wasm
// +build js,wasm

package elements

import (
	"strconv"
	"time"

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

func fadeIn(e *Element, timeMS int) {
	var transition = e.value.Get("style").Get("transition").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("opacity", "0")
	e.value.Get("classList").Call("add", FADEIN_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("opacity", "1")
		time.Sleep(time.Duration(timeMS) * time.Millisecond)
		e.value.Get("style").Set("transition", transition)
	})
}

func fadeOut(e *Element, timeMS int) {
	var transition = e.value.Get("style").Get("transition").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("opacity", "1")
	e.value.Get("classList").Call("add", FADEOUT_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("opacity", "0")
		time.Sleep(time.Duration(timeMS) * time.Millisecond)
		e.value.Get("style").Set("transition", transition)
	})
}

func fromTop(e *Element, timeMS int) {
	var transform = e.value.Get("style").Get("transform").String()
	var transition = e.value.Get("style").Get("transition").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateY(-100%)")
	e.value.Get("classList").Call("add", FROMTOP_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", transform)
		time.Sleep(time.Duration(timeMS) * time.Millisecond)
		e.value.Get("style").Set("transition", transition)
	})
}

func fromLeft(e *Element, timeMS int) {
	var transform = e.value.Get("style").Get("transform").String()
	var transition = e.value.Get("style").Get("transition").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateX(-100%)")
	e.value.Get("classList").Call("add", FROMLEFT_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", transform)
		time.Sleep(time.Duration(timeMS) * time.Millisecond)
		e.value.Get("style").Set("transition", transition)
	})
}

func fromRight(e *Element, timeMS int) {
	var transform = e.value.Get("style").Get("transform").String()
	var transition = e.value.Get("style").Get("transition").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateX(100%)")
	e.value.Get("classList").Call("add", FROMRIGHT_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", transform)
		time.Sleep(time.Duration(timeMS) * time.Millisecond)
		e.value.Get("style").Set("transition", transition)
	})

}

func fromBottom(e *Element, timeMS int) {
	var transform = e.value.Get("style").Get("transform").String()
	var transition = e.value.Get("style").Get("transition").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateY(100%)")
	e.value.Get("classList").Call("add", FROMBOTTOM_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", "translateY(0)")
		e.value.Get("style").Set("transform", transform)
		time.Sleep(time.Duration(timeMS) * time.Millisecond)
		e.value.Get("style").Set("transition", transition)
	})
}

func bounce(e *Element, timeMS int) {
	var transition = e.value.Get("style").Get("transition").String()
	var transform = e.value.Get("style").Get("transform").String()
	e.value.Get("style").Set("transition", "all "+strconv.Itoa(timeMS/2)+"ms ease-in-out")
	e.value.Get("classList").Call("add", BOUNCE_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", "scale(1.1)")
		go func() {
			time.Sleep(time.Duration(timeMS) * time.Millisecond / 2)
			e.value.Get("style").Set("transform", "scale(1)")
			time.Sleep(time.Duration(timeMS) * time.Millisecond)
			e.value.Get("style").Set("transition", transition)
			e.value.Get("style").Set("transform", transform)
		}()
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