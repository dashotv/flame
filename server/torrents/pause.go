package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Pause(c *gin.Context, infohash string) {
	if infohash == "" {
		PauseAll(c)
		return
	}

	err := app.Utorrent.Pause(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func PauseAll(c *gin.Context) {
	err := app.Utorrent.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
