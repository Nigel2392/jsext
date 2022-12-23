package forms

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Nigel2392/jsext"
	"github.com/Nigel2392/jsext/elements"
	"github.com/Nigel2392/jsext/helpers"
)

// Form delimiter for nested structs
const Delimiter = "___"

// Form struct
type Form struct {
	Inner      *elements.Element
	Validators map[string]func(string) error
}

// Get the form value from the form element
func (f *Form) Value() jsext.Element {
	return f.Inner.JSExtElement()
}

// Get the form value from the form element
func (f *Form) Element() *elements.Element {
	return f.Inner
}

// Render the form
func (f *Form) Render() jsext.Element {
	return f.Inner.Render()
}

// Set the form ID
func (f *Form) AttrID(id string) *Form {
	f.Inner.AttrID(id)
	return f
}

// Add a validator to the form (NYI fully, works for now)
func (f *Form) AddValidator(name string, fn func(string) error) {
	f.Validators[name] = fn
}

// Eventlistener for when the form is submitted.
func (f *Form) OnSubmit(cb func(data map[string]string, elements []jsext.Element)) {
	f.Inner.AddEventListener("submit", func(this jsext.Value, event jsext.Event) {
		event.PreventDefault()
		var data = make(map[string]string)
		var elements = this.Get("elements")
		var length = elements.Get("length").Int()
		var elemList = make([]jsext.Element, length)
		for i := 0; i < length; i++ {
			var element = elements.Index(i)
			var name = element.Get("name").String()
			var value = element.Get("value").String()
			if name == "" {
				continue
			}

			if fn, ok := f.Validators[name]; ok {
				var err = fn(value)
				if err != nil {
					return
				}
			}

			data[name] = value
			elemList[i] = element.ToElement()
		}
		cb(data, elemList)
	})
}

// Eventlistener for when the form is submitted, data is parsed into a struct
func (f *Form) OnSubmitToStruct(strct any, fn func(strct any, elements []jsext.Element)) {
	f.OnSubmit(func(data map[string]string, elements []jsext.Element) {
		var err = FormDataToStruct(data, strct)
		if err != nil {
			panic(err)
		}
		fn(strct, elements)
	})
}

// Render a struct into a form
func StructToForm(s any, labelClass, inputClass, action, method string) *Form {
	var form = elements.Form(action, method)
	var v = reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("not a struct")
	}
	form.AttrID(fmt.Sprintf("jsext_%s_form", v.Type().Name()))
	// Uses a list in the following format:
	// [[Label, Value, Type], [Label, Value, Type], [Label, Value, Type]]
	var List = createListFromStruct(v, s, "")
	// Create the form
	for _, item := range List {
		var label = item[0]
		var name = label
		var value = item[1]
		var typ = item[2]
		label = strings.ReplaceAll(label, Delimiter, " ")
		label = strings.ReplaceAll(label, "_", " ")
		var elemLabel = elements.Label(label, name)
		var elemInput = elements.Input(typ, name, label).AttrValue(value)
		if inputClass != "" {
			elemInput.AttrClass(inputClass)
		}
		if labelClass != "" {
			elemLabel.AttrClass(labelClass)
		}
		form.Append(elemLabel, elemInput)
	}
	return &Form{Inner: form, Validators: make(map[string]func(string) error)}
}

// Parse form data to a struct.
func FormDataToStruct(data map[string]string, s any) error {
	// Parse the form data into the struct
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("not a struct")
	}
	return formParse(data, v, s)
}

