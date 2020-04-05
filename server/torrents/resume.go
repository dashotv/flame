package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Resume(c *gin.Context) {
	infohash := c.Param("infohash")
	err := client.Resume(infohash)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func ResumeAll(c *gin.Context) {
	err := client.ResumeAll()
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
