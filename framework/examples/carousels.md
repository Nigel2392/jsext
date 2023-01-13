# Example image carousel component

A simple example of how to create an image carousel component.
    
```go
    var carosuel = carousels.Image([]string{
		"/static/logo.png",
		"/static/1.png",
		"/static/2.png",
		"/static/3.png",
	}, &carousels.Options{
		Width:         "600px",
		Height:        "400px",
		Background:    "#ffffff",
		ControlsColor: "rgba(0, 0, 0, 0.5)",
		Controls:      true,
		Indicators:    true,
		Border:        "1px solid #333333",
		ActiveItem:    2,
	}, true)
    Application.Render(carousel)
```