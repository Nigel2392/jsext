# A Countdown component example

A simple example of how to use the default countdown component.

```go
// Create a counter element for use in the component.
var counterElement = elements.H1()
// Create a new countdown component.
var counter = misc.NewTimeCounter(counterElement)
// Choose what happens when the counter updates.
counter.UpdateFunc(func(tt *convert.TimeTracker, *elements.Element) {
    // Regenerate component with current timer.
	e.Children = nil
	var yearsStr, monthsStr, daysStr, hoursStr, minutesStr, secondsStr string = tt.Strings()
	e.ClearInnerHTML()
    // Years
	if tt.Years > 0 {
		e.Append(elements.Span(yearsStr).AttrStyle("color: #ff0000"))
		if tt.Years > 1 {
			e.Append(elements.Span(" Years "))
		} else {
			e.Append(elements.Span(" Year "))
		}
	}
    // Months
	if tt.Months > 0 {
		e.Append(elements.Span(monthsStr).AttrStyle("color: #ff0000"))
		if tt.Months > 1 {
			e.Append(elements.Span(" Months "))
		} else {
			e.Append(elements.Span(" Month "))
		}
	}
    // Days
	if tt.Days > 0 {
		e.Append(elements.Span(daysStr).AttrStyle("color: #ff0000"))
		if tt.Days > 1 {
			e.Append(elements.Span(" Days "))
		} else {
			e.Append(elements.Span(" Day "))
		}
	}
    // Hours
	if tt.Hours > 0 || tt.Days >= 0 || tt.Months >= 0 || tt.Years >= 0 {
		e.Append(elements.Span(hoursStr).AttrStyle("color: #ff0000"))
		e.Append(elements.Span(":"))
	}
    // Minutes
	if tt.Minutes > 0 || tt.Hours >= 0 || tt.Days >= 0 || tt.Months >= 0 || tt.Years >= 0 {
		e.Append(elements.Span(minutesStr).AttrStyle("color: #ff0000"))
		e.Append(elements.Span(":"))
	}
    // Seconds
	e.Append(elements.Span(secondsStr).AttrStyle("color: #ff0000"))
	if tt.Years == 0 && tt.Months == 0 && tt.Days == 0 && tt.Hours == 0 && tt.Minutes == 0 {
		if tt.Seconds > 1 {
			e.Append(elements.Span(" Seconds"))
		} else {
			e.Append(elements.Span(" Second"))
		}
	}
})
// Set the counter date.
counter.Date(2023, 12, 13, 18, 5, 0, 0)
// Update the counter element every one second.
counter.Live()
// Render the counter element.
Application.Render(counterElement)
```