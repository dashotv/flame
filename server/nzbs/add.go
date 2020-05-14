package nzbs

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/dashotv/flame/nzbget"
	"github.com/dashotv/golem/web"
)

func Add(c *gin.Context, URL, cat, name string) {
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

	id, err := app.Nzbget.Add(u, options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "id": id})
}

func Remove(c *gin.Context, id int) {
	var err error

	if c.Query("delete") == "true" {
		err = app.Nzbget.Delete(id)
	} else {
		err = app.Nzbget.Remove(id)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Destroy(c *gin.Context, id int) {
	var err error

	err = app.Nzbget.Destroy(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
