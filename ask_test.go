package ask

import (
	"math"
	"reflect"
	"testing"
)

func TestFor(t *testing.T) {

	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": 100,
			},
		},
	}

	// OK
	answer := For(source, "a[0].b")
	if answer.value != 100 {
		t.Errorf(`For() = (%v); want (100)`, answer.value)
	}

	// Missing
	answer = For(source, "a[1].b")
	if answer.value != nil {
		t.Errorf(`For() = (%v); want (nil)`, answer.value)
	}

	// Missing slice
	answer = For(source, "d[1]")
	if answer.value != nil {
		t.Errorf("For() = (%v); want (nil)", answer.value)
	}

	// Empty path should return source
	answer = For(source, "")
	if !reflect.DeepEqual(answer.value, source) {
		t.Errorf(`For() = (%v); want (source)`, answer.value)
	}

	// Invalid path should return nil value
	answer = For(source, "---")
	if answer.value != nil {
		t.Errorf(`For() = (%v); want (nil)`, answer.value)
	}

}

func TestForPath(t *testing.T) {

	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b":   100,
				"c.d": true,
			},
		},
	}

	// OK
	answer := ForArgs(source, "a", 0, "b")
	if answer.value != 100 {
		t.Errorf(`ForArgs() = (%v); want (100)`, answer.value)
	}

	// Missing
	answer = ForArgs(source, "a", 1, "b")
	if answer.value != nil {
		t.Errorf(`ForArgs() = (%v); want (nil)`, answer.value)
	}

	// Missing slice
	answer = ForArgs(source, "d", 1)
	if answer.value != nil {
		t.Errorf("ForArgs() = (%v); want (nil)", answer.value)
	}

	// Part with dot in field name should return value
	answer = ForArgs(source, "a", 0, "c.d")
	if answer.value != true {
		t.Errorf(`ForArgs() = (%v); want (true)`, answer.value)
	}

	// Empty path should return source
	answer = ForArgs(source, "")
	if !reflect.DeepEqual(answer.value, source) {
		t.Errorf(`ForArgs() = (%v); want (source)`, answer.value)
	}

	// Invalid path should return nil value
	answer = ForArgs(source, "---")
	if answer.value != nil {
		t.Errorf(`ForArgs() = (%v); want (nil)`, answer.value)
	}

}

func TestPath(t *testing.T) {

	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": 100,
			},
		},
	}

	// OK
	answer := For(source, "a[0]").Path("b")
	if answer.value != 100 {
		t.Errorf(`Path() = (%v); want (100)`, answer.value)
	}

	// Missing
	answer = For(source, "a[1]").Path("b")
	if answer.value != nil {
		t.Errorf(`Path() = (%v); want (nil)`, answer.value)
	}

	// Empty path should return source
	answer = For(source, "").Path("")
	if !reflect.DeepEqual(answer.value, source) {
		t.Errorf(`Path() = (%v); want (source)`, answer.value)
	}

}

func TestPathArgs(t *testing.T) {

	source := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"b": 100,
			},
		},
	}

	// OK
	answer := For(source, "a").PathArgs(0, "b")
	if answer.value != 100 {
		t.Errorf(`Path() = (%v); want (100)`, answer.value)
	}

	// Missing
	answer = For(source, "a[1]").PathArgs("b")
	if answer.value != nil {
		t.Errorf(`Path() = (%v); want (nil)`, answer.value)
	}

	// Empty path should return source
	answer = For(source, "").PathArgs("")
	if !reflect.DeepEqual(answer.value, source) {
		t.Errorf(`Path() = (%v); want (source)`, answer.value)
	}

}

func TestString(t *testing.T) {

	source := map[string]interface{}{
		"string": "test",
		"number": 100,
	}

	// OK
	res, err := For(source, "string").String("default")
	if err != nil || res != "test" {
		t.Errorf(`String() = ("%s", %v); want ("test", <nil>)`, res, err)
	}

	// Wrong type
	res, err = For(source, "number").String("default")
	if err != ErrWrongType || res != "default" {
		t.Errorf(`String() = ("%s", %v); want ("default", wrong type)`, res, err)
	}

	// Missing
	res, err = For(source, "nothing").String("default")
	if err != ErrNotFound || res != "default" {
		t.Errorf(`String() = ("%s", %v); want ("default", not found)`, res, err)
	}

}

func TestBool(t *testing.T) {

	source := map[string]interface{}{
		"bool1":  true,
		"bool2":  false,
		"number": 100,
	}

	// OK
	res, err := For(source, "bool1").Bool(false)
	if err != nil || res != true {
		t.Errorf(`Bool() = (%t, %v); want (true, <nil>)`, res, err)
	}

	// OK
	res, err = For(source, "bool2").Bool(false)
	if err != nil || res != false {
		t.Errorf(`Bool() = (%t, %v); want (false, <nil>)`, res, err)
	}

	// Wrong type
	res, err = For(source, "number").Bool(false)
	if err != ErrWrongType || res != false {
		t.Errorf(`Bool() = (%t, %v); want (false, wrong type)`, res, err)
	}

	// Missing
	res, err = For(source, "nothing").Bool(false)
	if err != ErrNotFound || res != false {
		t.Errorf(`Bool() = (%t, %v); want (false, not found)`, res, err)
	}

}

