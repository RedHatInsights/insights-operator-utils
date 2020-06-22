package helpers

import "encoding/json"

// ToJSONString converts anything to JSON or panics if it's not possible
func ToJSONString(obj interface{}) string {
	jsonBytes, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		panic(err)
	}

	return string(jsonBytes)
}
