package utils

import (
	"errors"
	"math"
	"strconv"
)

// FixFloat returns fixed digits precision float64 value
func FixFloat(f float64, digit ...int) float64 {
	var useDigit = 2
	if len(digit) > 0 {
		useDigit = digit[0]
	}
	pow := math.Pow10(useDigit)
	i := int64(f * pow)
	return float64(i) / pow
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// FixFloatAndRound returns fixed digits precision and rounded float64 value
func FixFloatAndRound(f float64, digit ...int) float64 {
	var useDigit = 2
	if len(digit) > 0 {
		useDigit = digit[0]
	}
	pow := math.Pow10(useDigit)
	return float64(round(f*pow)) / pow
}

// FixZero returns 0 if lower than 0
func FixZero(f float64) float64 {
	if f < 0 {
		return 0
	}
	return f
}

// ToFloat converts interface to float64 value
func ToFloat(v interface{}) (float64, error) {
	if v == nil {
		return 0, errors.New("nil-value")
	}
	var (
		floatValue float64
		err        error
	)
	isValid := true
	switch value := v.(type) {
	case string:
		if floatValue, err = strconv.ParseFloat(value, 64); err != nil {
			isValid = false
		}
	case float64:
		floatValue = value
	case float32:
		floatValue = float64(value)
	case int:
		floatValue = float64(value)
	case int32:
		floatValue = float64(value)
	case int64:
		floatValue = float64(value)
	case uint:
		floatValue = float64(value)
	case uint32:
		floatValue = float64(value)
	case uint64:
		floatValue = float64(value)
	default:
		isValid = false
	}
	if !isValid {
		return 0, errors.New("wrong-value")
	}
	return floatValue, nil
}
