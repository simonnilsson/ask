package ask

import (
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

func TestString(t *testing.T) {

	source := map[string]interface{}{
		"value1": "test",
		"value2": 100,
	}

	// OK
	res, ok := For(source, "value1").String("default")
	if !ok || res != "test" {
		t.Errorf(`String() = ("%s", %t); want ("test", true)`, res, ok)
	}

	// Wrong type
	res, ok = For(source, "value2").String("default")
	if ok || res != "default" {
		t.Errorf(`String() = ("%s", %t); want ("default", false)`, res, ok)
	}

	// Missing
	res, ok = For(source, "value3").String("default")
	if ok || res != "default" {
		t.Errorf(`String() = ("%s", %t); want ("default", false)`, res, ok)
	}

}

func TestBool(t *testing.T) {

	source := map[string]interface{}{
		"value1": true,
		"value2": false,
		"value3": 100,
	}

	// OK
	res, ok := For(source, "value1").Bool(false)
	if !ok || res != true {
		t.Errorf(`Bool() = (%t, %t); want (true, true)`, res, ok)
	}

	// OK
	res, ok = For(source, "value2").Bool(false)
	if !ok || res != false {
		t.Errorf(`Bool() = (%t, %t); want (false, true)`, res, ok)
	}

	// Wrong type
	res, ok = For(source, "value3").Bool(false)
	if ok || res != false {
		t.Errorf(`Bool() = (%t, %t); want (false, false)`, res, ok)
	}

	// Missing
	res, ok = For(source, "value4").Bool(false)
	if ok || res != false {
		t.Errorf(`Bool() = (%t, %t); want (false, false)`, res, ok)
	}

}

func TestInt(t *testing.T) {

	source := map[string]interface{}{
		"value1": 100,
		"value2": -100,
		"value3": "test",
	}

	// OK
	res, ok := For(source, "value1").Int(5)
	if !ok || res != 100 {
		t.Errorf("Int() = (%d, %t); want (100, true)", res, ok)
	}

	// OK negative
	res, ok = For(source, "value2").Int(5)
	if !ok || res != -100 {
		t.Errorf("Int() = (%d, %t); want (-100, true)", res, ok)
	}

	// Wrong type
	res, ok = For(source, "value3").Int(5)
	if ok || res != 5 {
		t.Errorf("Int() = (%d, %t); want (5, false)", res, ok)
	}

	// Missing
	res, ok = For(source, "value4").Int(5)
	if ok || res != 5 {
		t.Errorf("Int() = (%d, %t); want (5, false)", res, ok)
	}

}

func TestUint(t *testing.T) {

	source := map[string]interface{}{
		"value1": 100,
		"value2": -100,
		"value3": "test",
	}

	// OK
	res, ok := For(source, "value1").Uint(5)
	if !ok || res != 100 {
		t.Errorf("Uint() = (%d, %t); want (100, true)", res, ok)
	}

	// Fail on negative
	res, ok = For(source, "value2").Uint(5)
	if ok || res != 5 {
		t.Errorf("Uint() = (%d, %t); want (5, false)", res, ok)
	}

	// Wrong type
	res, ok = For(source, "value3").Uint(5)
	if ok || res != 5 {
		t.Errorf("Uint() = (%d, %t); want (5, false)", res, ok)
	}

	// Missing
	res, ok = For(source, "value4").Uint(5)
	if ok || res != 5 {
		t.Errorf("Uint() = (%d, %t); want (5, false)", res, ok)
	}

}

func TestFloat(t *testing.T) {

	source := map[string]interface{}{
		"value1": 100.10,
		"value2": -100.10,
		"value3": "test",
	}

	// OK
	res, ok := For(source, "value1").Float(5)
	if !ok || res != 100.10 {
		t.Errorf("Float() = (%f, %t); want (100.10, true)", res, ok)
	}

	// OK negative
	res, ok = For(source, "value2").Float(5)
	if !ok || res != -100.10 {
		t.Errorf("Float() = (%f, %t); want (-100.10, true)", res, ok)
	}

	// Wrong type
	res, ok = For(source, "value3").Float(5)
	if ok || res != 5 {
		t.Errorf("Float() = (%f, %t); want (5, false)", res, ok)
	}

	// Missing
	res, ok = For(source, "value4").Float(5)
	if ok || res != 5 {
		t.Errorf("Float() = (%f, %t); want (5, false)", res, ok)
	}

}

func TestSlice(t *testing.T) {

	def := make([]interface{}, 1)

	source := map[string]interface{}{
		"value1": make([]interface{}, 5),
		"value2": make([]interface{}, 0),
		"value3": "test",
	}

	// OK
	res, ok := For(source, "value1").Slice(def)
	if !ok || len(res) != 5 {
		t.Errorf("Slice() = ([%d], %t); want ([5]], true)", len(res), ok)
	}

	// OK
	res, ok = For(source, "value2").Slice(def)
	if !ok || len(res) != 0 {
		t.Errorf("Slice() = ([%d], %t); want ([0], true)", len(res), ok)
	}

	// Wrong type
	res, ok = For(source, "value3").Slice(def)
	if ok || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %t); want ([1], false)", len(res), ok)
	}

	// Missing
	res, ok = For(source, "value4").Slice(def)
	if ok || len(res) != 1 {
		t.Errorf("Slice() = ([%d], %t); want ([1], false)", len(res), ok)
	}

}

func TestMap(t *testing.T) {

	def := map[string]interface{}{"value00": "test"}

	source := map[string]interface{}{
		"value1": map[string]interface{}{"value11": "test", "value12": "test"},
		"value2": "test",
	}

	// OK
	res, ok := For(source, "value1").Map(def)
	if !ok || len(res) != 2 {
		t.Errorf("Map() = ([%d], %t); want ([2]], true)", len(res), ok)
	}

	// Wrong type
	res, ok = For(source, "value2").Map(def)
	if ok || len(res) != 1 {
		t.Errorf("Map() = ([%d], %t); want ([1], false)", len(res), ok)
	}

	// Missing
	res, ok = For(source, "value3").Map(def)
	if ok || len(res) != 1 {
		t.Errorf("Map() = ([%d], %t); want ([1], false)", len(res), ok)
	}

}

func TestExists(t *testing.T) {

	source := map[string]interface{}{
		"value1": "test",
		"value2": 0,
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

	// Missing
	res = For(source, "value3").Exists()
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
		t.Errorf("Exists() = (%v); want (test)", res)
	}

	// OK
	res = For(source, "value2").Value()
	if res != 0 {
		t.Errorf("Exists() = (%v); want (0)", res)
	}

	// Missing
	res = For(source, "value3").Value()
	if res != nil {
		t.Errorf("Exists() = (%v); want (nil)", res)
	}

}