func TestInt(t *testing.T) {

	source := map[string]interface{}{
		"positive":       100,
		"negative":       -100,
		"unsigned":       uint(100),
		"float":          float64(100),
		"negative_float": float64(-100),
		"toobig":         uint64(math.MaxUint64),
		"string":         "test",
	}

	// OK
	res, err := For(source, "positive").Int(5)
	if err != nil || res != 100 {
		t.Errorf("Int() = (%d, %v); want (100, <nil>)", res, err)
	}

	// OK negative
	res, err = For(source, "negative").Int(5)
	if err != nil || res != -100 {
		t.Errorf("Int() = (%d, %v); want (-100, <nil>)", res, err)
	}

	// OK unsigned
	res, err = For(source, "unsigned").Int(5)
	if err != nil || res != 100 {
		t.Errorf("Int() = (%d, %v); want (100, <nil>)", res, err)
	}

	// OK float
	res, err = For(source, "float").Int(5)
	if err != nil || res != 100 {
		t.Errorf("Int() = (%d, %v); want (100, <nil>)", res, err)
	}

	// OK negative float
	res, err = For(source, "negative_float").Int(5)
	if err != nil || res != -100 {
		t.Errorf("Int() = (%d, %v); want (-100, <nil>)", res, err)
	}

	// Too big number
	res, err = For(source, "toobig").Int(5)
	if err != ErrWrongType || res != 5 {
		t.Errorf("Int() = (%d, %v); want (5, wrong type)", res, err)
	}

	// Wrong type
	res, err = For(source, "string").Int(5)
	if err != ErrWrongType || res != 5 {
		t.Errorf("Int() = (%d, %v); want (5, wrong type)", res, err)
	}

	// Missing
	res, err = For(source, "nothing").Int(5)
	if err != ErrNotFound || res != 5 {
		t.Errorf("Int() = (%d, %v); want (5, not found)", res, err)
	}

}

func TestUint(t *testing.T) {

	source := map[string]interface{}{
		"positive": 100,
		"negative": -100,
		"unsigned": uint(100),
		"float":    float64(100),
		"string":   "test",
	}

	// OK
	res, err := For(source, "positive").Uint(5)
	if err != nil || res != 100 {
		t.Errorf("Uint() = (%d, %v); want (100, <nil>)", res, err)
	}

	// OK unsigned
	res, err = For(source, "unsigned").Uint(5)
	if err != nil || res != 100 {
		t.Errorf("Uint() = (%d, %v); want (100, <nil>)", res, err)
	}

	// OK float
	res, err = For(source, "float").Uint(5)
	if err != nil || res != 100 {
		t.Errorf("Uint() = (%d, %v); want (100, <nil>)", res, err)
	}

	// Fail on negative
	res, err = For(source, "negative").Uint(5)
	if err != ErrWrongType || res != 5 {
		t.Errorf("Uint() = (%d, %v); want (5, wrong type)", res, err)
	}

	// Wrong type
	res, err = For(source, "string").Uint(5)
	if err != ErrWrongType || res != 5 {
		t.Errorf("Uint() = (%d, %v); want (5, wrong type)", res, err)
	}

	// Missing
	res, err = For(source, "nothing").Uint(5)
	if err != ErrNotFound || res != 5 {
		t.Errorf("Uint() = (%d, %v); want (5, not found)", res, err)
	}

}

func TestFloat(t *testing.T) {

	source := map[string]interface{}{
		"positive": 100,
		"negative": -100,
		"unsigned": uint(100),
		"float32":  float32(100.1),
		"float64":  float64(100.1),
		"string":   "test",
	}

	// OK
	res, err := For(source, "positive").Float(5)
	if err != nil || res != 100 {
		t.Errorf("Float() = (%f, %v); want (100, <nil>)", res, err)
	}

	// OK negative
	res, err = For(source, "negative").Float(5)
	if err != nil || res != -100 {
		t.Errorf("Float() = (%f, %v); want (-100, <nil>)", res, err)
	}

	// OK unsigned
	res, err = For(source, "unsigned").Float(5)
	if err != nil || res != 100 {
		t.Errorf("Float() = (%f, %v); want (100, <nil>)", res, err)
	}

	// OK float32
	res, err = For(source, "float32").Float(5)
	if err != nil || math.Abs(res-100.1) > .00001 {
		t.Errorf("Float() = (%f, %v); want (100.1, <nil>)", res, err)
	}

	// OK float64
	res, err = For(source, "float64").Float(5)
	if err != nil || res != 100.1 {
		t.Errorf("Float() = (%f, %v); want (100.1, <nil>)", res, err)
	}

	// Wrong type
	res, err = For(source, "string").Float(5)
	if err != ErrWrongType || res != 5 {
		t.Errorf("Float() = (%f, %v); want (5, wrong type)", res, err)
	}

	// Missing
	res, err = For(source, "nothing").Float(5)
	if err != ErrNotFound || res != 5 {
		t.Errorf("Float() = (%f, %v); want (5, not found)", res, err)
	}

}

