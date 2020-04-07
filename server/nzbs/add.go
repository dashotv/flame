package nzbs

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	URL := c.Query("url")
	b, err := base64.StdEncoding.DecodeString(URL)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := client.Add(string(b))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "id": id})
}

func Remove(c *gin.Context) {
	id, err := QueryInteger(c, "id")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if c.Query("delete") == "true" {
		err = client.Delete(id)
	} else {
		err = client.Remove(id)
	}
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false})
}

func Destroy(c *gin.Context) {
	id, err := QueryInteger(c, "id")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = client.Destroy(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false})
}
