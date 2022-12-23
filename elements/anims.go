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

func (e *Element) FadeIn(timeMS int) {
	e.addAnim(FADEIN, timeMS)
}

func (e *Element) FadeOut(timeMS int) {
	e.addAnim(FADEOUT, timeMS)
}

func (e *Element) Bounce(timeMS int) {
	e.addAnim(BOUNCE, timeMS)
}

func (e *Element) FromTop(timeMS int) {
	e.addAnim(FROMTOP, timeMS)
}

func (e *Element) FromLeft(timeMS int) {
	e.addAnim(FROMLEFT, timeMS)
}

func (e *Element) FromRight(timeMS int) {
	e.addAnim(FROMRIGHT, timeMS)
}

func (e *Element) FromBottom(timeMS int) {
	e.addAnim(FROMBOTTOM, timeMS)
}

func (e *Element) animate() {
	for _, anim := range e.animations {
		switch anim.Type {
		case FADEIN:
			e.fadeIn(anim.Duration)
		case FADEOUT:
			e.fadeOut(anim.Duration)
		case FROMTOP:
			e.fromTop(anim.Duration)
		case FROMBOTTOM:
			e.fromBottom(anim.Duration)
		case FROMLEFT:
			e.fromLeft(anim.Duration)
		case FROMRIGHT:
			e.fromRight(anim.Duration)
		case BOUNCE:
			e.bounce(anim.Duration)
		}
	}
}

func (e *Element) addAnim(typ, duration int) {
	e.animations = append(e.animations, Animation{Type: typ, Duration: duration})
	if !e.value.IsUndefined() {
		e.animate()
	}
}

func (e *Element) fadeIn(timeMS int) {
	e.Value().Get("style").Set("transition", "opacity "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("opacity", "0")
	e.Value().Get("classList").Call("add", FADEIN_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("opacity", "1")
	})
}

func (e *Element) fadeOut(timeMS int) {
	e.Value().Get("style").Set("transition", "opacity "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("opacity", "1")
	e.Value().Get("classList").Call("add", FADEOUT_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("opacity", "0")
	})
}

func (e *Element) fromTop(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateY(-100%)")
	e.Value().Get("classList").Call("add", FROMTOP_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", translate)
	})
}

func (e *Element) fromLeft(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateX(-100%)")
	e.Value().Get("classList").Call("add", FROMLEFT_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", translate)
	})
}

func (e *Element) fromRight(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateX(100%)")
	e.Value().Get("classList").Call("add", FROMRIGHT_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", translate)
	})

}

func (e *Element) fromBottom(timeMS int) {
	var translate = e.Value().Get("style").Get("transform").String()
	e.Value().Get("style").Set("transition", "transform "+strconv.Itoa(timeMS)+"ms")
	e.Value().Get("style").Set("transform", "translateY(100%)")
	e.Value().Get("classList").Call("add", FROMBOTTOM_CLASS)
	e.InViewListener(func(this jsext.Value, event jsext.Event) {
		e.Value().Get("style").Set("transform", "translateY(0)")
		e.Value().Get("style").Set("transform", translate)
	})
}

func (e *Element) bounce(timeMS int) {
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
