package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var cache *redis.Client

func Routes(c *redis.Client, e *gin.Engine) {
	cache = c
	r := e.Group("/torrents")
	r.GET("/", Torrents)
}

func Torrents(c *gin.Context) {
	// read the json string from cache
	res, err := cache.Get("flame").Result()
	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusOK, res)
}
