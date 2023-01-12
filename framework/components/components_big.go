//go:build js && wasm && big
// +build js,wasm,big

package components

import (
	"github.com/Nigel2392/jsext/canvas"
	"github.com/Nigel2392/jsext/framework/components/carousels"
	"github.com/Nigel2392/jsext/framework/components/loaders"
	"github.com/Nigel2392/jsext/framework/components/misc"
	"github.com/Nigel2392/jsext/framework/components/navbars"
	"github.com/Nigel2392/jsext/framework/elements"
	"github.com/Nigel2392/jsext/framework/graphs/charts"
	"github.com/Nigel2392/jsext/framework/graphs/options"
)

type carousels_module struct {
	Plain func(*carousels.Options) *elements.Element
	Image func([]string, *carousels.Options) *elements.Element
}

type loaders_module struct {
	LoaderRotatingBlock  func(appendTo string, className string, deleteOnFinish bool) Loader
	LoaderHexagonRolling func(appendTo string, className string, deleteOnFinish bool) Loader
	LoaderMultiRing      func(appendTo string, className string, deleteOnFinish bool) Loader
	LoaderRing           func(appendTo string, className string, deleteOnFinish bool) Loader
}

type charts_module struct {
	Bar      func(canvas.Canvas, options.GraphOptions)
	Line     func(canvas.Canvas, options.GraphOptions)
	Pie      func(canvas.Canvas, options.GraphOptions)
	Doughnut func(canvas.Canvas, options.GraphOptions)
}

type navbars_module struct {
	Official func(logo *navbars.Logo, urls *elements.URLs) *elements.Element
	Search   func(logo *navbars.Logo, urls *elements.URLs) (*elements.Element, []*elements.Element)
	Custom   func(logo *navbars.Logo, urls *elements.URLs, bg, fg string, middle ...*elements.Element) *elements.Element
}

var Carousels = carousels_module{
	Plain: carousels.Plain,
	Image: carousels.Image,
}

var Loaders = loaders_module{
	LoaderRotatingBlock: func(appendTo string, className string, deleteOnFinish bool) Loader {
		return loaders.NewLoader(appendTo, className, deleteOnFinish, loaders.LoaderRotatingBlock)
	},
	LoaderHexagonRolling: func(appendTo string, className string, deleteOnFinish bool) Loader {
		return loaders.NewLoader(appendTo, className, deleteOnFinish, loaders.LoaderHexagonRolling)
	},

	LoaderMultiRing: func(appendTo string, className string, deleteOnFinish bool) Loader {
		return loaders.NewLoader(appendTo, className, deleteOnFinish, loaders.LoaderMultiRing)
	},

	LoaderRing: func(appendTo string, className string, deleteOnFinish bool) Loader {
		return loaders.NewLoader(appendTo, className, deleteOnFinish, loaders.LoaderRing)
	},
}

var Charts = charts_module{
	Bar:  charts.Bar,
	Line: charts.Line,
	Pie: func(c canvas.Canvas, opts options.GraphOptions) {
		charts.Pie(c, opts, false)
	},
	Doughnut: func(c canvas.Canvas, opts options.GraphOptions) {
		charts.Pie(c, opts, true)
	},
}

var Navbars = navbars_module{
	Official: navbars.Official,
	Search:   navbars.Search,
	Custom:   navbars.Custom,
}

var SearchBar = misc.SearchBar
