package app

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/golem/web"
)

func (a *Application) NzbsAdd(c *gin.Context, URL, cat, name string) {
	pri, err := web.QueryDefaultInteger(c, "priority", nzbget.PriorityNormal)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := strings.Replace(string(b), "&amp;", "&", -1)

	options := nzbget.NewOptions()
	options.Category = cat
	options.Priority = pri
	if name != "" {
		options.NiceName = name
	}

	id, err := a.Nzb.Add(u, options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "id": id})
}

func (a *Application) NzbsRemove(c *gin.Context, id int) {
	var err error

	if c.Query("delete") == "true" {
		err = a.Nzb.Delete(id)
	} else {
		err = a.Nzb.Remove(id)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) NzbsDestroy(c *gin.Context, id int) {
	err := a.Nzb.Destroy(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) NzbsIndex(c *gin.Context) {
	// read the json string from cache
	res := ""
	ok, err := a.Cache.Get("flame-nzbs", &res)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.Errorf("failed to get %s from cache", "flame-nzbs")})
		return
	}

	c.String(http.StatusOK, res)
}

func (a *Application) NzbsHistory(c *gin.Context, hidden bool) {
	r, err := a.Nzb.History(hidden)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, nzbget.HistoryResponse{Response: &nzbget.Response{Timestamp: time.Now()}, Result: r})
}

func (a *Application) NzbsPause(c *gin.Context, id int) {
	var err error

	if id == -1 {
		a.NzbsPauseAll(c)
		return
	}

	err = a.Nzb.Pause(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) NzbsPauseAll(c *gin.Context) {
	err := a.Nzb.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) NzbsResume(c *gin.Context, id int) {
	var err error

	if id == -1 {
		a.NzbsResumeAll(c)
		return
	}
	err = a.Nzb.Resume(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func (a *Application) NzbsResumeAll(c *gin.Context) {
	err := a.Nzb.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
