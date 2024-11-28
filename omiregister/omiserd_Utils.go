package register

import (
	"encoding/json"
	"strings"
)

func mapToJsonStr(data map[string]string) string {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	return string(jsonStr)
}

func splitMessage(input, delimiter string) (string, string) {
	parts := strings.SplitN(input, delimiter, 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
