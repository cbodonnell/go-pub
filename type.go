package main

import (
	"encoding/json"
)

func Read(b []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(b, &result)
	return result
}

func GetType(o interface{}) interface{} {
	switch t := o.(type) {
	default:
		return t
	}
}

// TODO: Can make isIRI as well using url.Parse(...)
func IsString(o interface{}) bool {
	_, isString := o.(string)
	return isString
}
