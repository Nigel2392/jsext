# Simple scrollable SPA

Here you can see how to create a simple scrollable single page application.

```go
var Application = scroll.App("#app", &scroll.Options{
	ScrollAxis:    scroll.ScrollAxisY,
	ScrollThrough: true,
	// GradientTo:    "to bottom right",
	ClassPrefix: "scrollable-app",
})

// Main webassembly entry point
func main() {
	var Mainmenu = menus.NewMenuOptions(menus.Left)
	Mainmenu.CSSMap[menus.OverlayBackgroundColor] = "rgba(0,0,0,0.8)"
	Mainmenu.URLs.FromElements(true,
		elements.A("#home", "Home"),
		elements.A("#about", "About"),
		elements.A("#contact", "Contact"),
	)
	Application.SetNavbar(menus.Blurry(Mainmenu))
	var colorbgs = Application.Backgrounds(scroll.BackgroundTypeColor, "#9200ff", "#ff0000", "#0000ff")
	var imgbgs = Application.Backgrounds(scroll.BackgroundTypeImage, "/static/1.png")

	colorbgs.AddGradient("rgba(1,1,1,0.2)", "rgba(255, 255, 255, 0.4)")
	imgbgs.AddGradient("rgba(0,0,0,0.2)", "rgba(152, 92, 255, 0.2)")

	var (
		Home_Component = elements.Div().AttrStyle("color:white;").Append(
			elements.H1("Home"),
			elements.P("Welcome to the Home page."),
		)
		About_Component = elements.Div().AttrStyle("color:white;").Append(
			elements.H1("About"),
			elements.P("Welcome to the About page."),
		)
		Contact_Component = elements.Div().AttrStyle("color:white;").Append(
			elements.H1("Contact"),
			elements.P("Welcome to the Contact page."),
		)
	)

	Application.AddPage("home", Home_Component)
	Application.AddPage("about", About_Component)
	Application.AddPage("contact", Contact_Component)

	Application.Run()
}
```