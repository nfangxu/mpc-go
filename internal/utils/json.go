package utils

import "encoding/json"

func Json(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
