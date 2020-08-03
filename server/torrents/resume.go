package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Resume(c *gin.Context, infohash string) {
	if infohash == "" {
		ResumeAll(c)
		return
	}

	err := app.Utorrent.Resume(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ResumeAll(c *gin.Context) {
	err := app.Utorrent.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
