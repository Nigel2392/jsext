# Simple scrollable SPA

Here you can see how to create a simple scrollable single page application.

```go
func main() {
	var Application = scroll.App("#app", &scroll.Options{
		ScrollAxis:    scroll.ScrollAxisX,
		ScrollThrough: true,
		ClassPrefix:   "scrollable-app",
	})
	var Mainmenu = menus.NewMenuOptions(menus.Left)
	Mainmenu.CSSMap[menus.OverlayBackgroundColor] = "rgba(0,0,0,0.8)"
	Mainmenu.URLs.FromElements(true,
		elements.A("#home", "Home"),
		elements.A("#about", "About"),
		elements.A("#contact", "Contact"),
	)
	Application.SetNavbar(menus.Blurry(Mainmenu))
	var colorbgs = Application.Backgrounds(scroll.BackgroundTypeColor, "#ff0000", "#9200ff", "#0000ff")
	var imgbgs = Application.Backgrounds(scroll.BackgroundTypeImage, "/static/1.png")

	colorbgs[0].AddGradient(scroll.GradientTypeRadial, "circle", "rgba(0,0,0,0.4)", "rgba(255, 255, 255, 0.2)")
	colorbgs[1].AddGradient(scroll.GradientTypeLinear, "to right", "rgba(0,0,0,0.4)", "rgba(255, 255, 255, 0.2)")
	imgbgs.AddGradient(scroll.GradientTypeRadial, "circle", "rgba(0,0,0,0.4)", "rgba(152, 92, 255, 0.2)")

	var pageNames = []string{"home", "about", "contact"}
	for _, name := range pageNames {
		var title = elements.H1(name)
		var text = elements.P("Welcome to the " + name + " page.")
		title.Animations.FadeIn(500, elements.UseIntersectionObserver, true)
		text.Animations.FadeIn(500, elements.UseIntersectionObserver, true)
		var page = elements.Div().AttrStyle("color:white;").Append(
			title,
			text,
		)
		Application.AddPage(name, page)
	}

	Application.Run()
}
```