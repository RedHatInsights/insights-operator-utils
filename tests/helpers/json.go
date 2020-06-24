package helpers

import "encoding/json"

// ToJSONString converts anything to JSON or panics if it's not possible
func ToJSONString(obj interface{}) string {
	return toJSONString(obj, false)
}

// ToJSONPrettyString converts anything to indented JSON or panics if it's not possible
func ToJSONPrettyString(obj interface{}) string {
	return toJSONString(obj, true)
}

// toJSONString converts anything to JSON or panics if it's not possible
// isOutputPretty makes output indented
func toJSONString(obj interface{}, isOutputPretty bool) string {
	var (
		jsonBytes []byte
		err       error
	)
	if isOutputPretty {
		jsonBytes, err = json.MarshalIndent(obj, "", "\t")
	} else {
		jsonBytes, err = json.Marshal(obj)
	}
	if err != nil {
		panic(err)
	}

	return string(jsonBytes)
}
