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
	res, err := For(object, "a[0].b.c").Int(0)
	fmt.Println(res, err)

	// Attempt extracting a string at path .d that does not exist
	res2, err := For(object, "a[0].b.d").String("default")
	fmt.Println(res2, err)

	// Output:
	// 3 <nil>
	// default not found
}
