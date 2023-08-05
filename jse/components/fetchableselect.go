package shortcuts

import (
	"time"

	"github.com/Nigel2392/jsext/v2"
	"github.com/Nigel2392/jsext/v2/jse"
)

type SearchableSelectOption interface {
	Display() string
	Value() string
}

func NewSearchableSelectOption(display string, value string) SearchableSelectOption {
	return searchableSelectBoxOption{
		display: display,
		value:   value,
	}
}

type searchableSelectBoxOption struct {
	display string
	value   string
}

func (f searchableSelectBoxOption) Display() string {
	return f.display
}

func (f searchableSelectBoxOption) Value() string {
	return f.value
}

type SearchableSelectBoxOptions struct {
	NullOption     SearchableSelectOption
	InitialOption  SearchableSelectOption
	InitialOptions []SearchableSelectOption
	LoadFunc       func(string) (results []SearchableSelectOption, cancel bool)
	OnSelect       func(SearchableSelectOption)
	InputOptions   *jse.InputOptions
	Delay          time.Duration
	MinChars       int
}

func (opts *SearchableSelectBoxOptions) Defaults() {
	if opts.Delay == 0 {
		opts.Delay = 200 * time.Millisecond
	}
	if opts.MinChars == 0 {
		opts.MinChars = 3
	}
}

func (s *SelectBox) optionsFetchFunc(selectBox, input *jse.Element) func(_ *jse.Element, _ jsext.Event) {
	return func(this *jse.Element, e jsext.Event) {
		var value = this.Get("value").String()
		go func() {
			var results, cancel = s.opts.LoadFunc(value)
			if cancel {
				return
			}
			selectBox.InnerHTML("")
			if s.opts.NullOption != nil {
				s.appendOption(selectBox, this, s.opts.NullOption)
			}
			for _, result := range results {
				s.appendOption(selectBox, this, result)
			}
		}()
	}
}

func (s *SelectBox) defaultOptionsEventFunc(selectBox *jse.Element) func(_ *jse.Element, _ jsext.Event) {
	if s.opts.InitialOptions == nil {
		return nil
	}
	return func(_ *jse.Element, _ jsext.Event) {
		selectBox.InnerHTML("")
		if s.opts.NullOption != nil {
			s.appendOption(selectBox, nil, s.opts.NullOption)
		}
		for _, option := range s.opts.InitialOptions {
			s.appendOption(selectBox, nil, option)
		}
	}
}

type SelectBox struct {
	root     *jse.Element
	input    *jse.Element
	selected SearchableSelectOption
	onSelect func(*SelectBox, SearchableSelectOption)
	opts     *SearchableSelectBoxOptions
}

func (s *SelectBox) Element() *jse.Element {
	return s.root
}

func (s *SelectBox) Selected() SearchableSelectOption {
	return s.selected
}

func SearchableSelectBox(opts SearchableSelectBoxOptions) *SelectBox {
	opts.Defaults()

	var selectBoxContainer = jse.NewElement("select-box-container")
	var selectBoxInputHeader = selectBoxContainer.NewElement("select-box-header")
	var selectBoxInput = selectBoxInputHeader.Input("text", "select-box-input", opts.InputOptions)
	var selectBox = selectBoxContainer.NewElement("select-box")
	var s = &SelectBox{
		root:  selectBoxContainer,
		input: selectBoxInput,
		onSelect: func(s *SelectBox, option SearchableSelectOption) {
			s.selected = option
			if opts.OnSelect != nil {
				s.opts.OnSelect(option)
			}
		},
		opts: &opts,
	}
	if opts.NullOption != nil {
		s.appendOption(selectBox, selectBoxInput, opts.NullOption)
	}
	if opts.InitialOption != nil {
		s.selected = opts.InitialOption
		selectBoxInput.Set("value", opts.InitialOption.Display())
	}
	for _, option := range opts.InitialOptions {
		s.appendOption(selectBox, selectBoxInput, option)
	}
	AfterCharsAndDuration(
		opts.MinChars,
		opts.Delay,
		selectBoxInput,
		s.optionsFetchFunc(selectBox, selectBoxInput),
		s.defaultOptionsEventFunc(selectBox),
	)
	return s
}

func AfterCharsAndDuration(amount int, delay time.Duration, e *jse.Element, moreOrEqual, lessThan func(this *jse.Element, e jsext.Event)) {
	var timeout *time.Timer
	e.AddEventListener("input", func(this *jse.Element, e jsext.Event) {
		if timeout != nil {
			timeout.Stop()
		}
		timeout = time.AfterFunc(delay, func() {
			var thisValue = this.Get("value").String()
			if len(thisValue) >= amount && moreOrEqual != nil {
				moreOrEqual(this, e)
			} else if lessThan != nil {
				lessThan(this, e)
			}
		})
	})
}

func (s *SelectBox) appendOption(selectBox, input *jse.Element, option SearchableSelectOption) {
	var resultElement = selectBox.NewElement("select-box-result")
	resultElement.InnerText(option.Display())
	resultElement.Dataset().Set("value", option.Value())
	resultElement.AddEventListener("click", func(this *jse.Element, e jsext.Event) {
		if input != nil {
			input.Set("value", option.Display())
		}
		go s.onSelect(s, option)
	})
}
