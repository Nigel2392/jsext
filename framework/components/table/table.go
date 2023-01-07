//go:build js && wasm
// +build js,wasm

package table

import (
	"fmt"
	"reflect"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/framework/elements"
	"github.com/Nigel2392/jsext/framework/helpers"
)

type NoStructAvailable int

type Table[T any] struct {
	model          []T
	additionalCols map[string]func(model T) *elements.Element
	root           *elements.Element
	width          string
}

func New(width string) *Table[NoStructAvailable] {
	var t = Table[NoStructAvailable]{
		root:           elements.Div().AttrClass("jsext-table-root").AttrStyle("width:" + width),
		additionalCols: map[string]func(model NoStructAvailable) *elements.Element{},
		width:          width,
	}
	return &t
}

func (t *Table[T]) hasStruct() bool {
	return reflect.TypeOf(t.model).Kind() != reflect.Slice || reflect.TypeOf(t.model).Elem().Kind() != reflect.TypeOf(NoStructAvailable(0)).Kind()
}

func NewFromStruct[T any](width string, s []T, additionalCols map[string]func(model T) *elements.Element) *Table[T] {
	var cols = make(map[string]func(model T) *elements.Element, 0)
	if additionalCols != nil {
		cols = additionalCols
	}
	var t = Table[T]{
		model:          s,
		additionalCols: cols,
		root:           elements.Div().AttrClass("jsext-table-root").AttrStyle("width:" + width),
		width:          width,
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
	// Get the field naems from T
	var m T
	kind := reflect.TypeOf(m).Kind()
	if kind == reflect.Ptr {
		kind = reflect.TypeOf(m).Elem().Kind()
	}
	if kind != reflect.Struct {
		panic("model must be a struct")
	}
	var reflModel = reflect.TypeOf(m)
	var rowNames = GetStructFieldNames(reflModel)
	var table = t.root.Table().AttrClass("jsext-table").AttrStyle("width:" + t.width)
	var thead = table.Thead()
	var tbody = table.Tbody()
	var tr = thead.Tr()
	for _, rowName := range rowNames {
		tr.Th().AttrStyle("width:"+rowName.Width, "text-align:"+rowName.TextAlign).Span(rowName.Name)
	}

	for _, model := range t.model {
		tr = tbody.Tr()

		kind := reflect.TypeOf(model).Kind()
		if kind == reflect.Ptr {
			kind = reflect.TypeOf(model).Elem().Kind()
		}
		if kind != reflect.Struct {
			panic("model must be a struct")
		}
		var valueModel = reflect.TypeOf(model)
		var i = 0
		helpers.InlineLoopFields(valueModel, func(field reflect.StructField, parent reflect.Type, value reflect.Value) {
			// Get the value of the field
			var val = reflect.ValueOf(model).FieldByName(field.Name).Interface()
			var width = rowNames[i].Width
			var textAlign = rowNames[i].TextAlign
			var valueString = fmt.Sprintf("%v", val)
			tr.Td().AttrStyle("width:"+width, "text-align:"+textAlign).Span(valueString)
			i++
		})
	}
	return t.root
}

type Rows struct {
	Width     string
	Name      string
	TextAlign string
}

func GetStructFieldNames(reflModel reflect.Type) []Rows {
	var rowNames []Rows
	helpers.InlineLoopFields(reflModel, func(field reflect.StructField, parent reflect.Type, value reflect.Value) {
		var tag = field.Tag.Get("table")
		var widthTag = field.Tag.Get("width")
		var width = "auto"
		if widthTag != "" {
			width = widthTag
		}
		var textAlignTag = field.Tag.Get("align")
		var textAlign = "left"
		if textAlignTag != "" {
			textAlign = textAlignTag
		}
		rowNames = append(rowNames, Rows{
			Width:     width,
			Name:      tag,
			TextAlign: textAlign,
		})
	})
	return rowNames
}

func (t *Table[T]) Run() *elements.Element {
	return t.create()
}

// # Example:
//
//	var resultList = make([][]string, len(results)+1)
//	resultList[0] = []string{"Title", "Content", "Author", "Created", "Updated"}
//
//	for i, result := range results {
//		var innerList = make([]string, 5)
//		for j, v := range []string{result.Title, result.Body, result.Author.Username, result.CreatedAt().Format("2006-01-02 15:04:05"), result.UpdatedAt().Format("2006-01-02 15:04:05")} {
//			innerList[j] = v
//		}
//		resultList[i+1] = innerList
//	}
//
//	return table.NewListTable("100%", resultList)
type ListTable struct {
	innerList [][]string
	root      *elements.Element
	width     string
}

func NewListTable(width string, list [][]string) *ListTable {
	var t = ListTable{
		root:      elements.Div().AttrClass("jsext-table-root").AttrStyle("width:" + width),
		width:     width,
		innerList: list,
	}
	return &t
}

func (t *ListTable) Render() jsext.Element {
	return t.create().Render()
}

func (t *ListTable) create() *elements.Element {
	var table = t.root.Table().AttrClass("jsext-table").AttrStyle("width:" + t.width)
	var thead = table.Thead()
	var tbody = table.Tbody()
	var tr = thead.Tr()
	for _, rowName := range t.innerList[0] {
		tr.Th().AttrStyle("width:auto", "text-align:left").Span(fmt.Sprintf("%v", rowName))
	}

	for _, model := range t.innerList[1:] {
		tr = tbody.Tr()
		for _, val := range model {
			tr.Td().AttrStyle("width:auto", "text-align:left").Span(fmt.Sprintf("%v", val))
		}
	}
	return t.root
}
