package jsonutils

import (
	"encoding/json"
)

func SJsonPrint(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func SJsonPrettyPrint(obj interface{}) string {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}
