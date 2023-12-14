package app

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (a *Application) QbittorrentsAdd(c *gin.Context, URL string) {
	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u := string(b)

	infohash, err := app.Qbt.Add(u, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.Log.Infof("added: %s", infohash)

	c.JSON(http.StatusOK, gin.H{"error": false, "infohash": infohash})
}

func (a *Application) QbittorrentsRemove(c *gin.Context, infohash string, del bool) {
	err := app.Qbt.Delete(infohash, del)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsIndex(c *gin.Context) {
	// read the json string from cache
	// read the json string from cache
	res := ""
	ok, err := a.Cache.Get("flame-qbittorrents", &res)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.Errorf("failed to get %s from cache", "flame-qbittorrents")})
		return
	}

	c.String(http.StatusOK, res)
}

func (a *Application) QbittorrentsLabel(c *gin.Context, infohash, label string) {
	err := app.Qbt.SetLabel([]string{infohash}, label)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsPause(c *gin.Context, infohash string) {
	if infohash == "" {
		a.QbittorrentsPauseAll(c)
		return
	}
	err := app.Qbt.Pause(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsPauseAll(c *gin.Context) {
	err := app.Qbt.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsResume(c *gin.Context, infohash string) {
	if infohash == "" {
		a.QbittorrentsResumeAll(c)
		return
	}

	err := app.Qbt.Resume(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsResumeAll(c *gin.Context) {
	err := app.Qbt.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

//func (a *Application) QbittorrentsStart(c *gin.Context, infohash string) {
//	err := app.Qbt.Start(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}
//
//func (a *Application) QbittorrentsStop(c *gin.Context, infohash string) {
//	err := App().Utorrent.Stop(infohash)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"error": false})
//}

func (a *Application) QbittorrentsWant(c *gin.Context, infohash, files string) {
	if files == "none" {
		a.QbittorrentsWantNone(c, infohash)
		return
	}

	ids := strings.Split(files, ",")
	err := app.Qbt.Want(infohash, ids)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsWantNone(c *gin.Context, infohash string) {
	err := app.Qbt.WantNone(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) QbittorrentsWanted(c *gin.Context, infohash string) {
	want, err := app.Qbt.WantedCount(infohash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "wanted": want > 0, "count": want})
}
