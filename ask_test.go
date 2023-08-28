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
	res, ok := For(source, "string").String("default")
	if !ok || res != "test" {
		t.Errorf(`String() = ("%s", %t); want ("test", true)`, res, ok)
	}

	// Wrong type
	res, ok = For(source, "number").String("default")
	if ok || res != "default" {
		t.Errorf(`String() = ("%s", %t); want ("default", false)`, res, ok)
	}

	// Missing
	res, ok = For(source, "nothing").String("default")
	if ok || res != "default" {
		t.Errorf(`String() = ("%s", %t); want ("default", false)`, res, ok)
	}

}

func TestBool(t *testing.T) {

	source := map[string]interface{}{
		"bool1":  true,
		"bool2":  false,
		"number": 100,
	}

	// OK
	res, ok := For(source, "bool1").Bool(false)
	if !ok || res != true {
		t.Errorf(`Bool() = (%t, %t); want (true, true)`, res, ok)
	}

	// OK
	res, ok = For(source, "bool2").Bool(false)
	if !ok || res != false {
		t.Errorf(`Bool() = (%t, %t); want (false, true)`, res, ok)
	}

	// Wrong type
	res, ok = For(source, "number").Bool(false)
	if ok || res != false {
		t.Errorf(`Bool() = (%t, %t); want (false, false)`, res, ok)
	}

	// Missing
	res, ok = For(source, "nothing").Bool(false)
	if ok || res != false {
		t.Errorf(`Bool() = (%t, %t); want (false, false)`, res, ok)
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
	res, ok := For(source, "positive").Int(5)
	if !ok || res != 100 {
		t.Errorf("Int() = (%d, %t); want (100, true)", res, ok)
	}

	// OK negative
	res, ok = For(source, "negative").Int(5)
	if !ok || res != -100 {
		t.Errorf("Int() = (%d, %t); want (-100, true)", res, ok)
	}

	// OK unsigned
	res, ok = For(source, "unsigned").Int(5)
	if !ok || res != 100 {
		t.Errorf("Int() = (%d, %t); want (100, true)", res, ok)
	}

	// OK float
	res, ok = For(source, "float").Int(5)
	if !ok || res != 100 {
		t.Errorf("Int() = (%d, %t); want (100, true)", res, ok)
	}

	// OK negative float
	res, ok = For(source, "negative_float").Int(5)
	if !ok || res != -100 {
		t.Errorf("Int() = (%d, %t); want (-100, true)", res, ok)
	}

	// Too big number
	res, ok = For(source, "toobig").Int(5)
	if ok || res != 5 {
		t.Errorf("Int() = (%d, %t); want (5, false)", res, ok)
	}

	// Wrong type
	res, ok = For(source, "string").Int(5)
	if ok || res != 5 {
		t.Errorf("Int() = (%d, %t); want (5, false)", res, ok)
	}

	// Missing
	res, ok = For(source, "nothing").Int(5)
	if ok || res != 5 {
		t.Errorf("Int() = (%d, %t); want (5, false)", res, ok)
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
	res, ok := For(source, "positive").Uint(5)
	if !ok || res != 100 {
		t.Errorf("Uint() = (%d, %t); want (100, true)", res, ok)
	}

	// OK unsigned
	res, ok = For(source, "unsigned").Uint(5)
	if !ok || res != 100 {
		t.Errorf("Uint() = (%d, %t); want (100, true)", res, ok)
	}

	// OK float
	res, ok = For(source, "float").Uint(5)
	if !ok || res != 100 {
		t.Errorf("Uint() = (%d, %t); want (100, true)", res, ok)
	}

	// Fail on negative
	res, ok = For(source, "negative").Uint(5)
	if ok || res != 5 {
		t.Errorf("Uint() = (%d, %t); want (5, false)", res, ok)
	}

	// Wrong type
	res, ok = For(source, "string").Uint(5)
	if ok || res != 5 {
		t.Errorf("Uint() = (%d, %t); want (5, false)", res, ok)
	}

	// Missing
	res, ok = For(source, "nothing").Uint(5)
	if ok || res != 5 {
		t.Errorf("Uint() = (%d, %t); want (5, false)", res, ok)
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
	res, ok := For(source, "positive").Float(5)
	if !ok || res != 100 {
		t.Errorf("Float() = (%f, %t); want (100, true)", res, ok)
	}

	// OK negative
	res, ok = For(source, "negative").Float(5)
	if !ok || res != -100 {
		t.Errorf("Float() = (%f, %t); want (-100, true)", res, ok)
	}

	// OK unsigned
	res, ok = For(source, "unsigned").Float(5)
	if !ok || res != 100 {
		t.Errorf("Float() = (%f, %t); want (100, true)", res, ok)
	}

	// OK float32
	res, ok = For(source, "float32").Float(5)
	if !ok || math.Abs(res-100.1) > .00001 {
		t.Errorf("Float() = (%f, %t); want (100.1, true)", res, ok)
	}

	// OK float64
	res, ok = For(source, "float64").Float(5)
	if !ok || res != 100.1 {
		t.Errorf("Float() = (%f, %t); want (100.1, true)", res, ok)
	}

	// Wrong type
	res, ok = For(source, "string").Float(5)
	if ok || res != 5 {
		t.Errorf("Float() = (%f, %t); want (5, false)", res, ok)
	}

	// Missing
	res, ok = For(source, "nothing").Float(5)
	if ok || res != 5 {
		t.Errorf("Float() = (%f, %t); want (5, false)", res, ok)
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
	res, ok := For(source, "slice1").Slice(def)
	if !ok || len(res) != 5 {
		t.Errorf("Slice() = ([%d], %t); want ([5]], true)", len(res), ok)
	}

	// OK
	res, ok = For(source, "slice2").Slice(def)
	if !ok || len(res) != 0 {
		t.Errorf("Slice() = ([%d], %t); want ([0], true)", len(res), ok)
	}

	// Wrong type
	res, ok = For(source, "string").Slice(def)
	if ok || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %t); want ([1], false)", len(res), ok)
	}

	// Missing
	res, ok = For(source, "nothing").Slice(def)
	if ok || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %t); want ([1], false)", len(res), ok)
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
	res, ok := For(source, "slice1").Slice(def)
	if !ok || len(res) != 5 {
		t.Errorf("Slice() = ([%d], %t); want ([5]], true)", len(res), ok)
	}

	// OK
	res, ok = For(source, "slice2").Slice(def)
	if !ok || len(res) != 0 {
		t.Errorf("Slice() = ([%d], %t); want ([0], true)", len(res), ok)
	}

	// Wrong type
	res, ok = For(source, "string").Slice(def)
	if ok || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %t); want ([1], false)", len(res), ok)
	}

	// Missing
	res, ok = For(source, "nothing").Slice(def)
	if ok || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %t); want ([1], false)", len(res), ok)
	}

}

