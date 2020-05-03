package nzbs

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryInteger(c *gin.Context, name string) (int, error) {
	return QueryDefaultInteger(c, name, -1)
}

func QueryDefaultInteger(c *gin.Context, name string, def int) (int, error) {
	v := c.Query(name)
	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return def, err
	}

	if n < 0 {
		return def, fmt.Errorf("less than zero")
	}

	return n, nil
}
