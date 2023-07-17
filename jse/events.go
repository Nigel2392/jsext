package jse

import (
	"syscall/js"
	"time"

	"github.com/Nigel2392/jsext/v2"
)

// OnClick adds an event listener to the Element
func (e *Element) OnClick(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("click", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnChange adds an event listener to the Element
func (e *Element) OnChange(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("change", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnKeyUp adds an event listener to the Element
func (e *Element) OnKeyUp(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("keyup", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnKeyDown adds an event listener to the Element
func (e *Element) OnKeyDown(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("keydown", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnKeyPress adds an event listener to the Element
func (e *Element) OnKeyPress(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("keypress", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnFocus adds an event listener to the Element
func (e *Element) OnFocus(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("focus", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnBlur adds an event listener to the Element
func (e *Element) OnBlur(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("blur", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseDown adds an event listener to the Element
func (e *Element) OnMouseDown(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mousedown", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseUp adds an event listener to the Element
func (e *Element) OnMouseUp(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mouseup", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseOver adds an event listener to the Element
func (e *Element) OnMouseOver(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mouseover", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseOut adds an event listener to the Element
func (e *Element) OnMouseOut(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mouseout", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseMove adds an event listener to the Element
func (e *Element) OnMouseMove(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mousemove", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseEnter adds an event listener to the Element
func (e *Element) OnMouseEnter(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mouseenter", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnMouseLeave adds an event listener to the Element
func (e *Element) OnMouseLeave(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("mouseleave", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

// OnScroll adds an event listener to the Element
func (e *Element) OnScroll(callback func(this *Element, event jsext.Event)) js.Func {
	return e.AddEventListener("scroll", func(e *Element, event jsext.Event) {
		callback(e, event)
	})
}

type debouncer struct {
	timer  *time.Timer
	ticker *time.Ticker
}

func (d *debouncer) Stop() {
	if d == nil {
		return
	}
	if d.timer != nil {
		d.timer.Stop()
	}
	if d.ticker != nil {
		d.ticker.Stop()
	}
}

// afterFunc calls f after the duration d.
// It returns a debouncer that can be used to stop the timer.
func afterFunc(d, rep time.Duration, event jsext.Event, f func()) *debouncer {
	var debouncer = &debouncer{
		timer:  time.AfterFunc(d, f),
		ticker: time.NewTicker(rep),
	}
	event.Set("stopTimer", js.FuncOf(func(this js.Value, args []js.Value) any {
		debouncer.Stop()
		return nil
	}))
	return debouncer
}

// OnHoldKey adds an event listener to the Element
//
// wait is the time to wait before the first repeat
// repeat is the time between each function call
//
// Optionally, the function can call event.Get("stopTimer").Invoke() to stop the timer
func (e *Element) OnHoldKey(wait time.Duration, repeat time.Duration, f func(*Element, jsext.Event)) *Element {
	var debouncer *debouncer
	var setStopTimer = func(event jsext.Event) {
		event.Set("stopTimer", js.FuncOf(func(this js.Value, args []js.Value) any {
			debouncer.Stop()
			return nil
		}))
	}

	e.OnKeyDown(func(this *Element, event jsext.Event) {
		setStopTimer(event)
		debouncer.Stop()
		f(this, event)
		debouncer = afterFunc(wait, repeat, event, func() {
			for range debouncer.ticker.C {
				f(this, event)
			}
		})
	})
	e.OnKeyUp(func(this *Element, event jsext.Event) {
		debouncer.Stop()
	})
	e.OnBlur(func(this *Element, event jsext.Event) {
		debouncer.Stop()
	})
	return e
}

// OnHoldClick adds an event listener to the Element
//
// wait is the time to wait after the first click before starting to repeat
// repeat is the time between each function call
//
// Optionally, the function can call event.Get("stopTimer").Invoke() to stop the timer
func (e *Element) OnHoldClick(wait time.Duration, repeat time.Duration, f func(this *Element, event jsext.Event)) *Element {
	var debouncer *debouncer

	var setStopTimer = func(event jsext.Event) {
		event.Set("stopTimer", js.FuncOf(func(this js.Value, args []js.Value) any {
			debouncer.Stop()
			return nil
		}))
	}

	e.OnMouseDown(func(this *Element, event jsext.Event) {
		setStopTimer(event)
		debouncer.Stop()
		f(this, event)
		debouncer = afterFunc(wait, repeat, event, func() {
			for range debouncer.ticker.C {
				f(this, event)
			}
		})
	})
	e.OnMouseUp(func(_ *Element, _ jsext.Event) {
		debouncer.Stop()
	})
	e.OnMouseLeave(func(_ *Element, _ jsext.Event) {
		debouncer.Stop()
	})
	return e
}

// OnScrolledToBottom adds an event listener to the Element
//
// callEach is the minimum time between each function call
//
// pxFromBottom is the number of pixels from the bottom of the element to call the function
func (e *Element) OnScrolledToBottom(callEach time.Duration, pxFromBottom int, f func(*Element, jsext.Event)) js.Func {
	var lastCall = time.Time{}
	return e.AddEventListener("scroll", func(this *Element, event jsext.Event) {
		if time.Since(lastCall) < callEach {
			return
		}
		var target = event.Target()
		var offsetHeight = target.Get("offsetHeight").Int()
		var scrollTop = target.Get("scrollTop").Int()
		var scroll = offsetHeight + scrollTop
		var height = target.Get("scrollHeight").Int() - pxFromBottom
		if scroll >= height {
			lastCall = time.Now()
			f(this, event)
		}
	})
}
