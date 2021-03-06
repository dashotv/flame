package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Label(c *gin.Context, infohash, label string) {
	err := app.Utorrent.Label(infohash, label)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
