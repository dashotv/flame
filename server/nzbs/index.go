package nzbs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"github.com/dashotv/flame/nzbget"
)

var cache *redis.Client
var client *nzbget.Client

func Routes(red *redis.Client, c *nzbget.Client, e *gin.Engine) {
	cache = red
	client = c
	r := e.Group("/nzbs")
	r.GET("/", Index)
	r.GET("/add", Add)
	r.GET("/remove", Remove)
	r.GET("/pause", Pause)
	r.GET("/resume", Resume)
}

func Index(c *gin.Context) {
	// read the json string from cache
	res, err := cache.Get("flame-nzbs").Result()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.String(http.StatusOK, res)
}
