package app

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/dashotv/golem/web"
	"github.com/gin-gonic/gin"

	"github.com/dashotv/flame/nzbget"
)

func NzbsAdd(c *gin.Context, URL, cat, name string) {
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

	id, err := App().Nzbget.Add(u, options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "id": id})
}

func NzbsRemove(c *gin.Context, id int) {
	var err error

	if c.Query("delete") == "true" {
		err = App().Nzbget.Delete(id)
	} else {
		err = App().Nzbget.Remove(id)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func NzbsDestroy(c *gin.Context, id int) {
	err := App().Nzbget.Destroy(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}

func NzbsIndex(c *gin.Context) {
	// read the json string from cache
	res, err := App().Cache.Get(ctx, "flame-nzbs").Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.String(http.StatusOK, res)
}

func NzbsHistory(c *gin.Context, hidden bool) {
	r, err := App().Nzbget.History(hidden)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, nzbget.HistoryResponse{Response: &nzbget.Response{Timestamp: time.Now()}, Result: r})
}

func NzbsPause(c *gin.Context, id int) {
	var err error

	if id == -1 {
		NzbsPauseAll(c)
		return
	}

	err = App().Nzbget.Pause(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func NzbsPauseAll(c *gin.Context) {
	err := App().Nzbget.PauseAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func NzbsResume(c *gin.Context, id int) {
	var err error

	if id == -1 {
		NzbsResumeAll(c)
		return
	}
	err = App().Nzbget.Resume(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func NzbsResumeAll(c *gin.Context) {
	err := App().Nzbget.ResumeAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}
