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

func (e *Element) FadeIn(timeMS int) {
	e.Value().Get("style").Set("transition", "opacity "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("opacity", "0")
	e.Value().Get("classList").Call("add", FADEIN_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("opacity", "1")
	})
}

func (e *Element) FadeOut(timeMS int) {
	e.Value().Get("style").Set("transition", "opacity "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("opacity", "1")
	e.Value().Get("classList").Call("add", FADEOUT_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("opacity", "0")
	})
}

func (e *Element) FromTop(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateY(-100%)")
	e.Value().Get("classList").Call("add", FROMTOP_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", translate)
	})
}

func (e *Element) FromLeft(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateX(-100%)")
	e.Value().Get("classList").Call("add", FROMLEFT_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", translate)
	})
}

func (e *Element) FromRight(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateX(100%)")
	e.Value().Get("classList").Call("add", FROMRIGHT_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", translate)
	})

}

func (e *Element) FromBottom(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateY(100%)")
	e.Value().Get("classList").Call("add", FROMBOTTOM_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", "translateY(0)")
		e.Value().Get("style").Set("transform", translate)
	})
}

func (e *Element) Bounce(timeMS int) {
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("classList").Call("add", BOUNCE_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", "scale(1.1)")
	})
}

func (e *Element) InViewListener(cb func(this jsext.Value, event jsext.Event)) {
	if e.isInViewport() {
		cb(jsext.Value{}, jsext.Event{})
	}
	jsext.Element(jsext.Window).AddEventListener("scroll", func(this jsext.Value, event jsext.Event) {
		if e.isInViewport() {
			cb(this, event)
		}
	})
}

func (e *Element) isInViewport() bool {
	var (
		bounding   = e.Value().Call("getBoundingClientRect")
		elemHeight = e.Value().Get("offsetHeight").Int()
		elemWidth  = e.Value().Get("offsetWidth").Int()
	)

	if bounding.Get("top").Int() >= -elemHeight &&
		bounding.Get("left").Int() >= -elemWidth &&
		bounding.Get("bottom").Int() <= (jsext.Window.Get("innerHeight").Int()+elemHeight) &&
		bounding.Get("right").Int() <= (jsext.Window.Get("innerWidth").Int()+elemWidth) {
		return true
	}

	return false
}

func (e *Element) EaseIn(timeMS int) {
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms ease")
}

func (e *Element) EaseOut(timeMS int) {
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms ease-out")
}

func (e *Element) EaseInOut(timeMS int) {
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms ease-in-out")
}

func (e *Element) EaseInBack(timeMS int) {
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms cubic-bezier(0.600, -0.280, 0.735, 0.045)")
}

func (e *Element) EaseOutBack(timeMS int) {
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms cubic-bezier(0.175, 0.885, 0.320, 1.275)")
}
