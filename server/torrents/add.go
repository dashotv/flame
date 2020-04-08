package torrents

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	URL := c.Query("url")
	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	infohash, err := client.Add(string(b))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "infohash": infohash})
}

func Remove(c *gin.Context) {
	infohash := c.Query("infohash")
	delete := c.Query("delete") == "true"
	err := client.Remove(infohash, delete)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
