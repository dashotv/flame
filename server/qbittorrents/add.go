package qbittorrents

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
	u := string(b)

	infohash, err := app.Qbittorrent.Add(u, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app.Log.Infof("added: %s", infohash)

	c.JSON(http.StatusOK, gin.H{"error": false, "infohash": infohash})
}

func Remove(c *gin.Context, infohash string, del bool) {
	var err error

	if del {
		_, err = app.Qbittorrent.DeletePermanently([]string{infohash})
	} else {
		_, err = app.Qbittorrent.DeleteTemp([]string{infohash})
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
