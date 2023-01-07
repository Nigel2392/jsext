package convert

import "strconv"

func FormatNumber(value any) string {
	switch value.(type) {
	case int:
		return strconv.Itoa(value.(int))
	case int8:
		return strconv.Itoa(int(value.(int8)))
	case int16:
		return strconv.Itoa(int(value.(int16)))
	case int32:
		return strconv.Itoa(int(value.(int32)))
	case int64:
		return strconv.Itoa(int(value.(int64)))
	case uint:
		return strconv.Itoa(int(value.(uint)))
	case uint8:
		return strconv.Itoa(int(value.(uint8)))
	case uint16:
		return strconv.Itoa(int(value.(uint16)))
	case uint32:
		return strconv.Itoa(int(value.(uint32)))
	case uint64:
		return strconv.Itoa(int(value.(uint64)))
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	default:
		panic("Value must be a number type")
	}
}

func ToFloat(value any) float64 {
	switch value.(type) {
	case int:
		return float64(value.(int))
	case int8:
		return float64(value.(int8))
	case int16:
		return float64(value.(int16))
	case int32:
		return float64(value.(int32))
	case int64:
		return float64(value.(int64))
	case uint:
		return float64(value.(uint))
	case uint8:
		return float64(value.(uint8))
	case uint16:
		return float64(value.(uint16))
	case uint32:
		return float64(value.(uint32))
	case uint64:
		return float64(value.(uint64))
	case float32:
		return float64(value.(float32))
	case float64:
		return value.(float64)
	default:
		panic("Value must be a number type")
	}
}
