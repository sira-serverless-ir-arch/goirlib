package field

import (
	"encoding/json"
	"fmt"
)

const ID = "Id"

func StringToObject(document string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(document), &result)
	if err != nil {
		panic(err)
	}
	return result
}

func GetID(input map[string]interface{}) string {
	for k, v := range input {
		if k == ID {
			return fmt.Sprintf("%s", v)
		}
	}
	return ""
}

func Flatten(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	recFlatten(input, output, "")
	return output
}

func recFlatten(input map[string]interface{}, output map[string]interface{}, prefix string) {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch value := v.(type) {
		case map[string]interface{}:
			recFlatten(value, output, key)
		default:
			output[key] = value
		}
	}
}
