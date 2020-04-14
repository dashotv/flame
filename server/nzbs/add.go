package nzbs

import (
	"encoding/base64"
	"net/http"

	"github.com/dashotv/flame/nzbget"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	URL := c.Query("url")
	cat := c.Query("category")
	pri, err := QueryInteger(c, "priority")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	options := nzbget.NewOptions()
	options.Category = cat
	options.Priority = pri

	id, err := client.Add(string(b), options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "id": id})
}

func Remove(c *gin.Context) {
	id, err := QueryInteger(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.Query("delete") == "true" {
		err = client.Delete(id)
	} else {
		err = client.Remove(id)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Destroy(c *gin.Context) {
	id, err := QueryInteger(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = client.Destroy(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
