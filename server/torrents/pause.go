package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Pause(c *gin.Context) {
	infohash := c.Query("infohash")
	if infohash == "" {
		PauseAll(c)
		return
	}
	err := client.Pause(infohash)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func PauseAll(c *gin.Context) {
	err := client.PauseAll()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
