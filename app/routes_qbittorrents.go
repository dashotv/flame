package app

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func QbittorrentsAdd(c *gin.Context, URL string) {
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

func QbittorrentsRemove(c *gin.Context, infohash string, del bool) {
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

var ctx = context.Background()

func QbittorrentsIndex(c *gin.Context) {
	// read the json string from cache
	res, err := app.Cache.Get(ctx, "flame-qbittorrents").Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, res)
}

func QbittorrentsLabel(c *gin.Context, infohash, label string) {
	_, err := app.Qbittorrent.SetLabel([]string{infohash}, label)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsPause(c *gin.Context, infohash string) {
	if infohash == "" {
		PauseAll(c)
		return
	}
	_, err := app.Qbittorrent.Pause(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsPauseAll(c *gin.Context) {
	_, err := app.Qbittorrent.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsResume(c *gin.Context, infohash string) {
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

func QbittorrentsResumeAll(c *gin.Context) {
	_, err := app.Qbittorrent.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

//func QbittorrentsStart(c *gin.Context, infohash string) {
//	err := app.Qbittorrent.Start(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}
//
//func QbittorrentsStop(c *gin.Context, infohash string) {
//	err := app.Utorrent.Stop(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}

func QbittorrentsWant(c *gin.Context, infohash, files string) {
	if files == "none" {
		QbittorrentsWantNone(c, infohash)
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

func QbittorrentsWantNone(c *gin.Context, infohash string) {
	err := app.Qbittorrent.WantNone(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsWanted(c *gin.Context, infohash string) {
	want, err := app.Qbittorrent.WantedCount(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "wanted": want > 0, "count": want})
}
