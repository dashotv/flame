package qbittorrents

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Want(c *gin.Context, infohash, files string) {
	if files == "none" {
		WantNone(c, infohash)
		return
	}

	ids := strings.Split(files, ",")
	err := app.Qbittorrent.Want(infohash, ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func WantNone(c *gin.Context, infohash string) {
	err := app.Qbittorrent.WantNone(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Wanted(c *gin.Context, infohash string) {
	want, err := app.Qbittorrent.Wanted(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "wanted": want})
}