func formParse(data map[string]string, v reflect.Value, s any) error {
	for key, value := range data {
		if field, found := v.Type().FieldByName(key); found {
			var val, err = TransformValue(s, key, value)
			if err != nil {
				return err
			}
			SetValueStrict(s, field, v, key, val)
		} else {
			var keys = strings.Split(key, Delimiter)
			// Parse the inner struct recursively with another function
			var err = recurseKeys(keys, value, v, s, s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Load form data into a struct
//
// Struct looks like this:
//
//	type User struct {
//		Name string
//		Age  int
//		Sub  struct {
//			Name string
//			Age  int
//		}
//	}
//
// Form data looks like this:
//
//	map[name:John age:20 sub_name:John sub_age:20]
func recurseKeys(keys []string, value string, v reflect.Value, s any, parent any) error {
	if len(keys) == 1 {
		// We are at the end of the keys
		// Set the value
		var val, err = TransformValue(s, keys[0], value)
		if err != nil {
			panic(err)
		}
		if field, found := v.Type().FieldByName(keys[0]); found {
			SetValueStrict(s, field, v, keys[0], val)
		}
		return nil
	}
	var i = v.FieldByName(keys[0])
	if !i.IsValid() {
		return fmt.Errorf("field %s not found in %s", keys[0], v.Type())
	}
	// Check if the field is a pointer
	if i.Kind() == reflect.Ptr {
		// Check if the pointer is nil
		if i.IsNil() {
			// Create a new struct
			var newStruct = reflect.New(i.Type().Elem())
			// Set the value of the field to the new struct
			i.Set(newStruct)
			// Set the value of i to the new struct
			i = newStruct.Elem()
		} else {
			// Set the value of i to the struct that the pointer points to
			i = i.Elem()
		}
	} else if i.Kind() == reflect.Struct {
		// We found a struct
		// Recurse
		return recurseKeys(keys[1:], value, i, s, parent)
	}
	// Recurse
	return recurseKeys(keys[1:], value, i, s, parent)
}

func createListFromStruct(v reflect.Value, s any, prefix string) [][]string {
	var list [][]string
	for i := 0; i < v.NumField(); i++ {
		var label = v.Type().Field(i).Name
		var value = helpers.ValueToString(v.Field(i))
		var typ = ReflectInputType(helpers.GetValue(s, label))
		if prefix != "" {
			label = prefix + Delimiter + label
		}
		if v.Field(i).Kind() == reflect.Struct && !isValidTyp(typ) {
			var subList = createListFromStruct(v.Field(i), s, label)
			list = append(list, subList...)
			continue
		}
		list = append(list, []string{label, value, string(typ)})
	}
	return list
}

// Transform a form value to the correct type
func TransformValue(s any, field string, val any) (any, error) {
	mdl_val := helpers.GetValue(s, field)
	switch mdl_val.(type) {
	case int, int8, int16, int32, int64:
		integer, err := strconv.Atoi(fmt.Sprint(val))
		if err != nil {
			return nil, err
		}
		switch mdl_val.(type) {
		case int:
			return int(integer), nil
		case int8:
			return int8(integer), nil
		case int16:
			return int16(integer), nil
		case int32:
			return int32(integer), nil
		case int64:
			return int64(integer), nil
		default:
			return nil, errors.New("unknown integer type")
		}
	case uint, uint8, uint16, uint32, uint64:
		integer, err := strconv.ParseUint(fmt.Sprint(val), 10, 64)
		if err != nil {
			return nil, err
		}
		switch mdl_val.(type) {
		case uint:
			return uint(integer), nil
		case uint8:
			return uint8(integer), nil
		case uint16:
			return uint16(integer), nil
		case uint32:
			return uint32(integer), nil
		case uint64:
			return uint64(integer), nil
		default:
			return nil, errors.New("unknown integer type")
		}
	case float32, float64:
		floaty, err := strconv.ParseFloat(fmt.Sprint(val), 64)
		if err != nil {
			return nil, err
		}
		switch mdl_val.(type) {
		case float32:
			return float32(floaty), nil
		case float64:
			return float64(floaty), nil
		default:
			return nil, errors.New("unknown float type")
		}
	case bool:
		return strconv.ParseBool(fmt.Sprint(val))
	case string:
		return val, nil
	case []byte:
		return val, nil
	case time.Time:
		// This whole function is essentially just for this.
		// Convert HTML local datetime to time object
		t, err := time.Parse("2006-01-02T15:04", fmt.Sprint(val))
		if err != nil {
			return val, errors.New("error parsing time from value: " + fmt.Sprint(val) + " " + err.Error())
		}
		return t, nil
	default:
		return val, nil
	}
}

// Valid formtypes
type FORMTYPES string

// Check if two formtypes are equal
func (ft FORMTYPES) Equals(other FORMTYPES) bool {
	return strings.EqualFold(string(ft), string(other))
}

// Formtypes to use in forms
const (
	FORMTYP_TEXT     FORMTYPES = "text"
	FORMTYP_CHECKBOX FORMTYPES = "checkbox"
	FORMTYP_NUMBER   FORMTYPES = "number"
	FORMTYP_DATETIME FORMTYPES = "datetime-local"
	FORMTYP_FILE     FORMTYPES = "file"
	FORMTYP_STRUCT   FORMTYPES = "struct"
	// FORMTYP_LIST     FORMTYPES = "select"
	FORMTYP_INVALID FORMTYPES = "text"
)

// Reflect the form input type of a value
func ReflectInputType(val any) FORMTYPES {
	switch val.(type) {
	case string:
		return FORMTYP_TEXT
	case bool:
		return FORMTYP_CHECKBOX
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return FORMTYP_NUMBER
	case time.Time:
		return FORMTYP_DATETIME
	case []byte:
		return FORMTYP_FILE
	default:
		newV := reflect.ValueOf(val)
		if newV.Kind() == reflect.Ptr {
			newV = newV.Elem()
		}
		//if newV.Kind() == reflect.Slice || newV.Kind() == reflect.Array {
		//	return FORMTYP_LIST
		//} else
		if newV.Kind() == reflect.Struct {
			return FORMTYP_STRUCT
		} else {
			return FORMTYP_INVALID
		}
	}
}

func isValidTyp(typ FORMTYPES) bool {
	switch typ {
	case FORMTYP_TEXT, FORMTYP_CHECKBOX, FORMTYP_NUMBER, FORMTYP_DATETIME, FORMTYP_FILE: //, FORMTYP_LIST:
		return true
	default:
		return false
	}
}

// Format a value for display in a form
func FormatIfDateTime(val any) any {
	switch val := val.(type) {
	case time.Time:
		return val.Format("2006-01-02T15:04")
	default:
		return val
	}
}

// Set a value on a struct
func SetValue(s any, column string, value any) {
	// Validate kind
	kind := helpers.StructKind(s)
	// Loop through all fields in the struct
	for i := 0; i < kind.NumField(); i++ {
		f_kind := kind.Field(i)
		// Get the name of the struct field
		if strings.EqualFold(f_kind.Name, column) {
			// Set the value of the struct field
			// Check if types match
			reflect.ValueOf(s).Elem().Field(i).Set(reflect.ValueOf(value))
			return
		}
	}
}

// Set a value on a struct, faster option.
func SetValueStrict(s any, f reflect.StructField, v reflect.Value, column string, value any) {
	newVal := reflect.ValueOf(value)
	// Check if the field is a pointer
	if f.Type.Kind() == reflect.Ptr {
		// Check if the value is a pointer
		if newVal.Kind() == reflect.Ptr {
			// Set the value
			v.FieldByName(column).Set(newVal)
		} else {
			// Create a pointer to the value
			var newV = reflect.New(newVal.Type())
			newV.Elem().Set(newVal)
			// Set the value
			v.FieldByName(column).Set(newV)
		}
	} else {
		// Set the value
		v.FieldByName(column).Set(newVal)
	}
}