func TestSlice(t *testing.T) {

	def := make([]interface{}, 1)

	source := map[string]interface{}{
		"slice1": make([]interface{}, 5),
		"slice2": make([]interface{}, 0),
		"string": "test",
	}

	// OK
	res, err := For(source, "slice1").Slice(def)
	if err != nil || len(res) != 5 {
		t.Errorf("Slice() = ([%d], %v); want ([5]], <nil>)", len(res), err)
	}

	// OK
	res, err = For(source, "slice2").Slice(def)
	if err != nil || len(res) != 0 {
		t.Errorf("Slice() = ([%d], %v); want ([0], <nil>)", len(res), err)
	}

	// Wrong type
	res, err = For(source, "string").Slice(def)
	if err != ErrWrongType || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %v); want ([1], wrong type)", len(res), err)
	}

	// Missing
	res, err = For(source, "nothing").Slice(def)
	if err != ErrNotFound || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %v); want ([1], not found)", len(res), err)
	}

}

func TestCustomSlice(t *testing.T) {

	type CustomSlice []interface{}

	def := make(CustomSlice, 1)

	source := map[string]interface{}{
		"slice1": make(CustomSlice, 5),
		"slice2": make(CustomSlice, 0),
		"string": "test",
	}

	// OK
	res, err := For(source, "slice1").Slice(def)
	if err != nil || len(res) != 5 {
		t.Errorf("Slice() = ([%d], %v); want ([5]], <nil>)", len(res), err)
	}

	// OK
	res, err = For(source, "slice2").Slice(def)
	if err != nil || len(res) != 0 {
		t.Errorf("Slice() = ([%d], %v); want ([0], <nil>)", len(res), err)
	}

	// Wrong type
	res, err = For(source, "string").Slice(def)
	if err != ErrWrongType || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %v); want ([1], wrong type)", len(res), err)
	}

	// Missing
	res, err = For(source, "nothing").Slice(def)
	if err != ErrNotFound || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %v); want ([1], not found)", len(res), err)
	}

}

func TestMap(t *testing.T) {

	def := map[string]interface{}{"value00": "test"}

	source := map[string]interface{}{
		"map":    map[string]interface{}{"value11": "test", "value12": "test"},
		"string": "test",
	}

	// OK
	res, err := For(source, "map").Map(def)
	if err != nil || len(res) != 2 {
		t.Errorf("Map() = ([%d], %v); want ([2]], <nil>)", len(res), err)
	}

	// Wrong type
	res, err = For(source, "string").Map(def)
	if err != ErrWrongType || len(res) != 1 {
		t.Errorf("Map() = ([%d], %v); want ([1], wrong type)", len(res), err)
	}

	// Missing
	res, err = For(source, "nothing").Map(def)
	if err != ErrNotFound || len(res) != 1 {
		t.Errorf("Map() = ([%d], %v); want ([1], not found)", len(res), err)
	}

}

func TestCustomMap(t *testing.T) {

	type CustomMap map[string]interface{}

	def := CustomMap{"value00": "test"}

	source := CustomMap{
		"map":    CustomMap{"value11": "test", "value12": "test"},
		"string": "test",
	}

	// OK
	res, err := For(source, "map").Map(def)
	if err != nil || len(res) != 2 {
		t.Errorf("Map() = ([%d], %v); want ([2]], <nil>)", len(res), err)
	}

	// Wrong type
	res, err = For(source, "string").Map(def)
	if err != ErrWrongType || len(res) != 1 {
		t.Errorf("Map() = ([%d], %v); want ([1], wrong type)", len(res), err)
	}

	// Missing
	res, err = For(source, "nothing").Map(def)
	if err != ErrNotFound || len(res) != 1 {
		t.Errorf("Map() = ([%d], %v); want ([1], not found)", len(res), err)
	}

}

func TestExists(t *testing.T) {

	source := map[string]interface{}{
		"value1": "test",
		"value2": 0,
		"nil":    nil,
	}

	// OK
	res := For(source, "value1").Exists()
	if !res {
		t.Errorf("Exists() = (%t); want (true)", res)
	}

	// OK
	res = For(source, "value2").Exists()
	if !res {
		t.Errorf("Exists() = (%t); want (true)", res)
	}

	// Nil value is the same as missing
	res = For(source, "nil").Exists()
	if res {
		t.Errorf("Exists() = (%t); want (false)", res)
	}

	// Missing
	res = For(source, "nothing").Exists()
	if res {
		t.Errorf("Exists() = (%t); want (false)", res)
	}

}

func TestValue(t *testing.T) {

	source := map[string]interface{}{
		"value1": "test",
		"value2": 0,
	}

	// OK
	res := For(source, "value1").Value()
	if res != "test" {
		t.Errorf("Value() = (%v); want (test)", res)
	}

	// OK
	res = For(source, "value2").Value()
	if res != 0 {
		t.Errorf("Value() = (%v); want (0)", res)
	}

	// Missing
	res = For(source, "nothing").Value()
	if res != nil {
		t.Errorf("Value() = (%v); want (<nil>)", res)
	}

}
