package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// QueryString retrieves a string param from the gin request querystring
func QueryString(c *gin.Context, name string) string {
	return c.Query(name)
}

// QueryInt retrieves an integer param from the gin request querystring
func QueryInt(c *gin.Context, name string) int {
	v := c.Query(name)
	i, _ := strconv.Atoi(v)
	return i
}

// QueryDefaultInt retrieves an integer param from the gin request querystring
// defaults to def argument if not found
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

// QueryBool retrieves a boolean param from the gin request querystring
func QueryBool(c *gin.Context, name string) bool {
	return c.Query(name) == "true"
}

// WithTimeout runs a delegate function with a timeout,
//
// Example: Wait for a channel
//
//	if value, ok := WithTimeout(func()interface{}{return <- inbox}, time.Second); ok {
//	    // returned
//	} else {
//	    // didn't return
//	}
//
// Example: To send to a channel
//
//	_, ok := WithTimeout(func()interface{}{outbox <- myValue; return nil}, time.Second)
//	if !ok {
//	    // didn't send
//	}
func WithTimeout(delegate func() interface{}, timeout time.Duration) (ret interface{}, ok bool) {
	ch := make(chan interface{}, 1) // buffered
	go func() { ch <- delegate() }()
	select {
	case ret = <-ch:
		return ret, true
	case <-time.After(timeout):
	}
	return nil, false
}
