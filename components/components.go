package components

import (
	"create_elems/jsext"
)

type Component interface {
	Render() jsext.Element
}

type Loader interface {
	Stop()        // Stop the loader.
	Show()        // Show the loader.
	Run(f func()) // Run the function, finalize loader automatically.
	Finalize()    // Finalize loader.
}

type URL struct {
	Name string
	Url  string
}
