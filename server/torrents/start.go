package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(c *gin.Context) {
	infohash := c.Param("infohash")
	err := client.Start(infohash)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Stop(c *gin.Context) {
	infohash := c.Param("infohash")
	err := client.Stop(infohash)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
