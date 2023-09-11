package tool

import (
	"encoding/json"
)

func SimpleJson(data interface{}) string {
	d, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(d)
}

func Json(data interface{}) (string, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
