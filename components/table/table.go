package table

import (
	"reflect"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/components"
	"github.com/Nigel2392/jsext/elements"
)

type NoStructAvailable int

type Table[T any] struct {
	model          []T
	additionalCols map[string]func(model T) *elements.Element
	root           *elements.Element
}

func New() *Table[NoStructAvailable] {
	var t = Table[NoStructAvailable]{
		root: elements.Div().AttrClass("jsext-table-root"),
	}
	return &t
}

func (t *Table[T]) hasStruct() bool {
	return reflect.TypeOf(t.model).Kind() != reflect.Slice || reflect.TypeOf(t.model).Elem().Kind() != reflect.TypeOf(NoStructAvailable(0)).Kind()
}

func NewFromStruct[T any](s []T, additionalCols map[string]func(model T) *elements.Element) *Table[T] {
	var t = Table[T]{
		model:          s,
		additionalCols: additionalCols,
		root:           elements.Div().AttrClass("jsext-table-root"),
	}
	return &t
}

func (t *Table[T]) Render() jsext.Element {
	if t.hasStruct() {
		return t.create().Render()
	}
	return t.root.Render()
}

func (t *Table[T]) create() *elements.Element {
	var reflModel = components.StructKind(t.model)
	var rowNames = GetStructFieldNames(reflModel)
	var table = t.root.Table().AttrClass("jsext-table")
	var thead = table.Thead()
	var tbody = thead.Tbody()
	var tr = thead.Tr()
	for _, rowName := range rowNames {
		tr.Th().Span(rowName)
	}
	for _, model := range t.model {
		var tr = tbody.Tr()
		for _, colTag := range rowNames {
			tr.Td().Span(components.ValueToString(reflect.ValueOf(model).FieldByName(colTag)))
		}
	}
	return t.root
}

func GetStructFieldNames(reflModel reflect.Type) []string {
	var rowNames []string
	InlineLoopFields(reflModel, func(field reflect.StructField) {
		var tag = field.Tag.Get("table")
		rowNames = append(rowNames, tag)
	})
	return rowNames
}

func InlineLoopFields(reflModel reflect.Type, callback func(field reflect.StructField)) {
	for i := 0; i < reflModel.NumField(); i++ {
		var field = reflModel.Field(i)
		if !isValidField(field) {
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			InlineLoopFields(field.Type, callback)
		} else {
			callback(field)
		}
	}
}

func isValidField(field reflect.StructField) bool {
	// Get the struct tags for the field
	var tag = field.Tag.Get("table")
	// Check if the field is ignored
	if tag == "-" || tag == "" {
		return false
	}
	return true
}
