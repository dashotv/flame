package torrents

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context, URL string) {
	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	infohash, err := app.Utorrent.Add(string(b))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "infohash": infohash})
}

func Remove(c *gin.Context, infohash string, delete bool) {
	err := app.Utorrent.Remove(infohash, delete)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
