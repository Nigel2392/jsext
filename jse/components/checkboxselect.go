package shortcuts

import "github.com/Nigel2392/jsext/v2/jse"

type CheckboxSelect struct {
	display string
	*jse.Element
	container *jse.Element
	Boxes     []*jse.Element
}

func NewCheckboxSelect(display string) *CheckboxSelect {
	return &CheckboxSelect{
		display: display,
	}
}

func (c *CheckboxSelect) checkNew() {
	if c.Element == nil {
		c.Element = jse.NewElement("checkbox-select")
		c.Element.NewElement("checkbox-select-label", c.display)
		c.Boxes = make([]*jse.Element, 0)
		c.container = c.Element.NewElement("checkbox-select-container")
	}
}

func (c *CheckboxSelect) SetDisplay(display string) {
	c.checkNew()
	c.display = display
	var e = c.Element.FirstElementChild()
	e.InnerText(display)
}

func (c *CheckboxSelect) Option(name string, value string, selected bool) *jse.Element {
	c.checkNew()
	var box = c.container.NewElement("checkbox-select-box")
	var input = box.Input("checkbox", name, nil)
	if value != "" {
		input.Set("value", value)
	}
	if selected {
		input.Set("checked", true)
	}
	box.Label(name, name)
	c.Boxes = append(c.Boxes, box)
	return box
}

func (c *CheckboxSelect) Selected() []string {
	c.checkNew()
	var selected = make([]string, 0, len(c.Boxes))
	for _, box := range c.Boxes {
		var firstChild = box.FirstElementChild()
		if !firstChild.IsNull() && !firstChild.IsUndefined() {
			var checked = firstChild.Get("checked")
			if checked.IsNull() || checked.IsUndefined() {
				continue
			}
			if checked.Bool() {
				var value = firstChild.Get("value")
				if value.IsNull() || value.IsUndefined() {
					value = firstChild.Get("name")
				}
				selected = append(selected, value.String())
			}
		}
	}
	return selected
}

func (c *CheckboxSelect) SetSelected(selected []string) {
	c.checkNew()
	for _, box := range c.Boxes {
		var firstChild = box.FirstElementChild()
		if !firstChild.IsNull() && !firstChild.IsUndefined() {
			var value = firstChild.Get("value")
			if value.IsNull() || value.IsUndefined() {
				value = firstChild.Get("name")
			}
			for _, s := range selected {
				if s == value.String() {
					firstChild.Set("checked", true)
					break
				}
			}
		}
	}
}

func (c *CheckboxSelect) Clear() {
	c.checkNew()
	for _, box := range c.Boxes {
		var firstChild = box.FirstElementChild()
		if !firstChild.IsNull() && !firstChild.IsUndefined() {
			firstChild.Set("checked", false)
		}
	}
}

func (c *CheckboxSelect) RootElement() *jse.Element {
	c.checkNew()
	return c.Element
}

func (c *CheckboxSelect) LabelElement() *jse.Element {
	c.checkNew()
	var e = c.Element.FirstElementChild()
	return &e
}

func (c *CheckboxSelect) BoxElement() *jse.Element {
	c.checkNew()
	return c.container
}
