package traytor

import "encoding/json"

// jsonObjectType returns the value of a "type" key in a json object
func jsonObjectType(data []byte) (string, error) {
	typeHolder := &struct {
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(data, &typeHolder)
	return typeHolder.Type, err
}
