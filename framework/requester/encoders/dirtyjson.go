package encoders

import (
	"bytes"
	"strconv"
	"strings"
)

// This package includes a json encoder for maps and lists.
// Decoding is not supported.

func MarshalMap(data map[string]interface{}) []byte {
	var b bytes.Buffer
	marshalMapToBuffer(data, &b)
	return b.Bytes()
}

func MarshalList(data []interface{}) []byte {
	var b bytes.Buffer
	marshalListToBuffer(data, &b)
	return b.Bytes()
}

func formatStringSafe(value interface{}) string {
	return "\"" + strings.ReplaceAll(value.(string), `"`, `\"`) + "\""
}

func formatString(value interface{}) string {
	return value.(string)
}

func formatInt(value interface{}) string {
	switch val := value.(type) {
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.Itoa(int(val))
	case int16:
		return strconv.Itoa(int(val))
	case int32:
		return strconv.Itoa(int(val))
	case int64:
		return strconv.Itoa(int(val))
	default:
		return ""
	}
}

func formatUint(value interface{}) string {
	switch val := value.(type) {
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	default:
		return ""
	}
}

func formatFloat(value interface{}) string {
	switch val := value.(type) {
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return ""
	}
}

func formatComplex(value interface{}) string {
	switch val := value.(type) {
	case complex64:
		return strconv.FormatComplex(complex128(val), 'f', -1, 64)
	case complex128:
		return strconv.FormatComplex(val, 'f', -1, 128)
	default:
		return ""
	}
}

func formatBool(value interface{}) string {
	return strconv.FormatBool(value.(bool))
}

func StringMapToAnyMap(input map[string]string) map[string]interface{} {
	var output = make(map[string]interface{}, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func marshalListToBuffer(data []interface{}, b *bytes.Buffer) {
	b.WriteString("[")
	for i, value := range data {
		//lint:ignore S1034 this is a switch statement
		switch value.(type) {
		case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float64, float32, complex64, complex128, bool, []byte, []rune:
			b.WriteString(toString(value))
		case map[string]interface{}:
			marshalMapToBuffer(value.(map[string]interface{}), b)
		case []interface{}:
			marshalListToBuffer(value.([]interface{}), b)
		case nil:
			b.WriteString("null")
		default:
			b.WriteString("null")
		}
		if i < len(data)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("]")
}

func marshalMapToBuffer(data map[string]interface{}, b *bytes.Buffer) {
	b.WriteString("{")

	var i = 0
	for key, value := range data {
		//lint:ignore S1034 this is a switch statement
		switch value.(type) {
		case string:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatStringSafe(value))
		case int, int8, int16, int32, int64:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatInt(value))
		case uint, uint8, uint16, uint32, uint64:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatUint(value))
		case float64, float32:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatFloat(value))
		case complex64, complex128:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatComplex(value))
		case bool:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatBool(value))
		case map[string]interface{}:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			marshalMapToBuffer(value.(map[string]interface{}), b)
		case []interface{}:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			marshalListToBuffer(value.([]interface{}), b)
		case []byte:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString(formatStringSafe(string(value.([]byte))))
		case nil:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString("null")
		default:
			b.WriteString(formatStringSafe(key))
			b.WriteString(":")
			b.WriteString("null")
		}
		if i < len(data)-1 {
			b.WriteString(",")
		}
		i++
	}
	b.WriteString("}")
}

func toString(data interface{}) string {
	// string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float64, float32, complex64, complex128, bool
	switch data := data.(type) {
	case bool:
		return formatBool(data)
	case int, int8, int16, int32, int64:
		return formatInt(data)
	case uint, uint8, uint16, uint32, uint64:
		return formatUint(data)
	case float32, float64:
		return formatFloat(data)
	case complex128, complex64:
		return formatComplex(data)
	case string:
		return formatString(data)
	case []byte:
		return formatString(string(data))
	case []rune:
		return formatString(string(data))
	default:
		return ""
	}
}
