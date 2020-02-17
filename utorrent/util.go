package utorrent

import (
	"fmt"
	"strings"
)

func query(params map[string]string) string {
	list := []string{}
	for k, v := range params {
		list = append(list, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(list, "&")
}

func getString(value interface{}) string {
	if value != nil {
		return value.(string)
	}
	return ""
}

func getFloat64(value interface{}) float64 {
	if value != nil {
		return value.(float64)
	}
	return 0.0
}

func getInt(value interface{}) int {
	if value != nil {
		return value.(int)
	}
	return 0
}
