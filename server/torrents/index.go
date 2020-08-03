package torrents

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

func Index(c *gin.Context) {
	// read the json string from cache
	res, err := app.Cache.Get(ctx, "flame-qbittorrents").Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, res)
}
