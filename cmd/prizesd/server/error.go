package server

import (
	"encoding/json"
)

func parseResult(v interface{}) (string, bool) {
	if stderr, ok := v.(error); ok {
		response := Response{Error: stderr.Error()}
		jsonStr, _ := json.Marshal(response)
		return string(jsonStr), false
	} else {
		response := Response{Result: v}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			out, _ := parseResult(err)
			return string(out), false
		}
		return string(jsonResponse), true
	}
}
func parseError(stderr error) string {
	response := Response{Error: stderr.Error()}
	jsonStr, _ := json.Marshal(response)
	return string(jsonStr)
}
