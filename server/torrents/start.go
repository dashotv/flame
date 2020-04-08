package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(c *gin.Context) {
	infohash := c.Query("infohash")
	err := client.Start(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Stop(c *gin.Context) {
	infohash := c.Query("infohash")
	err := client.Stop(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
