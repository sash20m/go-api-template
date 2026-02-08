package utils

import (
	"encoding/json"
	"net/http"
)

func GetJSONBody[T any](r *http.Request) (T, error) {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return body, err
	}
	return body, nil
}
