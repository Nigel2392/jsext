package elements

import (
	"strconv"

	"github.com/Nigel2392/jsext"
)

const (
	FADEIN_CLASS  = "jsext-fade-in"
	FADEOUT_CLASS = "jsext-fade-out"
	BOUNCE_CLASS  = "jsext-bounce"

	FROMTOP_CLASS    = "jsext-from-top"
	FROMLEFT_CLASS   = "jsext-from-left"
	FROMRIGHT_CLASS  = "jsext-from-right"
	FROMBOTTOM_CLASS = "jsext-from-bottom"
)

type Animation struct {
	Type     int
	Duration int
}

const (
	FADEIN     = 0
	FADEOUT    = 1
	BOUNCE     = 2
	FROMTOP    = 3
	FROMLEFT   = 4
	FROMRIGHT  = 5
	FROMBOTTOM = 6
)

var AnimationMap = map[int]func(e *Element, timeMS int){
	FADEIN:     fadeIn,
	FADEOUT:    fadeOut,
	BOUNCE:     bounce,
	FROMTOP:    fromTop,
	FROMLEFT:   fromLeft,
	FROMRIGHT:  fromRight,
	FROMBOTTOM: fromBottom,
}

func (e *Element) FadeIn(timeMS int) *Element {
	e.addAnim(FADEIN, timeMS)
	return e
}

func (e *Element) FadeOut(timeMS int) *Element {
	e.addAnim(FADEOUT, timeMS)
	return e
}

func (e *Element) Bounce(timeMS int) *Element {
	e.addAnim(BOUNCE, timeMS)
	return e
}

func (e *Element) FromTop(timeMS int) *Element {
	e.addAnim(FROMTOP, timeMS)
	return e
}

func (e *Element) FromLeft(timeMS int) *Element {
	e.addAnim(FROMLEFT, timeMS)
	return e
}

func (e *Element) FromRight(timeMS int) *Element {
	e.addAnim(FROMRIGHT, timeMS)
	return e
}

func (e *Element) FromBottom(timeMS int) *Element {
	e.addAnim(FROMBOTTOM, timeMS)
	return e
}

func (e *Element) animate() {
	for _, anim := range e.animations {
		var f, ok = AnimationMap[anim.Type]
		if ok {
			f(e, anim.Duration)
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
	e.value.Get("style").Set("transition", "opacity "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("opacity", "0")
	e.value.Get("classList").Call("add", FADEIN_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("opacity", "1")
	})
}

func fadeOut(e *Element, timeMS int) {
	e.value.Get("style").Set("transition", "opacity "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("opacity", "1")
	e.value.Get("classList").Call("add", FADEOUT_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("opacity", "0")
	})
}

func fromTop(e *Element, timeMS int) {
	var translate = e.value.Get("style").Get("transform").String()
	e.value.Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateY(-100%)")
	e.value.Get("classList").Call("add", FROMTOP_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", translate)
	})
}

func fromLeft(e *Element, timeMS int) {
	var translate = e.value.Get("style").Get("transform").String()
	e.value.Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateX(-100%)")
	e.value.Get("classList").Call("add", FROMLEFT_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", translate)
	})
}

func fromRight(e *Element, timeMS int) {
	var translate = e.value.Get("style").Get("transform").String()
	e.value.Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateX(100%)")
	e.value.Get("classList").Call("add", FROMRIGHT_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", translate)
	})

}

func fromBottom(e *Element, timeMS int) {
	var translate = e.value.Get("style").Get("transform").String()
	e.value.Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("style").Set("transform", "translateY(100%)")
	e.value.Get("classList").Call("add", FROMBOTTOM_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", "translateY(0)")
		e.value.Get("style").Set("transform", translate)
	})
}

func bounce(e *Element, timeMS int) {
	e.value.Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.value.Get("classList").Call("add", BOUNCE_CLASS)
	InViewListener(e, func(this jsext.Value, event jsext.Event) {
		e.value.Get("style").Set("transform", "scale(1.1)")
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
