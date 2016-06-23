package jsonutil

import "encoding/json"

// ObjectType returns the value of a "type" key in a json object
func ObjectType(data []byte) (string, error) {
	typeHolder := &struct {
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(data, &typeHolder)
	return typeHolder.Type, err
}
