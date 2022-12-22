//go:build js && wasm
// +build js,wasm

package messages

import (
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/Nigel2392/jsext/helpers"
	"github.com/Nigel2392/jsext/helpers/csshelpers"
)

const tickerWaitMS = 100

var ActiveMessages = RunNewMessages(tickerWaitMS)

type Messages struct {
	Active          []*Message
	StopMessageLoop chan struct{}
	ticker          *time.Ticker
	mu              *sync.Mutex
}

func NewMessages() *Messages {
	var m = &Messages{StopMessageLoop: make(chan struct{}), mu: &sync.Mutex{}, Active: make([]*Message, 0)}
	return m
}

func RunNewMessages(tickerTimeMS int) *Messages {
	m := NewMessages()
	m.ticker = time.NewTicker(time.Duration(tickerTimeMS) * time.Millisecond)
	go m.collect()
	return m
}

func (m *Messages) Stop() {
	m.StopMessageLoop <- struct{}{}
}

func (m *Messages) New(t string, c string, expireMS int, widthPX int) *Message {
	var id string = "gohtml-message-" + t + "-" + helpers.RandStringBytesMaskImprSrcUnsafe(10)

	var messageColor string

	switch t {
	case "error":
		messageColor = "#ff0000"
	case "warning":
		messageColor = "#ff9900"
	case "info":
		messageColor = "#0000ff"
	case "success":
		messageColor = "#00ff00"
	default:
		messageColor = "#000000"
	}

	var msg = &Message{Type: t, Content: c, id: id}
	// Create the message in the DOM
	msg.jsVal = js.Global().Get("document").Call("createElement", "div")
	// Set the id
	msg.jsVal.Call("setAttribute", "id", id)
	// Set the content
	msg.jsVal.Set("innerHTML", c)
	// Set the classes
	msg.jsVal.Get("classList").Call("add", "gohtml-message", t)

	var toTop int = 0
	if len(m.Active) > 0 {
		for _, v := range m.Active {
			var oh = v.jsVal.Get("offsetHeight")
			if oh.IsUndefined() || oh.IsNull() {
				continue
			}
			toTop += oh.Int() + widthPX/15
		}
	}
	toTop += widthPX / 15

	// Create the style
	// var height = int(float64(widthPX) / 2)
	var style = []string{
		"position: fixed",
		// "bottom: calc(5% +" + ToPX((height+20)*len(ActiveMessages.Active)) + ")",
		"bottom: " + csshelpers.ToPX(toTop),
		"right: 2%",
		"width: " + csshelpers.ToPX(widthPX),
		"font-size: " + csshelpers.ToPX(int(float64(widthPX)/9)),
		"padding: " + csshelpers.ToPX(int(float64(widthPX)/15)),
		"text-align: center",
		"border-radius: " + csshelpers.ToPX(int(float64(widthPX)/15)),
		"box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.5)",
		"background-color: " + messageColor,
		"color: #ffffff",
		"cursor: pointer",
		"z-index: 1000",
	}
	// Set the style
	msg.jsVal.Call("setAttribute", "style", strings.Join(style, ";"))
	// Delete the message when clicked
	// msg.jsVal.Call("setAttribute", "onclick", "this.remove()")
	msg.jsVal.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m.Delete(msg)
		return nil
	}))
	// Set the expire time
	msg.Expire = time.Now().Add(time.Millisecond * time.Duration(expireMS))
	// Add message to ActiveMessages
	m.Active = append(m.Active, msg)
	// Append the message to the DOM
	js.Global().Get("document").Get("body").Call("appendChild", msg.jsVal)
	return msg
}

func (m *Messages) collect() {
	for {
		select {
		case <-m.StopMessageLoop:
			return
		case <-m.ticker.C:
			for _, msg := range m.Active {
				if msg.Expired() {
					m.Delete(msg)
				}
			}
		}
	}
}

func (m *Messages) Delete(msg *Message) {
	msg.Remove()
	for i, v := range m.Active {
		if v == msg {
			m.mu.Lock()
			if i == len(m.Active)-1 {
				m.Active = m.Active[:i]
			} else if i == 0 {
				m.Active = m.Active[1:]
			} else {
				m.Active = append(m.Active[:i], m.Active[i+1:]...)
			}
			m.mu.Unlock()
		}
	}
}

func (m *Messages) NewError(c string, expire int, widthPX int) *Message {
	return m.New("error", c, expire, widthPX)
}

func (m *Messages) NewWarning(c string, expire int, widthPX int) *Message {
	return m.New("warning", c, expire, widthPX)
}

func (m *Messages) NewInfo(c string, expire int, widthPX int) *Message {
	return m.New("info", c, expire, widthPX)
}

func (m *Messages) NewSuccess(c string, expire int, widthPX int) *Message {
	return m.New("success", c, expire, widthPX)
}

type Message struct {
	Type    string
	Content string
	id      string
	jsVal   js.Value
	Expire  time.Time
}

func (m *Message) Expired() bool {
	return m.Expire.Before(time.Now()) && !m.Expire.IsZero()
}

func (m *Message) Remove() {
	var elem = js.Global().Get("document").Call("getElementById", m.id)
	if elem.IsUndefined() {
		return
	}
	elem.Call("remove")
}