func TestMap(t *testing.T) {

	def := map[string]interface{}{"value00": "test"}

	source := map[string]interface{}{
		"map":    map[string]interface{}{"value11": "test", "value12": "test"},
		"string": "test",
	}

	// OK
	res, ok := For(source, "map").Map(def)
	if !ok || len(res) != 2 {
		t.Errorf("Map() = ([%d], %t); want ([2]], true)", len(res), ok)
	}

	// Wrong type
	res, ok = For(source, "string").Map(def)
	if ok || len(res) != 1 {
		t.Errorf("Map() = ([%d], %t); want ([1], false)", len(res), ok)
	}

	// Missing
	res, ok = For(source, "nothing").Map(def)
	if ok || len(res) != 1 {
		t.Errorf("Map() = ([%d], %t); want ([1], false)", len(res), ok)
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
	res, ok := For(source, "map").Map(def)
	if !ok || len(res) != 2 {
		t.Errorf("Map() = ([%d], %t); want ([2]], true)", len(res), ok)
	}

	// Wrong type
	res, ok = For(source, "string").Map(def)
	if ok || len(res) != 1 {
		t.Errorf("Map() = ([%d], %t); want ([1], false)", len(res), ok)
	}

	// Missing
	res, ok = For(source, "nothing").Map(def)
	if ok || len(res) != 1 {
		t.Errorf("Map() = ([%d], %t); want ([1], false)", len(res), ok)
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
		t.Errorf("Value() = (%v); want (nil)", res)
	}

}
