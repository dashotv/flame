package flame

import (
	"fmt"
	"strings"
)

func query(params map[string]string) (string) {
	list := []string{}
	for k, v := range params {
		list = append(list, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(list, "&")
}
