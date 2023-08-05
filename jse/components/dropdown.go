package shortcuts

import (
	"github.com/Nigel2392/jsext/v2/jse"
)

//	<div class="jse-dropdown-container">
//		<button class="jse-dropdown-button">
//			Open
//		</button>
//		<div class="jse-dropdown-dropbox">
//			<div style="display:flex;flex-direction:row;">
//				<div> <button>ITEM 1</button> </div>
//				<div> <button>ITEM 2</button> </div>
//				<div> <button>ITEM 3</button> </div>
//				<div> <button>ITEM 4</button> </div>
//			</div>
//			<div> ITEM 1 </div>
//			<div> ITEM 2 </div>
//			<div> ITEM 3 </div>
//			<div> ITEM 4 </div>
//		</div>
//	</div>

type DropdownOptions struct {
	Button      *jse.Element
	Menu        []*jse.Element
	ClassPrefix string
}

func (d *DropdownOptions) Defaults() {
	if d.ClassPrefix == "" {
		d.ClassPrefix = "jse-"
	}
	if d.Button == nil || d.Button != nil && d.Button.IsUndefined() {
		d.Button = jse.NewElement("button")
	}
	if d.Menu == nil {
		d.Menu = []*jse.Element{}
	}
}

func Dropdown(options DropdownOptions) *jse.Element {

	options.Defaults()

	var d = jse.NewElement("div")
	d.ClassList().Add(options.ClassPrefix + "dropdown-container")

	var button = options.Button
	button.ClassList().Add(options.ClassPrefix + "dropdown-button")

	var dropbox = jse.NewElement("div")
	dropbox.ClassList().Add(options.ClassPrefix + "dropdown-dropbox")

	for _, menu := range options.Menu {
		dropbox.AppendChild(menu)
	}

	d.AppendChild(button)
	d.AppendChild(dropbox)

	return d
}
