package server

import "encoding/json"

func parseError(err string) string {
	response := Response{Err: err}
	jsonStr, _ := json.Marshal(response)
	return string(jsonStr)
}
