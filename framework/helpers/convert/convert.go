package convert

import "strconv"

func FormatNumber(value any) string {
	switch value := value.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.Itoa(int(value))
	case uint:
		return strconv.Itoa(int(value))
	case uint8:
		return strconv.Itoa(int(value))
	case uint16:
		return strconv.Itoa(int(value))
	case uint32:
		return strconv.Itoa(int(value))
	case uint64:
		return strconv.Itoa(int(value))
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		panic("Value must be a number type")
	}
}

func ToFloat(value any) float64 {
	switch value := value.(type) {
	case int:
		return float64(value)
	case int8:
		return float64(value)
	case int16:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	case uint:
		return float64(value)
	case uint8:
		return float64(value)
	case uint16:
		return float64(value)
	case uint32:
		return float64(value)
	case uint64:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	default:
		panic("Value must be a number type")
	}
}
