package app

import (
	"fmt"
	"strconv"
	"time"
)

// Application datamap to store and retrieve data.
type DataMap map[string]any

// Set data on the data map.
func (d DataMap) Set(k string, v any) {
	d[k] = v
}

// Get data from the data map.
func (d DataMap) Get(key string) interface{} {
	return d[key]
}

// Get data from the data map in the form of an int.
func (d DataMap) GetInt(key string) int {
	var data = d[key]
	switch data := data.(type) {
	case int:
		return data
	case int8:
		return int(data)
	case int16:
		return int(data)
	case int32:
		return int(data)
	case int64:
		return int(data)
	}
	return 0
}

// Get data from the data map in the form of a uint.
func (d DataMap) GetUint(key string) uint {
	var data = d[key]
	switch data := data.(type) {
	case uint:
		return data
	case uint8:
		return uint(data)
	case uint16:
		return uint(data)
	case uint32:
		return uint(data)
	case uint64:
		return uint(data)
	}
	return 0
}

// Get data from the data map in the form of a string.
func (d DataMap) GetString(key string) string {
	if d, ok := d[key]; ok {
		if s, ok := d.(string); ok {
			return s
		}
		return formatString(d)
	}
	return ""
}

// Get data from the data map in the form of a bool.
func (d DataMap) GetBool(key string) bool {
	if d, ok := d[key]; ok {
		if b, ok := d.(bool); ok {
			return b
		}
		return false
	}
	return false
}

// Get data from the data map in the form of a float64.
func (d DataMap) GetFloat(key string) float64 {
	var data = d[key]
	switch data := data.(type) {
	case float64:
		return data
	case float32:
		return float64(data)
	}
	return 0
}

// Get data from the data map in the form of a complex128.
func (d DataMap) GetComplex(key string) complex128 {
	var data = d[key]
	switch data := data.(type) {
	case complex128:
		return data
	case complex64:
		return complex128(data)
	}
	return 0
}

// Get data from the data map in the form of a time.Time.
func (d DataMap) GetTime(key string) time.Time {
	var data = d[key]
	switch data := data.(type) {
	case time.Time:
		return data
	}
	return time.Time{}
}

// Get data from the data map in the form of a time.Duration.
func (d DataMap) GetDuration(key string) time.Duration {
	var data = d[key]
	switch data := data.(type) {
	case time.Duration:
		return data
	}
	return time.Duration(0)
}

// Get data from the data map in the form of a []byte.
func (d DataMap) GetBytes(key string) []byte {
	var data = d[key]
	switch data := data.(type) {
	case []byte:
		return data
	}
	return nil
}

// Get data from the data map in the form of a []rune.
func (d DataMap) GetRunes(key string) []rune {
	var data = d[key]
	switch data := data.(type) {
	case []rune:
		return data
	}
	return nil
}

func formatString(data interface{}) string {
	switch data := data.(type) {
	case string:
		return data
	case int:
		return strconv.Itoa(data)
	case int8:
		return strconv.Itoa(int(data))
	case int16:
		return strconv.Itoa(int(data))
	case int32:
		return strconv.Itoa(int(data))
	case int64:
		return strconv.Itoa(int(data))
	case uint:
		return strconv.FormatUint(uint64(data), 10)
	case uint8:
		return strconv.FormatUint(uint64(data), 10)
	case uint16:
		return strconv.FormatUint(uint64(data), 10)
	case uint32:
		return strconv.FormatUint(uint64(data), 10)
	case uint64:
		return strconv.FormatUint(data, 10)
	case float32:
		return strconv.FormatFloat(float64(data), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data, 'f', -1, 64)
	case complex64:
		return strconv.FormatComplex(complex128(data), 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(data, 'f', -1, 128)
	case bool:
		return strconv.FormatBool(data)
	case []byte:
		return string(data)
	case []rune:
		return string(data)
	case time.Time:
		return data.Format(time.RFC3339)
	case time.Duration:
		return data.String()
	case error:
		return data.Error()
	case fmt.Stringer:
		return data.String()
	default:
		panic("cannot format this type to string.")
	}
}

// Get any type from the datamap.
// Panics if the type is not correct.
// Returns a new instance of the type if the key is not found.
func GetType[T any](d DataMap, key string) (T, bool) {
	if d, ok := d[key]; ok {
		if s, ok := d.(T); ok {
			return s, true
		}
		panic("cannot convert to type")
	}
	return *new(T), false // return a new instance of the type
}
