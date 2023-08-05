package shortcuts

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Nigel2392/jsext/v2/jse"
)

type TimeCounter struct {
	time       time.Time
	element    *jse.Element
	ticker     *time.Ticker
	formatFunc func(*TimeTracker) string
	updateFunc func(*TimeTracker, *jse.Element) error
}

func NewTimeCounter(elem *jse.Element) *TimeCounter {
	return &TimeCounter{
		time:    time.Now(),
		element: elem,
		formatFunc: func(t *TimeTracker) string {
			return t.Format("%YR years, %MO months, %DD days, %HH hours, %MM minutes, %SS seconds")
		},
	}
}

func (c *TimeCounter) Increment() { c.Add(time.Second) }

func (c *TimeCounter) FormatFunc(f func(*TimeTracker) string)              { c.formatFunc = f }
func (c *TimeCounter) UpdateFunc(f func(*TimeTracker, *jse.Element) error) { c.updateFunc = f }

func (c *TimeCounter) Display(Time time.Time) error {
	if c.element == nil || c.element.IsNull() || c.element.IsUndefined() {
		return errors.New("element is nil")
	}
	if c.formatFunc == nil {
		return errors.New("formatFunc is nil")
	}
	if c.updateFunc != nil {
		return c.updateFunc(c.Tracker(), c.element)
	}
	// Display time until
	timeTracker := NewTimeTracker(Time)
	c.element.InnerHTML(c.formatFunc(timeTracker))
	return nil
}

func (c *TimeCounter) Reset() error {
	c.time = time.Now()
	return c.Display(c.time)
}

func (c *TimeCounter) Set(t time.Time) error {
	c.time = t
	return c.Display(c.time)
}

func (c *TimeCounter) Get() time.Time { return c.time }

func (c *TimeCounter) Add(t time.Duration) error {
	c.time = c.time.Add(t)
	return c.Display(c.time)
}

func (c *TimeCounter) Sub(t time.Duration) error {
	c.time = c.time.Add(t)
	return c.Display(c.time)
}

func durationFromAny(t time.Time, arg any) (time.Duration, bool, error) {
	switch v := arg.(type) {
	case time.Duration:
		return v, false, nil
	case func(time.Time) time.Duration:
		return v(t), true, nil
	}
	return 0, false, errors.New("data type mismatch")
}

// each must either be:
// - a time.Duration
// - a function which takes a time.Time and returns a time.Duration
func (c *TimeCounter) Live(each any) {
	c.Display(c.time)
	var tickerDuration, isFunc, err = durationFromAny(c.time, each)
	if err != nil {
		panic(err)
	}
	c.ticker = time.NewTicker(tickerDuration)
	go func() {
		for range c.ticker.C {
			err = c.Display(c.time)
			if err != nil {
				c.ticker.Stop()
				return
			}
			if isFunc {
				tickerDuration, _, err = durationFromAny(c.time, each)
				if err != nil {
					c.ticker.Stop()
					return
				}
				c.ticker.Stop()
				c.Live(tickerDuration)
			}
		}
	}()
}

func (c *TimeCounter) StopLive() { c.ticker.Stop() }

func (c *TimeCounter) Tracker() *TimeTracker { return NewTimeTracker(c.time) }

func (c *TimeCounter) Date(year int, month time.Month, day int, hour int, min int, sec int, nsec int) {
	c.Set(time.Now().Add(time.Since(time.Date(year, month, day, hour, min, sec, nsec, time.UTC)) * -1))
}

type TimeTracker struct {
	Years   int
	Months  int
	Days    int
	Hours   int
	Minutes int
	Seconds int
	time    time.Time
}

func NewTimeTracker(Time time.Time) *TimeTracker {
	if Time.IsZero() {
		return &TimeTracker{}
	}
	var t time.Duration
	if Time.After(time.Now()) {
		t = time.Until(Time)
	} else {
		t = time.Since(Time)
	}
	var total = int(t.Seconds())
	var years = int(total / (60 * 60 * 24 * 365))
	var months = int(total / (60 * 60 * 24 * 30) % 12)
	var days = int(total / (60 * 60 * 24) % 30)
	var hours = int(total / (60 * 60) % 24)
	var minutes = int(total/60) % 60
	var seconds = int(total % 60)

	return &TimeTracker{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
		time:    Time,
	}
}

func (t *TimeTracker) Strings() (string, string, string, string, string, string) {
	var string = t.Format("%YR-%MO-%DD-%HH-%MM-%SS")
	var split = strings.Split(string, "-")
	return split[0], split[1], split[2], split[3], split[4], split[5]
}

func (t *TimeTracker) Format(format string) string {
	format = strings.ReplaceAll(format, "%YR", strconv.Itoa(t.Years))
	format = strings.ReplaceAll(format, "%MO", strconv.Itoa(t.Months))
	format = strings.ReplaceAll(format, "%DD", strconv.Itoa(t.Days))
	format = strings.ReplaceAll(format, "%HH", strconv.Itoa(t.Hours))
	format = strings.ReplaceAll(format, "%MM", strconv.Itoa(t.Minutes))
	format = strings.ReplaceAll(format, "%SS", strconv.Itoa(t.Seconds))
	return format
}

func (t *TimeTracker) IsZero() bool {
	return t.time.IsZero()
}

func (t *TimeTracker) IsPast() bool {
	return t.time.Before(time.Now())
}

func (t *TimeTracker) IsFuture() bool {
	return t.time.After(time.Now())
}

func (t *TimeTracker) Time() time.Time {
	return t.time
}
