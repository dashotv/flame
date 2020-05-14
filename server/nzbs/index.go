package nzbs

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/flame/nzbget"
)

func Index(c *gin.Context) {
	// read the json string from cache
	res, err := app.Cache.Get("flame-nzbs").Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.String(http.StatusOK, res)
}

func History(c *gin.Context, hidden bool) {
	r, err := app.Nzbget.History(hidden)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, nzbget.HistoryResponse{Response: &nzbget.Response{Timestamp: time.Now()}, Result: r})
}
