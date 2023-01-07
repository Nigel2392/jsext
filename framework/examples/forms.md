```go
	var div = elements.Div()
    // Struct to parse the form into. If values are set in the struct, they will be used as the default values in the form.
	type HelloW struct{ Enter_Your_Name string }

	var hw = HelloW{}

    // Create a form from the struct
	var form = forms.StructToForm(&hw, "jsext-example-label", "jsext-example-input", "/", "")

    // Listen for the form to be submitted and parse the values into the struct.
	form.OnSubmitToStruct(&hw, func(strct any, elements []jsext.Element) {
		jsext.QuerySelector(".jsext-example-label").InnerHTML("Hello " + hw.Enter_Your_Name + "!")
	})

	form.Element().Append(
        // Add a submit button to the form.
		elements.Button().InnerHTML("Submit").AttrType("submit"),
	)

    // Render the form.
	div.Append(form.Element())
```