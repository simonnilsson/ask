// Package ask provides a simple way of accessing nested properties in maps and arrays.
// Works great in combination with encoding/json and other packages that "Unmarshal" arbitrary data into Go data-types.
// Inspired by the get function in the lodash javascript library.
package ask

import (
	"errors"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var tokenMatcher = regexp.MustCompile(`([^[]+)?(?:\[(\d+)])?`)
var mapType = reflect.TypeOf(map[string]interface{}{})
var sliceType = reflect.TypeOf([]interface{}{})

var ErrNotFound = errors.New("not found")
var ErrWrongType = errors.New("wrong type")

// Answer holds result of call to For, use one of its methods to extract a value.
type Answer struct {
	value interface{}
}

func handleIntPart(current interface{}, part int) (interface{}, error) {
	val := reflect.ValueOf(current)
	if val.IsValid() && val.CanConvert(sliceType) {
		s := val.Convert(sliceType).Interface().([]interface{})
		if part >= 0 && part < len(s) {
			return s[part], nil
		}
	}
	return current, ErrNotFound
}

func handleStringPart(current interface{}, part string) (interface{}, error) {
	var err error = nil
	match := tokenMatcher.FindStringSubmatch(strings.TrimSpace(part))

	if len(match) == 3 {

		if match[1] != "" {
			val := reflect.ValueOf(current)
			if val.IsValid() && val.CanConvert(mapType) {
				current = val.Convert(mapType).Interface().(map[string]interface{})[match[1]]
			} else {
				err = ErrNotFound
			}
		}

		if match[2] != "" {
			index, _ := strconv.Atoi(match[2])
			return handleIntPart(current, index)
		}

	}

	return current, err
}

// For is used to select a path from source to return as answer.
func For(source interface{}, path string) *Answer {

	parts := strings.Split(path, ".")
	var err error = nil
	current := source

	for _, part := range parts {
		current, err = handleStringPart(current, part)
		if err != nil {
			return &Answer{}
		}
	}

	return &Answer{value: current}
}

// ForArgs is used to select a path using individual arguments from source to return as answer.
func ForArgs(source interface{}, parts ...interface{}) *Answer {

	current := source
	var err error = nil

	for _, part := range parts {

		switch vt := part.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
			index := reflect.ValueOf(vt).Int()
			current, err = handleIntPart(current, int(index))
			if err != nil {
				return &Answer{}
			}
		case string:
			current, err = handleStringPart(current, vt)
			if err != nil {
				return &Answer{}
			}
		}

	}

	return &Answer{value: current}
}

// Path does the same thing as For but uses existing answer as source.
func (a *Answer) Path(path string) *Answer {
	return For(a.value, path)
}

// PathArgs does the same thing as ForArgs but uses existing answer as source.
func (a *Answer) PathArgs(parts ...interface{}) *Answer {
	return ForArgs(a.value, parts...)
}

// Exists returns a boolean indicating if the answer exists (not nil).
func (a *Answer) Exists() bool {
	return a.value != nil
}

// Value returns the raw value as type interface{}, can be nil if no value is available.
func (a *Answer) Value() interface{} {
	return a.value
}

// Slice attempts asserting answer as a []interface{}.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Slice(d []interface{}) ([]interface{}, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	val := reflect.ValueOf(a.value)
	if val.IsValid() && val.CanConvert(sliceType) {
		return val.Convert(sliceType).Interface().([]interface{}), nil
	}
	return d, ErrWrongType
}

// Map attempts asserting answer as a map[string]interface{}.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Map(d map[string]interface{}) (map[string]interface{}, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	val := reflect.ValueOf(a.value)
	if val.IsValid() && val.CanConvert(mapType) {
		return val.Convert(mapType).Interface().(map[string]interface{}), nil
	}
	return d, ErrWrongType
}

// String attempts asserting answer as a string.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) String(d string) (string, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	res, ok := a.value.(string)
	if ok {
		return res, nil
	}
	return d, ErrWrongType
}

// Bool attempts asserting answer as a bool.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Bool(d bool) (bool, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	res, ok := a.value.(bool)
	if ok {
		return res, nil
	}
	return d, ErrWrongType
}

// Int attempts asserting answer as a int64. Casting from other number types will be done if necessary.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Int(d int64) (int64, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	switch vt := a.value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(vt).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		val := reflect.ValueOf(vt).Uint()
		if val <= math.MaxInt64 {
			return int64(val), nil
		}
	case float32, float64:
		val := reflect.ValueOf(vt).Float()
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			return int64(val), nil
		}
	}
	return d, ErrWrongType
}

// Uint attempts asserting answer as a uint64. Casting from other number types will be done if necessary.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Uint(d uint64) (uint64, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	switch vt := a.value.(type) {
	case int, int8, int16, int32, int64:
		val := reflect.ValueOf(vt).Int()
		if val >= 0 {
			return uint64(val), nil
		}
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(vt).Uint(), nil
	case float32, float64:
		val := reflect.ValueOf(vt).Float()
		if val >= 0 && val <= math.MaxUint64 {
			return uint64(val), nil
		}
	}
	return d, ErrWrongType
}

// Float attempts asserting answer as a float64. Casting from other number types will be done if necessary.
// The first return value is the result, and the second indicates if the operation was successful.
// If not successful the first return value will be set to the d parameter.
func (a *Answer) Float(d float64) (float64, error) {
	if a.value == nil {
		return d, ErrNotFound
	}
	switch vt := a.value.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(vt).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(vt).Uint()), nil
	case float32:
		return float64(vt), nil
	case float64:
		return vt, nil
	}
	return d, ErrWrongType
}
