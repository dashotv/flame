package nzbs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	URL := c.Query("url")
	id, err := client.Add(URL)
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
