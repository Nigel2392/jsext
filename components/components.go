package components

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Nigel2392/jsext"
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

func ValueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(v.Complex(), 'f', -1, 64)
	default:
		// Time
		if v.Type().String() == "time.Time" {
			return v.Interface().(time.Time).Format("2006-01-02T15:04")
		}
		return ""
	}
}

func ReflectModelName(s any) string {
	var structName string
	kind := reflect.TypeOf(s).Elem().Kind()
	if kind == reflect.Pointer {
		structName = reflect.TypeOf(s).Elem().Elem().Name()
	} else if kind == reflect.Struct {
		structName = reflect.TypeOf(s).Elem().Name()
	}
	structName = toValidName(structName)
	return structName
}

func toValidName(name string) string {
	var s = strings.ReplaceAll(name, "_", "")
	return s
}

// Get the kind of the model (Reflect.TYPE)
func StructKind(s any) reflect.Type {
	// Validate kind
	kind := reflect.TypeOf(s)
	if kind.Kind() == reflect.Ptr {
		kind = kind.Elem()
	}
	if kind.Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	return kind
}

// Get a value from a model struct
func GetValue(s any, column string) any {
	// Validate kind
	kind := StructKind(s)
	// Loop through all fields in the struct
	for i := 0; i < kind.NumField(); i++ {
		f_kind := kind.Field(i)
		// Get the name of the struct field
		if strings.EqualFold(f_kind.Name, column) {
			var val any
			// Get the value of the field
			if f_kind.Type.Kind() == reflect.Ptr {
				val = reflect.ValueOf(s).Elem().Field(i).Elem().Interface()
			} else {
				val = reflect.ValueOf(s).Elem().Field(i).Interface()
			}
			return val
		}
	}
	return nil
}

func InlineLoopFields(reflModel reflect.Type, callback func(field reflect.StructField, parent reflect.Type, value reflect.Value), tags ...string) {
	for i := 0; i < reflModel.NumField(); i++ {
		var field = reflModel.Field(i)
		if !isValidField(field, tags...) {
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			InlineLoopFields(field.Type, callback)
		} else {
			callback(field, reflModel, reflect.ValueOf(field))
		}
	}
}

func isValidField(field reflect.StructField, tags ...string) bool {
	for _, tag := range tags {
		if field.Tag.Get(tag) == "-" {
			return false
		}
	}
	return true
}
