//go:build js && wasm && !tinygo
// +build js,wasm,!tinygo

package messages

import (
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/Nigel2392/jsext/helpers"
	"github.com/Nigel2392/jsext/helpers/csshelpers"
)

// For some reason this package causes a nil pointer dereference in tinygo.
// Thus, we exclude it from the build.
// This unfortunately does mean that we cannot use the message framework at all.

// Time to wait for garbage collection of messages.
const tickerWaitMS = 100

// Add to messages to this struct to display them.
// Can be implemented in a custom way.
var ActiveMessages = RunNewMessages(tickerWaitMS)

// Message struct to be used for all messages.
type Messages struct {
	Active          []*Message
	StopMessageLoop chan struct{}
	ticker          *time.Ticker
	mu              *sync.Mutex
}

// Initialize a new message queue.
func NewMessages() *Messages {
	var m = &Messages{StopMessageLoop: make(chan struct{}), mu: &sync.Mutex{}, Active: make([]*Message, 0)}
	return m
}

// Add a new message to the queue, run it.
func RunNewMessages(tickerTimeMS int) *Messages {
	m := NewMessages()
	m.ticker = time.NewTicker(time.Duration(tickerTimeMS) * time.Millisecond)
	go m.collect()
	return m
}

// Stop a message queue.
func (m *Messages) Stop() {
	m.StopMessageLoop <- struct{}{}
}

// Create a new message.
func (m *Messages) New(t string, c string, expireMS int, widthPX int) *Message {
	var id string = "gohtml-message-" + t + "-" + helpers.RandStringBytesMaskImprSrcUnsafe(10)

	var messageColor string

	switch t {
	case "error":
		messageColor = "#aa0000"
	case "warning":
		messageColor = "#aa7700"
	case "info":
		messageColor = "#0000aa"
	case "success":
		messageColor = "#00aa00"
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

// Delete message from the queue.
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

// Delete message.
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

// Create a new error message.
func (m *Messages) NewError(c string, expire int, widthPX int) *Message {
	return m.New("error", c, expire, widthPX)
}

// Create a new warning message.
func (m *Messages) NewWarning(c string, expire int, widthPX int) *Message {
	return m.New("warning", c, expire, widthPX)
}

// Create a new info message.
func (m *Messages) NewInfo(c string, expire int, widthPX int) *Message {
	return m.New("info", c, expire, widthPX)
}

// Create a new success message.
func (m *Messages) NewSuccess(c string, expire int, widthPX int) *Message {
	return m.New("success", c, expire, widthPX)
}

// Message struct.
type Message struct {
	Type    string
	Content string
	id      string
	jsVal   js.Value
	Expire  time.Time
}

// Expired returns true if the message has expired.
func (m *Message) Expired() bool {
	return m.Expire.Before(time.Now()) && !m.Expire.IsZero()
}

// Remove the message from the DOM.
func (m *Message) Remove() {
	var elem = js.Global().Get("document").Call("getElementById", m.id)
	if elem.IsUndefined() {
		return
	}
	elem.Call("remove")
}
