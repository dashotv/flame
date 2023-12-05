package app

import (
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

	infohash, err := qb.Add(u, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("added: %s", infohash)

	c.JSON(http.StatusOK, gin.H{"error": false, "infohash": infohash})
}

func QbittorrentsRemove(c *gin.Context, infohash string, del bool) {
	err := qb.Delete(infohash, del)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsIndex(c *gin.Context) {
	// read the json string from cache
	res, err := cache.Get(ctx, "flame-qbittorrents").Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, res)
}

func QbittorrentsLabel(c *gin.Context, infohash, label string) {
	err := qb.SetLabel([]string{infohash}, label)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsPause(c *gin.Context, infohash string) {
	if infohash == "" {
		QbittorrentsPauseAll(c)
		return
	}
	err := qb.Pause(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsPauseAll(c *gin.Context) {
	err := qb.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsResume(c *gin.Context, infohash string) {
	if infohash == "" {
		QbittorrentsResumeAll(c)
		return
	}

	err := qb.Resume(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsResumeAll(c *gin.Context) {
	err := qb.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

//func QbittorrentsStart(c *gin.Context, infohash string) {
//	err := qb.Start(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}
//
//func QbittorrentsStop(c *gin.Context, infohash string) {
//	err := App().Utorrent.Stop(infohash)
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
	err := qb.Want(infohash, ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsWantNone(c *gin.Context, infohash string) {
	err := qb.WantNone(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func QbittorrentsWanted(c *gin.Context, infohash string) {
	want, err := qb.WantedCount(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "wanted": want > 0, "count": want})
}
