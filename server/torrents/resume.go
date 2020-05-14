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
	_, err := app.Qbittorrent.Resume(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ResumeAll(c *gin.Context) {
	_, err := app.Qbittorrent.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
