package ask

import (
	"encoding/json"
	"fmt"
)

func Example() {

	// Use parsed JSON as source data
	var object map[string]interface{}
	json.Unmarshal([]byte(`{
		"a": [{
			"b": {
				"c": 3
			}
		}]
 	}`), &object)

	// Extract the 3
	res, ok := For(object, "a[0].b.c").Int(0)
	fmt.Println(res, ok)

	// Attempt extracting a string at path .d that does not exist
	res2, ok := For(object, "a[0].b.d").String("nothing")
	fmt.Println(res2, ok)

	// Output:
	// 3 true
	// nothing false
}
