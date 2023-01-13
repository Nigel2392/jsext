package convert

import (
	"strconv"
	"strings"
	"time"
)

func FormatNumber(value any) string {
	switch value := value.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.Itoa(int(value))
	case uint:
		return strconv.Itoa(int(value))
	case uint8:
		return strconv.Itoa(int(value))
	case uint16:
		return strconv.Itoa(int(value))
	case uint32:
		return strconv.Itoa(int(value))
	case uint64:
		return strconv.Itoa(int(value))
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		panic("Value must be a number type")
	}
}

func ToFloat(value any) float64 {
	switch value := value.(type) {
	case int:
		return float64(value)
	case int8:
		return float64(value)
	case int16:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	case uint:
		return float64(value)
	case uint8:
		return float64(value)
	case uint16:
		return float64(value)
	case uint32:
		return float64(value)
	case uint64:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	default:
		panic("Value must be a number type")
	}
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
