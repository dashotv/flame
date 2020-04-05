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
	r.GET("/pause", PauseAll)
	r.GET("/resume", ResumeAll)
	r.GET("/:infohash/remove", Remove)
	r.GET("/:infohash/pause", Pause)
	r.GET("/:infohash/resume", Resume)
	r.GET("/:infohash/start", Start)
	r.GET("/:infohash/stop", Stop)
	r.GET("/:infohash/label", Label)
	r.GET("/:infohash/want/none", WantNone)
	r.GET("/:infohash/want/:files", Want)
	r.GET("/:infohash/wanted", Wanted)
}

func Index(c *gin.Context) {
	// read the json string from cache
	res, err := cache.Get("flame").Result()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.String(http.StatusOK, res)
}
