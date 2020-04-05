package torrents

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	URL := c.Query("url")
	infohash, err := client.Add(URL)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "infohash": infohash})
}

func Remove(c *gin.Context) {
	infohash := c.Param("infohash")
	delete := c.Query("delete") == "true"
	err := client.Remove(infohash, delete)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
