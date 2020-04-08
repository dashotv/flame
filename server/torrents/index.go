package torrents

import (
	"net/http"

	"github.com/dashotv/flame/utorrent"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var cache *redis.Client
var client *utorrent.Client

func Routes(red *redis.Client, c *utorrent.Client, e *gin.Engine) {
	cache = red
	client = c
	r := e.Group("/torrents")
	r.GET("/", Index)
	r.GET("/add", Add)
	r.GET("/remove", Remove)
	r.GET("/pause", Pause)
	r.GET("/resume", Resume)
	r.GET("/start", Start)
	r.GET("/stop", Stop)
	r.GET("/label", Label)
	r.GET("/want", Want)
	r.GET("/wanted", Wanted)
}

func Index(c *gin.Context) {
	// read the json string from cache
	res, err := cache.Get("flame-torrents").Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, res)
}
